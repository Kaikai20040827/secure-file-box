package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Kaikai20040827/graduation/internal/model"
	"gorm.io/gorm"
)

type FileService struct {
	db      *gorm.DB
	dirpath string
	encKey  []byte
	macKey  []byte
}

const (
	fileMagic = "SFB1"
	headerLen = 4 + 16 // magic + iv
	macLen    = 32
	streamBuf = 32 * 1024
)

func NewFileService(db *gorm.DB, storagePath string, base64Key string) *FileService {
	_ = os.MkdirAll(storagePath, 0755)
	fmt.Println("✓ Creating a new file service done")
	encKey, macKey := deriveKeys(base64Key)
	return &FileService{db: db, dirpath: storagePath, encKey: encKey, macKey: macKey}
}

func (f *FileService) UploadFile(fileReader io.Reader, filename string, uploaderID uint, description string) (*model.File, error) {
	storedName := fmt.Sprintf("%d_%d_%s", uploaderID, time.Now().UnixNano(), filename)
	dst := filepath.Join(f.dirpath, storedName)
	size, err := f.encryptToFile(fileReader, dst)
	if err != nil {
		return nil, err
	}

	file := &model.File{
		Filename:    filename,
		StoragePath: dst,
		Size:        size,
		Description: description,
		UploaderID:  fmt.Sprintf("%d", uploaderID),
		CreatedAt:   time.Now(),
	}

	if err := f.db.Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func (f *FileService) SaveUserAvatar(fileReader io.Reader, filename string, userID uint) (string, int64, error) {
	storedName := fmt.Sprintf("avatar_%d_%d_%s", userID, time.Now().UnixNano(), filename)
	dst := filepath.Join(f.dirpath, storedName)
	size, err := f.encryptToFile(fileReader, dst)
	if err != nil {
		return "", 0, err
	}
	return dst, size, nil
}

func (f *FileService) RemoveStoredFile(path string) error {
	if path == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (f *FileService) UpdateFile(id uint, fileReader io.Reader, filename *string, description *string) (*model.File, error) {
	var file model.File
	if err := f.db.First(&file, id).Error; err != nil {
		return nil, err
	}

	if fileReader != nil {
		tmpPath := fmt.Sprintf("%s.tmp.%d", file.StoragePath, time.Now().UnixNano())
		size, err := f.encryptToFile(fileReader, tmpPath)
		if err != nil {
			_ = os.Remove(tmpPath)
			return nil, err
		}
		if err := os.Rename(tmpPath, file.StoragePath); err != nil {
			_ = os.Remove(tmpPath)
			return nil, err
		}
		file.Size = size
		if filename != nil && *filename != "" {
			file.Filename = *filename
		}
	}

	if description != nil {
		file.Description = *description
	}

	if err := f.db.Save(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (f *FileService) DeleteFile(id uint) error {
	var file model.File
	if err := f.db.First(&file, id).Error; err != nil {
		return err
	}

	if err := os.Remove(file.StoragePath); err != nil && !os.IsNotExist(err) {
		// ignore non-existent
		return err
	}
	return f.db.Delete(&file).Error
}

func (f *FileService) GetFileByID(id uint) (*model.File, error) {
	var file model.File
	if err := f.db.First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *FileService) ListFiles(page, size int) (total int64, files []model.File, err error) {
	offset := (page - 1) * size
	if err = s.db.Model(&model.File{}).Count(&total).Error; err != nil {
		return
	}
	err = s.db.Order("created_at desc").Limit(size).Offset(offset).Find(&files).Error
	return
}

func deriveKeys(base64Key string) ([]byte, []byte) {
	raw, err := base64.RawURLEncoding.DecodeString(base64Key)
	if err != nil || len(raw) < 32 {
		return nil, nil
	}
	encKey := hmac.New(sha256.New, raw)
	encKey.Write([]byte("enc"))
	macKey := hmac.New(sha256.New, raw)
	macKey.Write([]byte("mac"))
	return encKey.Sum(nil), macKey.Sum(nil)
}

func (f *FileService) encryptToFile(src io.Reader, dstPath string) (int64, error) {
	if len(f.encKey) != 32 || len(f.macKey) != 32 {
		return 0, errors.New("file crypto key not configured")
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	iv := make([]byte, 16)
	if _, err := rand.Read(iv); err != nil {
		return 0, err
	}

	header := make([]byte, 0, headerLen)
	header = append(header, []byte(fileMagic)...)
	header = append(header, iv...)
	if _, err := out.Write(header); err != nil {
		return 0, err
	}

	block, err := aes.NewCipher(f.encKey)
	if err != nil {
		return 0, err
	}
	stream := cipher.NewCTR(block, iv)

	h := hmac.New(sha256.New, f.macKey)
	if _, err := h.Write(header); err != nil {
		return 0, err
	}

	buf := make([]byte, streamBuf)
	var total int64
	for {
		n, rerr := src.Read(buf)
		if n > 0 {
			total += int64(n)
			stream.XORKeyStream(buf[:n], buf[:n])
			if _, err := out.Write(buf[:n]); err != nil {
				return 0, err
			}
			if _, err := h.Write(buf[:n]); err != nil {
				return 0, err
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			return 0, rerr
		}
	}

	mac := h.Sum(nil)
	if _, err := out.Write(mac); err != nil {
		return 0, err
	}

	return total, nil
}

func (f *FileService) DecryptToWriter(w io.Writer, srcPath string) error {
	if len(f.encKey) != 32 || len(f.macKey) != 32 {
		return errors.New("file crypto key not configured")
	}

	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}
	if info.Size() < headerLen+macLen {
		return errors.New("encrypted file is too small")
	}

	header := make([]byte, headerLen)
	if _, err := io.ReadFull(in, header); err != nil {
		return err
	}
	if string(header[:4]) != fileMagic {
		return errors.New("invalid file magic")
	}
	iv := header[4:]

	expectedMac := make([]byte, macLen)
	if _, err := in.ReadAt(expectedMac, info.Size()-macLen); err != nil {
		return err
	}

	dataLen := info.Size() - macLen
	h := hmac.New(sha256.New, f.macKey)
	section := io.NewSectionReader(in, 0, dataLen)
	if _, err := io.CopyBuffer(h, section, make([]byte, streamBuf)); err != nil {
		return err
	}
	if !hmac.Equal(expectedMac, h.Sum(nil)) {
		return errors.New("file integrity check failed")
	}

	block, err := aes.NewCipher(f.encKey)
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)
	cipherSection := io.NewSectionReader(in, int64(headerLen), dataLen-int64(headerLen))
	buf := make([]byte, streamBuf)
	for {
		n, rerr := cipherSection.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf[:n], buf[:n])
			if _, err := w.Write(buf[:n]); err != nil {
				return err
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			return rerr
		}
	}
	return nil
}
