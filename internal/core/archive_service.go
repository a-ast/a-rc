package core

// ArchiveService orchestrates archiving and uploading a single job.
type ArchiveService struct {
	archiver Archiver
	uploader Uploader
}

func NewArchiveService(a Archiver, u Uploader) *ArchiveService {
	return &ArchiveService{archiver: a, uploader: u}
}

// RunJob archives the job's source path and uploads the result to the configured GDrive folder.
func (s *ArchiveService) RunJob(job Job) error {
	archivePath, err := s.archiver.Archive(job)
	if err != nil {
		return err
	}
	return s.uploader.Upload(archivePath)
}
