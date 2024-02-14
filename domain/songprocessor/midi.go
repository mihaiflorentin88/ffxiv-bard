package songprocessor

import (
	"errors"
	"ffxvi-bard/domain/song"
	"ffxvi-bard/port/contract"
	"fmt"
	"path/filepath"
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

func (m MidiProcessor) getUnprocessedFilePath(currentSong contract.SongInterface) string {
	return filepath.Join(m.UnprocessedPath, currentSong.GetFileCode()+".mid")
}

func (m MidiProcessor) getProcessedFilePath(currentSong contract.SongInterface) string {
	return filepath.Join(m.ProcessedPath, currentSong.GetFileCode()+".mid")
}

func (m MidiProcessor) ProcessSong(currentSong contract.SongInterface) error {
	currentSong.ChangeStatus(song.Processing, "Processing MIDI file")
	file, err := m.filesystem.ReadFile(m.getUnprocessedFilePath(currentSong))
	if err != nil {
		msg := "Error reading file"
		currentSong.ChangeStatus(song.Failed, msg)
		_ = m.RemoveUnprocessedSong(currentSong)
		return errors.New(fmt.Sprintf("%s: %s", msg, err.Error()))
	}
	m.File = file
	if !m.IsCorrectFormat() {
		msg := "file is not a MIDI file"
		currentSong.ChangeStatus(song.Failed, msg)
		_ = m.RemoveUnprocessedSong(currentSong)
		return errors.New(msg)
	}

	err = m.filesystem.WriteFile(m.getProcessedFilePath(currentSong), m.File)
	if err != nil {
		msg := "Error writing file"
		currentSong.ChangeStatus(song.Failed, msg)
		_ = m.RemoveUnprocessedSong(currentSong)
		return errors.New(fmt.Sprintf("%s: %s", msg, err.Error()))
	}
	currentSong.ChangeStatus(song.Processed, "MIDI file processed")
	_ = m.RemoveUnprocessedSong(currentSong)
	return nil
}

func (m MidiProcessor) WriteUnprocessedSong(currentSong contract.SongInterface) error {
	err := m.filesystem.WriteFile(m.getUnprocessedFilePath(currentSong), currentSong.GetFile())
	if err != nil {
		msg := "Error writing file"
		currentSong.ChangeStatus(song.Failed, msg)
		return errors.New(fmt.Sprintf("%s: %s", msg, err.Error()))
	}
	return nil
}

func (m MidiProcessor) RemoveUnprocessedSong(currentSong contract.SongInterface) error {
	err := m.filesystem.RemoveFile(m.getUnprocessedFilePath(currentSong))
	if err != nil {
		msg := "Error removing file"
		currentSong.ChangeStatus(song.Failed, msg)
		return errors.New(fmt.Sprintf("%s: %s", msg, err.Error()))
	}
	return nil
}
