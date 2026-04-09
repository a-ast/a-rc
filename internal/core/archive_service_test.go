package core

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// stubArchiver is a test double for the Archiver port.
type stubArchiver struct {
	path string
	err  error
}

func (s *stubArchiver) Archive(_ Job) (string, error) { return s.path, s.err }

// stubUploader is a test double for the Uploader port.
type stubUploader struct {
	received string
	err      error
}

func (s *stubUploader) Upload(localPath string) error {
	s.received = localPath
	return s.err
}

func TestRunJob_ArchivesAndUploads(t *testing.T) {
	// Create a real temp file so RemoveAll has something to clean up.
	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "test.zip")
	if err := os.WriteFile(zipPath, []byte("zip"), 0o644); err != nil {
		t.Fatal(err)
	}

	arc := &stubArchiver{path: zipPath}
	upl := &stubUploader{}
	svc := NewArchiveService(arc, upl)

	if err := svc.RunJob(Job{Path: "/some/path"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if upl.received != zipPath {
		t.Errorf("uploader received %q, want %q", upl.received, zipPath)
	}
}

func TestRunJob_CleansTempDirAfterUpload(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "a-rc-svc-test-*")
	if err != nil {
		t.Fatal(err)
	}
	zipPath := filepath.Join(tmpDir, "test.zip")
	if err := os.WriteFile(zipPath, []byte("zip"), 0o644); err != nil {
		t.Fatal(err)
	}

	arc := &stubArchiver{path: zipPath}
	upl := &stubUploader{}
	svc := NewArchiveService(arc, upl)

	_ = svc.RunJob(Job{Path: "/some/path"})

	if _, err := os.Stat(tmpDir); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("temp dir %q still exists after RunJob", tmpDir)
	}
}

func TestRunJob_ReturnsArchiveError(t *testing.T) {
	arcErr := errors.New("archive failed")
	arc := &stubArchiver{err: arcErr}
	upl := &stubUploader{}
	svc := NewArchiveService(arc, upl)

	if err := svc.RunJob(Job{}); !errors.Is(err, arcErr) {
		t.Errorf("got %v, want %v", err, arcErr)
	}
	if upl.received != "" {
		t.Error("uploader should not be called when archiver fails")
	}
}

func TestRunJob_ReturnsUploadError(t *testing.T) {
	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "test.zip")
	_ = os.WriteFile(zipPath, []byte("zip"), 0o644)

	uplErr := errors.New("upload failed")
	arc := &stubArchiver{path: zipPath}
	upl := &stubUploader{err: uplErr}
	svc := NewArchiveService(arc, upl)

	if err := svc.RunJob(Job{}); !errors.Is(err, uplErr) {
		t.Errorf("got %v, want %v", err, uplErr)
	}
}
