// Package logstore is a log-structured key-value store.
//
// The design, in one paragraph: writes append to a single log file and never
// modify what is already there. An in-memory index maps each key to the byte
// offset of its most recent record, so a read is one map lookup and one seek.
// Deletes append a tombstone rather than removing anything. The log therefore
// grows forever until compaction rewrites it, keeping only the records the
// index still points at.
//
// Everything below is deliberately unimplemented. The tests describe the
// behaviour; making them pass is the exercise.
package logstore

import "errors"

// ErrKeyNotFound is returned by Get for a key that was never written, or that
// was written and later deleted.
var ErrKeyNotFound = errors.New("logstore: key not found")

// Store is a log-structured key-value store backed by a single append-only file.
type Store struct {
	// path is the log file on disk.
	path string

	// index maps a key to the offset of its most recent record in the log.
	// Rebuilt by scanning the log at Open.
	index map[string]int64

	// TODO: the file handle, a write offset, and a mutex once concurrency matters.
}

// Open opens (or creates) a store at path and rebuilds the in-memory index by
// scanning the log from the beginning.
//
// The scan is the interesting part: every record is read in order, and later
// records for the same key overwrite earlier ones in the index. That is what
// makes crash recovery free — the log *is* the recovery mechanism.
func Open(path string) (*Store, error) {
	// TODO
	return nil, errors.New("not implemented")
}

// Put appends a key/value record to the log and updates the index.
//
// The write must be durable before Put returns, or a crash immediately after
// will lose an acknowledged write. Decide explicitly whether that means fsync
// per write (slow, correct) or a batched sync (fast, a bounded window of loss)
// — and write the decision down in the README.
func (s *Store) Put(key string, value []byte) error {
	// TODO
	return errors.New("not implemented")
}

// Get returns the most recent value for key, or ErrKeyNotFound.
func (s *Store) Get(key string) ([]byte, error) {
	// TODO
	return nil, errors.New("not implemented")
}

// Delete appends a tombstone record. The key must then read as ErrKeyNotFound,
// and the tombstone must survive a reopen — a delete that a restart undoes is
// worse than no delete at all.
func (s *Store) Delete(key string) error {
	// TODO
	return errors.New("not implemented")
}

// Close flushes and closes the underlying file.
func (s *Store) Close() error {
	// TODO
	return errors.New("not implemented")
}

// Compact rewrites the log keeping only live records — the newest record per
// key, and no tombstones — then atomically swaps it into place.
//
// Atomically is the whole trick: write to a temp file, fsync it, rename over
// the original. A crash at any point must leave either the old log or the new
// one, never a half-written file.
func (s *Store) Compact() error {
	// TODO
	return errors.New("not implemented")
}
