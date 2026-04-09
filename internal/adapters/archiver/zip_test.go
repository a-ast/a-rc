package archiver

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"a-rc/internal/domain"
)

func TestArchive_CreatesZip(t *testing.T) {
	src := t.TempDir()
	writeFile(t, src, "file.txt", "hello")
	writeFile(t, src, filepath.Join("sub", "nested.txt"), "world")

	a := New()
	zipPath, err := a.Archive(domain.Job{Name: "mybackup", Path: src})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.RemoveAll(filepath.Dir(zipPath))

	// zip file must exist
	if _, err := os.Stat(zipPath); err != nil {
		t.Fatalf("zip file not found: %v", err)
	}

	// zip name must be exactly {job.Name}.zip
	if filepath.Base(zipPath) != "mybackup.zip" {
		t.Errorf("zip name %q, want %q", filepath.Base(zipPath), "mybackup.zip")
	}

	base := filepath.Base(src)
	entries := zipEntries(t, zipPath)
	assertContains(t, entries, filepath.Join(base, "file.txt"))
	assertContains(t, entries, filepath.Join(base, "sub", "nested.txt"))
}

func TestArchive_ErrorOnMissingSource(t *testing.T) {
	a := New()
	_, err := a.Archive(domain.Job{Path: "/nonexistent/path"})
	if err == nil {
		t.Fatal("expected error for missing source, got nil")
	}
}

func TestArchive_ErrorOnFile(t *testing.T) {
	f, err := os.CreateTemp("", "a-rc-test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	a := New()
	_, err = a.Archive(domain.Job{Path: f.Name()})
	if err == nil {
		t.Fatal("expected error for non-directory source, got nil")
	}
}

// helpers

func writeFile(t *testing.T, dir, rel, content string) {
	t.Helper()
	full := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func zipEntries(t *testing.T, zipPath string) []string {
	t.Helper()
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("opening zip: %v", err)
	}
	defer r.Close()
	var names []string
	for _, f := range r.File {
		if !strings.HasSuffix(f.Name, "/") { // skip dir entries
			names = append(names, f.Name)
		}
	}
	return names
}

func assertContains(t *testing.T, entries []string, want string) {
	t.Helper()
	for _, e := range entries {
		if e == want {
			return
		}
	}
	t.Errorf("zip does not contain %q; entries: %v", want, entries)
}
