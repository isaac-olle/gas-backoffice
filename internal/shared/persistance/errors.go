package persistance

import "errors"

// ErrConcurrentModification is returned by a Repository.Save implementation
// when an optimistic-concurrency check (a version column) detects that the
// aggregate was modified by another writer since it was loaded. Infrastructure
// returns this raw sentinel per CLAUDE.md convention; the application service
// layer is responsible for mapping it to an app_errors.ConflictError.
var ErrConcurrentModification = errors.New("persistance: concurrent modification detected")
