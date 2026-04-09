package archiver

import (
	"errors"

	"a-rc/internal/core"
)

// Adapter implements core.Archiver.
// Phase 2: create a zip/tar.gz of job.Path inside os.MkdirTemp, return the file path.
// The caller (ArchiveService) is responsible for removing the temp file after upload.
type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) Archive(job core.Job) (string, error) {
	return "", errors.New("archiver: not implemented (Phase 2)")
}
