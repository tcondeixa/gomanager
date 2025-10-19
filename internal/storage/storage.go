package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Storage[T any] struct {
	filePath  string
	updatedAt time.Time
	binaries  map[string]T
}

func New[T any](filePath string) *Storage[T] {
	return &Storage[T]{
		filePath:  filePath,
		updatedAt: time.Now(),
		binaries:  map[string]T{},
	}
}

func (s *Storage[T]) Load() error {
	err := s.ensureFile()
	if err != nil {
		return fmt.Errorf("failed to ensure storage file: %w", err)
	}

	err = s.loadFile()
	if err != nil {
		return fmt.Errorf("failed to load storage file: %w", err)
	}

	return nil
}

func (s *Storage[T]) Save() error {
	return s.save(s.filePath, os.O_WRONLY)
}

func (s *Storage[T]) Dump(file string) error {
	return s.save(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
}

func (s *Storage[T]) save(file string, flag int) error {
	fp, err := os.OpenFile(file, flag, 0o644)
	if err != nil {
		return err
	}
	defer fp.Close()

	err = json.NewEncoder(fp).Encode(s.binaries)
	if err != nil {
		return err
	}

	err = fp.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage[T]) ensureFile() error {
	_, err := os.Stat(s.filePath)
	if err == nil {
		return nil
	}

	fp, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer fp.Close()

	err = json.NewEncoder(fp).Encode(map[string]T{})
	if err != nil {
		return err
	}

	err = fp.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage[T]) loadFile() error {
	fp, err := os.OpenFile(s.filePath, os.O_RDONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open storage file: %w", err)
	}
	defer fp.Close()

	err = json.NewDecoder(fp).Decode(&s.binaries)
	if err != nil {
		return fmt.Errorf("failed to decode storage file: %w", err)
	}

	return nil
}

func (s *Storage[T]) SaveItem(key string, item T) {
	s.binaries[key] = item
}

func (s *Storage[T]) DeleteItem(key string) {
	delete(s.binaries, key)
}

func (s *Storage[T]) GetItem(key string) (T, bool) {
	val, ok := s.binaries[key]
	return val, ok
}

func (s *Storage[T]) GetAllItems() map[string]T {
	return s.binaries
}
