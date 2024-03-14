package config

import (
	"github.com/joho/godotenv"
)

type FileConfigInstance struct {
}

func (f *FileConfigInstance) GetContent() (map[string]string, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	return godotenv.Read()
}
