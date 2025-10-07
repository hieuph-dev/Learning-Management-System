package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var allowExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}
var allowMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

const maxSize = 5 << 20

func ValidateAndSaveFile(fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	// Check extension in filename
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowExts[ext] {
		return "", errors.New("unsupported file extension")
	}

	// Check size
	if fileHeader.Size > maxSize {
		return "", errors.New("file too large (max 5MB)")
	}

	// Check file type
	file, err := fileHeader.Open()
	if err != nil {
		return "", errors.New("cannot open file")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", errors.New("cannot read file")
	}

	mimeType := http.DetectContentType(buffer)
	if !allowMimeTypes[mimeType] {
		return "", fmt.Errorf("invalid MIME type: %s", mimeType)
	}

	// Change file name
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Create folder if not exist
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", errors.New("cannot create upload folder")
	}

	// uploadDir "./upload" + filename "abc.jpg"
	savePath := filepath.Join(uploadDir, fileName)
	if err := saveFile(fileHeader, savePath); err != nil {
		return "", err
	}

	return fileName, nil
}

func saveFile(fileHeader *multipart.FileHeader, destination string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)

	return err
}

func ValidateAvatarFile(fileHeader *multipart.FileHeader) error {
	// Kiểm tra kích thước file (max 2MB cho avatar)
	maxSize := int64(2 << 20) // 2MB
	if fileHeader.Size > maxSize {
		return errors.New("avatar file too large (max 2MB)")
	}

	// Kiểm tra extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	if !allowedExts[ext] {
		return errors.New("unsupported avatar file type. Only JPG, JPEG, PNG, GIF are allowed")
	}

	// Kiểm tra MIME type
	file, err := fileHeader.Open()
	if err != nil {
		return errors.New("cannot open avatar file")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return errors.New("cannot read avatar file")
	}

	mimeType := http.DetectContentType(buffer)
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	if !allowedMimeTypes[mimeType] {
		return fmt.Errorf("invalid avatar MIME type: %s", mimeType)
	}

	return nil
}
