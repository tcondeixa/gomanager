package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const version = "v1"

type File[T any] struct {
	Version   string       `json:"version"`
	UpdatedAt time.Time    `json:"updated_at"`
	Binaries  map[string]T `json:"binaries"`
}

func NewFile[T any]() File[T] {
	return File[T]{
		UpdatedAt: time.Now(),
		Version:   version,
		Binaries:  map[string]T{},
	}
}

type Provider[T any] struct {
	filePath   string
	fileFormat File[T]
}

func New[T any](filePath string) *Provider[T] {
	return &Provider[T]{
		filePath:   filePath,
		fileFormat: NewFile[T](),
	}
}

func (s *Provider[T]) Start() error {
	err := s.ensureFile()
	if err != nil {
		return fmt.Errorf("failed to ensure storage file: %w", err)
	}

	err = s.loadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to load storage file: %w", err)
	}

	return nil
}

func (s *Provider[T]) Import(file string) error {
	err := s.loadFile(file)
	if err != nil {
		return fmt.Errorf("failed to import storage file: %w", err)
	}

	return s.saveFile(s.filePath)
}

func (s *Provider[T]) Export(file string) error {
	return s.saveFile(file)
}

func (s *Provider[T]) saveFile(file string) error {
	bytes, err := json.MarshalIndent(s.fileFormat, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode items to json: %w", err)
	}

	err = os.WriteFile(file, bytes, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write items to file: %w", err)
	}

	return nil
}

func (s *Provider[T]) ensureFile() error {
	_, err := os.Stat(s.filePath)
	if err == nil {
		return nil
	}

	err = s.saveFile(s.filePath)
	if err != nil {
		return err
	}

	return nil
}

func (s *Provider[T]) loadFile(file string) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read storage file: %w", err)
	}
	err = json.Unmarshal(bytes, &s.fileFormat)
	if err != nil {
		return fmt.Errorf("failed to unmarshal storage file: %w", err)
	}

	return nil
}

func (s *Provider[T]) SaveItem(key string, item T) error {
	s.fileFormat.Binaries[key] = item
	s.fileFormat.UpdatedAt = time.Now()

	return s.saveFile(s.filePath)
}

func (s *Provider[T]) DeleteItem(key string) error {
	delete(s.fileFormat.Binaries, key)
	s.fileFormat.UpdatedAt = time.Now()

	return s.saveFile(s.filePath)
}

func (s *Provider[T]) GetItem(key string) (T, bool) {
	val, ok := s.fileFormat.Binaries[key]
	return val, ok
}

func (s *Provider[T]) GetAllItems() map[string]T {
	return s.fileFormat.Binaries
}
