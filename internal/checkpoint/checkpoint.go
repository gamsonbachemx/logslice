// Package checkpoint persists and restores stream positions so that
// logslice can resume from where it left off after a restart.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// State holds the last-seen position for a given URL.
type State struct {
	URL       string    `json:"url"`
	Offset    int64     `json:"offset"`
	Timestamp time.Time `json:"timestamp"`
}

// Store manages checkpoint persistence to a JSON file on disk.
type Store struct {
	mu   sync.Mutex
	path string
	data map[string]State
}

// New opens or creates a checkpoint store at the given file path.
func New(path string) (*Store, error) {
	s := &Store{
		path: path,
		data: make(map[string]State),
	}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Get returns the stored State for a URL, and whether one exists.
func (s *Store) Get(url string) (State, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.data[url]
	return st, ok
}

// Set updates the State for a URL and flushes to disk.
func (s *Store) Set(url string, offset int64, ts time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[url] = State{URL: url, Offset: offset, Timestamp: ts}
	return s.flush()
}

// Delete removes the checkpoint entry for a URL and flushes to disk.
func (s *Store) Delete(url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, url)
	return s.flush()
}

func (s *Store) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.data)
}

func (s *Store) flush() error {
	f, err := os.CreateTemp("", "checkpoint-*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	if err := json.NewEncoder(f).Encode(s.data); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	f.Close()
	return os.Rename(tmp, s.path)
}
