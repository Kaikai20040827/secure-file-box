package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Kaikai20040827/graduation/internal/model"
	"gorm.io/gorm"
)

type FileService struct {
	db       *gorm.DB
	dirpath string
}

func NewFileService(db *gorm.DB, storagePath string) *FileService {
	_ = os.MkdirAll(storagePath, 0755)
	fmt.Println("✓ Creating a new file service done")
	return &FileService{db: db, dirpath: storagePath}
}

func (f *FileService) UploadFile(fileReader io.Reader, filename string, uploaderID uint, description string) (*model.File, error) {
	
	//确保名称在电脑里面唯一
	dst := filepath.Join(f.dirpath, fmt.Sprintf("%d_%s", uploaderID, filename))
	newFile, err := os.Create(dst)
	if err != nil {
		return nil, err
	}
	defer newFile.Close()

	//获取文件大小
	size, err := io.Copy(newFile, fileReader)
	if err != nil {
		return nil, err
	}

	file := &model.File{
		Filename:filename,
		StoragePath: dst,
		Size: size,
		Description: description,
		CreatedAt: time.Now(),
	}

	if err := f.db.Create(f).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func (f *FileService) DeleteFile(id uint) error {
	var file model.File
	if err := f.db.First(&f, id).Error; err != nil {
		return err
	}

	if err := os.Remove(file.StoragePath); err != nil && !os.IsNotExist(err) {
		// ignore non-existent
		return err
	}
	return f.db.Delete(&f).Error
}

func (f *FileService) GetFileByID(id uint) (*model.File, error) {
	var file model.File
	if err := f.db.First(&f, id).Error; err != nil {
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
