package repository

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// FileRepositorier is the interface that wraps the basic Read method
type FileRepositorier interface {
	Read(ctx context.Context) ([]string, string, error)
}

// FileRepository is the struct that implements the FileRepositorier interface
type FileRepository struct {
	logger *zap.SugaredLogger
}

// NewFileRepository returns a new FileRepository
func NewFileRepository(logger *zap.SugaredLogger) *FileRepository {
	return &FileRepository{logger}
}

// Read returns the data and the hash of the file
func (fr *FileRepository) Read(ctx context.Context) ([]string, string, error) {
	fileName := viper.GetString("FILE_NAME")

	file, err := os.Open(fileName)
	if err != nil {
		return nil, *new(string), err
	}

	defer file.Close()

	hasher := md5.New()
	scanner := bufio.NewScanner(file)
	data := make([]string, 0)

	for scanner.Scan() {
		data = append(data, scanner.Text())
		hasher.Write([]byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		message := fmt.Sprintf("Could not read file %s", fileName)
		return nil, *new(string), errors.New(message)
	}

	return data, hex.EncodeToString(hasher.Sum(nil)), nil
}
