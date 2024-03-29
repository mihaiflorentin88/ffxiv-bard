package processor

import (
	"errors"
	"ffxvi-bard/port/contract"
	"fmt"
	"github.com/google/uuid"
	"path/filepath"
	"strings"
)

type MidiProcessor struct {
	UnprocessedPath string
	ProcessedPath   string
	File            []byte
	filesystem      contract.FileSystemInterface
}

func NewMidiProcessor(UnprocessedPath string, FilepathProcessed string, filesystem contract.FileSystemInterface) *MidiProcessor {
	return &MidiProcessor{
		UnprocessedPath: UnprocessedPath,
		ProcessedPath:   FilepathProcessed,
		filesystem:      filesystem}
}

func (m MidiProcessor) IsCorrectFormat() bool {
	if len(m.File) < 4 {
		return false
	}
	isMidi := string(m.File[:4]) == "MThd"
	return isMidi
}

func (m MidiProcessor) GetUnprocessedFilePath(songFilename string) string {
	if !strings.HasSuffix(songFilename, ".mid") {
		songFilename = songFilename + ".mid"
	}
	return filepath.Join(m.UnprocessedPath, songFilename)
}

func (m MidiProcessor) GetProcessedFilePath(songFilename string) string {
	if !strings.HasSuffix(songFilename, ".mid") {
		songFilename = songFilename + ".mid"
	}
	return filepath.Join(m.ProcessedPath, songFilename)
}

func (m MidiProcessor) ProcessSong(songFilename string) error {
	file, err := m.filesystem.ReadFile(m.GetUnprocessedFilePath(songFilename))
	if err != nil {
		msg := "Error reading file"
		_ = m.RemoveUnprocessedSong(songFilename)
		return errors.New(fmt.Sprintf("%s: %s", msg, err.Error()))
	}
	m.File = file
	if !m.IsCorrectFormat() {
		msg := "file is not a MIDI file"
		_ = m.RemoveUnprocessedSong(songFilename)
		return errors.New(msg)
	}

	err = m.filesystem.WriteFile(m.GetProcessedFilePath(songFilename), m.File)
	if err != nil {
		msg := "Error writing file"
		return errors.New(fmt.Sprintf("%s: %s", msg, err.Error()))
	}
	_ = m.RemoveUnprocessedSong(songFilename)
	return nil
}

func (m MidiProcessor) WriteUnprocessedSong(songFilename string, song []byte) error {
	err := m.filesystem.WriteFile(m.GetUnprocessedFilePath(songFilename), song)
	if err != nil {
		msg := "Error writing file"
		return errors.New(fmt.Sprintf("%s: %s", msg, err.Error()))
	}
	return nil
}

func (m MidiProcessor) RemoveUnprocessedSong(songFilename string) error {
	err := m.filesystem.RemoveFile(m.GetUnprocessedFilePath(songFilename))
	if err != nil {
		msg := "Error removing file"
		return errors.New(fmt.Sprintf("%s: %s", msg, err.Error()))
	}
	return nil
}

func (m MidiProcessor) GenerateUniqueFilename() string {
	return uuid.New().String() + ".mid"
}
