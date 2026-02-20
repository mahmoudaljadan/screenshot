package app

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type PreferenceStore struct {
	mu   sync.Mutex
	path string
	data map[string]string
}

func NewPreferenceStore(path string) (*PreferenceStore, error) {
	if path == "" {
		cfgDir, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(cfgDir, "go-wails-shot", "preferences.json")
	}
	store := &PreferenceStore{path: path, data: map[string]string{}}
	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

func (p *PreferenceStore) Set(key, value string) error {
	if key == "" {
		return errors.New("preference key is required")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data[key] = value
	return p.saveLocked()
}

func (p *PreferenceStore) Get(key string) string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.data[key]
}

func (p *PreferenceStore) load() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	b, err := os.ReadFile(p.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return json.Unmarshal(b, &p.data)
}

func (p *PreferenceStore) saveLocked() error {
	if err := os.MkdirAll(filepath.Dir(p.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(p.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.path, b, 0o644)
}
