package manifest

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/codemint/codemint-cli/internal/util"
)

const CurrentVersion = "1"

type Item struct {
	CatalogID   string    `json:"catalogId"`
	Ref         string    `json:"ref"`
	Type        string    `json:"type"`
	Slug        string    `json:"slug"`
	Tool        string    `json:"tool,omitempty"`
	Version     string    `json:"version"`
	Checksum    string    `json:"checksum"`
	InstalledAt time.Time `json:"installedAt"`
	Path        string    `json:"path"`
}

type File struct {
	Version   string `json:"version"`
	Installed []Item `json:"installed"`
}

type Store struct {
	Root string
}

type Settings struct {
	AITool string `json:"aiTool,omitempty"`
}

func New(root string) *Store {
	return &Store{Root: root}
}

func (s *Store) BaseDir() string {
	return filepath.Join(s.Root, ".codemint")
}

func (s *Store) Path() string {
	return filepath.Join(s.BaseDir(), "manifest.json")
}

func (s *Store) SettingsPath() string {
	return filepath.Join(s.BaseDir(), "settings.json")
}

func (s *Store) Load() (File, error) {
	path := s.Path()
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return File{Version: CurrentVersion, Installed: []Item{}}, nil
	}
	if err != nil {
		return File{}, err
	}
	var mf File
	if err := json.Unmarshal(b, &mf); err != nil {
		return File{}, err
	}
	if mf.Version == "" {
		mf.Version = CurrentVersion
	}
	if mf.Installed == nil {
		mf.Installed = []Item{}
	}
	return mf, nil
}

func (s *Store) Save(mf File) error {
	if mf.Version == "" {
		mf.Version = CurrentVersion
	}
	sort.Slice(mf.Installed, func(i, j int) bool {
		if mf.Installed[i].Type == mf.Installed[j].Type {
			return mf.Installed[i].Slug < mf.Installed[j].Slug
		}
		return mf.Installed[i].Type < mf.Installed[j].Type
	})
	b, err := json.MarshalIndent(mf, "", "  ")
	if err != nil {
		return err
	}
	return util.AtomicWriteFile(s.Path(), append(b, '\n'), 0o644)
}

func (s *Store) LoadSettings() (Settings, error) {
	b, err := os.ReadFile(s.SettingsPath())
	if errors.Is(err, os.ErrNotExist) {
		return Settings{}, nil
	}
	if err != nil {
		return Settings{}, err
	}
	var settings Settings
	if err := json.Unmarshal(b, &settings); err != nil {
		return Settings{}, err
	}
	return settings, nil
}

func (s *Store) SaveSettings(settings Settings) error {
	b, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	return util.AtomicWriteFile(s.SettingsPath(), append(b, '\n'), 0o644)
}

func FindByCatalogID(items []Item, catalogID string) (int, bool) {
	for i, it := range items {
		if it.CatalogID == catalogID {
			return i, true
		}
	}
	return -1, false
}

func FindByRef(items []Item, ref string) (int, bool) {
	for i, it := range items {
		if it.Ref == ref {
			return i, true
		}
	}
	return -1, false
}
