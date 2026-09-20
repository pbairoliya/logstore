package logstore

import (
	"bytes"
	"path/filepath"
	"testing"
)

// The first test to make pass. Everything else builds on it.
func TestPutThenGet(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "log"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	if err := s.Put("name", []byte("pratik")); err != nil {
		t.Fatalf("Put: %v", err)
	}
	got, err := s.Get("name")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !bytes.Equal(got, []byte("pratik")) {
		t.Fatalf("Get = %q, want %q", got, "pratik")
	}
}

func TestGetMissingKey(t *testing.T) {
	s, _ := Open(filepath.Join(t.TempDir(), "log"))
	defer s.Close()

	if _, err := s.Get("absent"); err != ErrKeyNotFound {
		t.Fatalf("Get missing = %v, want ErrKeyNotFound", err)
	}
}

func TestOverwriteReturnsNewestValue(t *testing.T) {
	s, _ := Open(filepath.Join(t.TempDir(), "log"))
	defer s.Close()

	s.Put("k", []byte("first"))
	s.Put("k", []byte("second"))

	got, err := s.Get("k")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !bytes.Equal(got, []byte("second")) {
		t.Fatalf("Get = %q, want the newest write", got)
	}
}

// Reopen is where a log-structured store earns its name: the index is
// reconstructed purely by replaying the log.
func TestSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log")

	s, _ := Open(path)
	s.Put("persisted", []byte("yes"))
	s.Close()

	reopened, err := Open(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer reopened.Close()

	got, err := reopened.Get("persisted")
	if err != nil {
		t.Fatalf("Get after reopen: %v", err)
	}
	if !bytes.Equal(got, []byte("yes")) {
		t.Fatalf("Get after reopen = %q, want %q", got, "yes")
	}
}

func TestDeleteSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log")

	s, _ := Open(path)
	s.Put("doomed", []byte("v"))
	s.Delete("doomed")
	s.Close()

	reopened, _ := Open(path)
	defer reopened.Close()

	if _, err := reopened.Get("doomed"); err != ErrKeyNotFound {
		t.Fatalf("deleted key came back after reopen: %v", err)
	}
}

// Compaction must be invisible from the outside: same reads, smaller file.
func TestCompactPreservesLiveData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log")
	s, _ := Open(path)
	defer s.Close()

	s.Put("keep", []byte("final"))
	s.Put("keep", []byte("final-2"))
	s.Put("drop", []byte("x"))
	s.Delete("drop")

	if err := s.Compact(); err != nil {
		t.Fatalf("Compact: %v", err)
	}

	got, err := s.Get("keep")
	if err != nil || !bytes.Equal(got, []byte("final-2")) {
		t.Fatalf("after compact Get(keep) = %q, %v", got, err)
	}
	if _, err := s.Get("drop"); err != ErrKeyNotFound {
		t.Fatalf("compaction resurrected a deleted key")
	}
}
