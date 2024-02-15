package processor

import (
	"ffxvi-bard/mocks"
	"testing"
)

func getFileSystemMock() *mocks.FilesystemMock {
	return &mocks.FilesystemMock{}
}

func getSongsMock() *mocks.MockSongInterface {
	return &mocks.MockSongInterface{}
}

func getNewMidiProcessor() *MidiProcessor {
	return NewMidiProcessor("unprocessed", "processed", getFileSystemMock())
}

// write tests

func TestNewMidiProcessor(t *testing.T) {
	mp := getNewMidiProcessor()
	if mp.UnprocessedPath != "unprocessed" {
		t.Errorf("NewMidiProcessor should set UnprocessedPath to 'unprocessed'")
	}
	if mp.ProcessedPath != "processed" {
		t.Errorf("NewMidiProcessor should set ProcessedPath to 'processed'")
	}
}

func TestMidiProcessorIsCorrectFormat(t *testing.T) {
	mp := getNewMidiProcessor()
	mp.File = []byte("MThd")
	if !mp.IsCorrectFormat() {
		t.Errorf("IsCorrectFormat should return true")
	}
	mp.File = []byte("test")
	if mp.IsCorrectFormat() {
		t.Errorf("IsCorrectFormat should return false")
	}
}

func TestMidiProcessorGetUnprocessedFilePath(t *testing.T) {
	mp := getNewMidiProcessor()
	a := mp.getUnprocessedFilePath(getSongsMock())
	println(a)
	if mp.getUnprocessedFilePath(getSongsMock()) != "unprocessed/test123.mid" {
		t.Errorf("getUnprocessedFilePath should return 'unprocessed/test.mid'")
	}
}

func TestMidiProcessorGetProcessedFilePath(t *testing.T) {
	mp := getNewMidiProcessor()
	if mp.getProcessedFilePath(getSongsMock()) != "processed/test123.mid" {
		t.Errorf("getProcessedFilePath should return 'processed/test.mid'")
	}
}
