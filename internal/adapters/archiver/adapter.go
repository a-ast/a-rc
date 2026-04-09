package archiver

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"a-rc/internal/core"
)

// Adapter implements core.Archiver using zip.
type Adapter struct{}

func New() *Adapter { return &Adapter{} }

// Archive creates a zip of `job.Path` in a temp directory and returns the zip file path.
// The caller is responsible for removing the file after use.
func (a *Adapter) Archive(job core.Job) (string, error) {
	src := job.Path
	info, err := os.Stat(src)
	if err != nil {
		return "", fmt.Errorf("stat source %q: %w", src, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("source %q is not a directory", src)
	}

	tmpDir, err := os.MkdirTemp("", "a-rc-*")
	if err != nil {
		return "", fmt.Errorf("creating temp dir: %w", err)
	}

	name := fmt.Sprintf("%s-%s.zip", filepath.Base(src), time.Now().Format("2006-01-02T15-04-05"))
	zipPath := filepath.Join(tmpDir, name)

	if err := writeZip(zipPath, src); err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}
	return zipPath, nil
}

func writeZip(zipPath, src string) error {
	f, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("creating zip file: %w", err)
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	// base is the parent of src so paths inside the zip are relative to src's parent,
	// preserving the top-level folder name.
	base := filepath.Dir(src)

	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Add directory entry (trailing slash required by zip spec).
			_, err = w.Create(rel + "/")
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = rel
		header.Method = zip.Deflate

		fw, err := w.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(fw, file)
		return err
	})
}
