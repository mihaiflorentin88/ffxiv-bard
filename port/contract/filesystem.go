package contract

type FileSystemInterface interface {
	EnsureDir(dirPath string) error
	ListFiles(directory string) ([]string, error)
	RemoveFile(filepath string) error
	WriteFile(filepath string, data []byte) error
	ReadFile(filepath string) ([]byte, error)
}
