package mocks

type FilesystemMock struct{}

func (lf FilesystemMock) EnsureDir(dirPath string) error {
	return nil
}

func (lf FilesystemMock) ListFiles(directory string) ([]string, error) {
	files := []string{"file1", "file2", "file3"}
	return files, nil
}

func (lf FilesystemMock) RemoveFile(filepath string) error {
	return nil
}

func (lf FilesystemMock) WriteFile(filepath string, data []byte) error {
	return nil
}

func (lf FilesystemMock) ReadFile(filepath string) ([]byte, error) {
	header := []byte{
		0x4D, 0x54, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x06,
		0x00, 0x01,
		0x00, 0x01,
		0x00, 0x60,
	}

	trackStart := []byte{
		0x4D, 0x54, 0x72, 0x6B,
		0x00, 0x00, 0x00, 0x04,
	}

	trackEnd := []byte{
		0x00,
		0xFF, 0x2F,
		0x00,
	}
	return append(append(header, trackStart...), trackEnd...), nil
}
