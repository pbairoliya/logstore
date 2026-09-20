# logstore

A log-structured key-value store in Go, written to answer a question I couldn't answer properly:

> **What is a database actually doing when I call `put`?**

I work on batch systems that move a lot of data through storage engines I didn't write. This is
me writing the smallest honest version of one.

**Status: scaffold.** The API and the test suite exist; the implementation doesn't. Tests fail on
purpose — the point is to make them pass one at a time.

## The design

Writes append to a single file and never modify what's already there. An in-memory index maps
each key to the byte offset of its newest record, so a read is a map lookup and a seek. Deletes
append a tombstone. The log grows until compaction rewrites it, keeping only what the index still
points at.

```
Put("k","v")  ──► append record ──► index["k"] = offset
Get("k")      ──► index lookup  ──► seek + read
Delete("k")   ──► append tombstone
Compact()     ──► rewrite live records to temp ──► fsync ──► rename
```

Three properties make it interesting rather than trivial:

1. **The log is the recovery mechanism.** There is no separate WAL yet, because the log already
   is one. Reopening replays it from byte zero to rebuild the index.
2. **Durability is a decision, not a default.** fsync per write is correct and slow; batching is
   fast and loses a bounded window. Whichever I pick has to be written down and defended.
3. **Compaction has to be atomic.** Temp file, fsync, rename. A crash mid-compaction must leave
   either the old log or the new one, never a torn file.

## Order to build it

| # | Test | What it forces you to decide |
|---|---|---|
| 1 | `TestPutThenGet` | The record format — key length, value length, checksum? |
| 2 | `TestGetMissingKey` | Error semantics |
| 3 | `TestOverwriteReturnsNewestValue` | That the index points at the newest offset |
| 4 | `TestSurvivesReopen` | The replay scan — this is the real one |
| 5 | `TestDeleteSurvivesReopen` | Tombstones, and that replay honours them |
| 6 | `TestCompactPreservesLiveData` | Atomic rewrite |

```bash
make test
```

## After the tests pass

- A crash test: kill the process mid-write in a loop, reopen, assert no acknowledged write was lost.
- Benchmarks: writes/sec at fsync-per-write vs batched, read latency as the log grows.
- Then a real WAL, and then hint files so reopen doesn't scan the whole log.
- Then swap the index for a B-tree and measure what that costs and buys.

## Prior art worth reading

Bitcask (Riak) is the canonical version of this design, and the
[Bitcask paper](https://riak.com/assets/bitcask-intro.pdf) is nine pages and very readable.
