package repository

import (
	"compiler-playground-api/internal/domain"
	"os"
)

type tempFileRepo struct{}

func NewTempFileRepository() domain.TempFileRepository {
	return &tempFileRepo{}
}

func (r *tempFileRepo) Create(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "code*.py")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func (r *tempFileRepo) Delete(name string) error {
	return os.Remove(name)
}
