// Package lock provides a cross-platform exclusive process lock backed by a lock
// file. The daemon acquires it at startup so a second instance fails fast rather
// than starting a second scheduler (which would double-execute every task). The
// OS releases the lock automatically if the process dies, so a crashed daemon
// does not leave a stale lock.
//
// Acquire and the Lock type are defined per-platform (flock on Unix, LockFileEx
// on Windows).
package lock
