package usecase

import (
	"bytes"
	"compiler-playground-api/internal/entity"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type CodeRepository interface {
	SaveCodeMetadata(id, language string, createdAt string) error
}

type MinioService interface {
	SaveFile(id, sourceCode string, userInput string) error
	GetFile(id string) (string, string, error)
}

type CodeUseCase struct {
	repo         CodeRepository
	minioService MinioService
}

func NewCodeUseCase(repo CodeRepository, minioService MinioService) *CodeUseCase {
	return &CodeUseCase{
		repo:         repo,
		minioService: minioService,
	}
}

func (uc *CodeUseCase) SaveCode(code *entity.Code) (string, error) {
	if err := uc.minioService.SaveFile(code.ID, code.SourceCode, code.Input); err != nil {
		return "", err
	}

	if err := uc.repo.SaveCodeMetadata(code.ID, code.Language, code.CreatedAt.String()); err != nil {
		return "", err
	}

	return code.ID, nil
}

func (uc *CodeUseCase) ExecuteCode(id string) (string, error) {
	sourceCode, userInput, err := uc.minioService.GetFile(id)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", "python-*.py")
	if err != nil {
		return "", errors.New("failed to create temp file")
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(sourceCode); err != nil {
		return "", errors.New("failed to write source code to temp file")
	}
	if err := tmpFile.Close(); err != nil {
		return "", errors.New("failed to close temp file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "python3", tmpFile.Name())
	cmd.Stdin = bytes.NewBufferString(userInput)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}

	fmt.Println(out.String())

	output := out.String()
	if !strings.HasSuffix(output, "\n") {
		output += "\n"
	}

	return output, nil
}
