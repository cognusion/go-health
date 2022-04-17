package health

import (
	"errors"
	"sync"
)

var (
	// ErrNoSuchEntryError is returned when the requested element does not exist in the Registry
	ErrNoSuchEntryError = errors.New("no such element exists")
)

// StatusRegistry is a gorosafe map of services to their Status objects
type StatusRegistry struct {
	sync.RWMutex
	stats map[string]Status
}

// NewStatusRegistry returns an initialized StatusRegistry
func NewStatusRegistry() *StatusRegistry {
	return &StatusRegistry{
		stats: make(map[string]Status),
	}
}

// Add or update an entry in StatusRegistry
func (s *StatusRegistry) Add(name, status string, Value, ExpectedValue interface{}) {
	sname := SafeLabel(name)
	s.Lock()
	stat := Status{
		Name:          sname,
		Status:        status,
		Value:         Value,
		ExpectedValue: ExpectedValue,
	}
	s.stats[name] = stat
	s.Unlock()
}

// Remove an entry from the StatusRegistry
func (s *StatusRegistry) Remove(name string) {
	s.Lock()
	delete(s.stats, name)
	s.Unlock()
}

// Keys returns a list of names from the StatusRegistry
func (s *StatusRegistry) Keys() []string {
	keys := make([]string, len(s.stats))
	i := 0
	s.RLock()
	for k := range s.stats {
		keys[i] = k
		i++
	}
	s.RUnlock()
	return keys
}

// Get returns the requested Status, or ErrNoSuchEntryError
func (s *StatusRegistry) Get(name string) (*Status, error) {
	s.RLock()
	defer s.RUnlock()

	if stat, ok := s.stats[name]; ok {
		return &stat, nil
	}
	return nil, ErrNoSuchEntryError
}
