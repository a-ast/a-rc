package core

import (
	"os"
	"path/filepath"
)

// ArchiveService orchestrates archiving and uploading a single job.
type ArchiveService struct {
	archiver Archiver
	uploader Uploader
}

func NewArchiveService(a Archiver, u Uploader) *ArchiveService {
	return &ArchiveService{archiver: a, uploader: u}
}

// RunJob archives the job's source path, uploads the result to the configured GDrive folder,
// then removes the temporary archive file.
func (s *ArchiveService) RunJob(job Job) error {
	archivePath, err := s.archiver.Archive(job)
	if err != nil {
		return err
	}
	defer os.RemoveAll(filepath.Dir(archivePath)) // remove the entire temp dir

	return s.uploader.Upload(archivePath)
}
