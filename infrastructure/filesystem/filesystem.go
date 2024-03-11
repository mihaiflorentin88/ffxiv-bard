package filesystem

import (
	"os"
)

type FileSystem struct{}

func NewFileSystem() *FileSystem {
	return &FileSystem{}
}

func (lf FileSystem) EnsureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm) // os.ModePerm is 0777, allowing read, write, and execute
		if err != nil {
			return err
		}
	}
	return nil
}

func (lf FileSystem) ListFiles(directory string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func (lf FileSystem) RemoveFile(filepath string) error {
	return os.Remove(filepath)
}

func (lf FileSystem) WriteFile(filepath string, data []byte) error {
	return os.WriteFile(filepath, data, 0644)
}

func (lf FileSystem) ReadFile(filepath string) ([]byte, error) {
	return os.ReadFile(filepath)
}
