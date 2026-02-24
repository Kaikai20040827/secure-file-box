package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Kaikai20040827/graduation/internal/model"
	"gorm.io/gorm"
)

type FileService struct {
	db        *gorm.DB
	dirpath   string
	cipherKey []byte
}

func NewFileService(db *gorm.DB, storagePath, encryptSecret string) *FileService {
	_ = os.MkdirAll(storagePath, 0755)
	fmt.Println("✓ Creating a new file service done")
	key := sha256.Sum256([]byte(encryptSecret))
	return &FileService{db: db, dirpath: storagePath, cipherKey: key[:]}
}

func (f *FileService) UploadFile(fileReader io.Reader, filename string, uploaderID uint, description string) (*model.File, error) {
	dst := filepath.Join(f.dirpath, fmt.Sprintf("%d_%s", uploaderID, filename))
	plainBytes, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}

	encryptedBytes, err := encryptBytes(plainBytes, f.cipherKey)
	if err != nil {
		return nil, err
	}

	if err = os.WriteFile(dst, encryptedBytes, 0644); err != nil {
		return nil, err
	}

	file := &model.File{
		Filename:    filename,
		StoragePath: dst,
		Size:        int64(len(plainBytes)),
		Description: description,
		UploaderID:  fmt.Sprintf("%d", uploaderID),
		CreatedAt:   time.Now(),
	}

	if err := f.db.Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func (f *FileService) DeleteFile(id uint) error {
	var file model.File
	if err := f.db.First(&file, id).Error; err != nil {
		return err
	}

	if err := os.Remove(file.StoragePath); err != nil && !os.IsNotExist(err) {
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

func (f *FileService) ReadDecryptedFile(id uint) (*model.File, []byte, error) {
	file, err := f.GetFileByID(id)
	if err != nil {
		return nil, nil, err
	}
	encryptedBytes, err := os.ReadFile(file.StoragePath)
	if err != nil {
		return nil, nil, err
	}
	plainBytes, err := decryptBytes(encryptedBytes, f.cipherKey)
	if err != nil {
		return nil, nil, err
	}
	return file, plainBytes, nil
}

func (s *FileService) ListFiles(page, size int) (total int64, files []model.File, err error) {
	offset := (page - 1) * size
	if err = s.db.Model(&model.File{}).Count(&total).Error; err != nil {
		return
	}
	err = s.db.Order("created_at desc").Limit(size).Offset(offset).Find(&files).Error
	return
}

func encryptBytes(plainData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	cipherBytes := gcm.Seal(nil, nonce, plainData, nil)
	encoded := make([]byte, hex.EncodedLen(len(nonce)+len(cipherBytes)))
	hex.Encode(encoded, append(nonce, cipherBytes...))
	return encoded, nil
}

func decryptBytes(encryptedData, key []byte) ([]byte, error) {
	raw := make([]byte, hex.DecodedLen(len(encryptedData)))
	n, err := hex.Decode(raw, encryptedData)
	if err != nil {
		return nil, err
	}
	raw = raw[:n]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(raw) < nonceSize {
		return nil, fmt.Errorf("invalid encrypted file data")
	}
	nonce, cipherBytes := raw[:nonceSize], raw[nonceSize:]
	return gcm.Open(nil, nonce, cipherBytes, nil)
}
