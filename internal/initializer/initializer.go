package initializer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrAlreadyInitialized = errors.New("AOF is already initialized in this project")

type ConflictError struct{ Paths []string }

func (e *ConflictError) Error() string {
	return fmt.Sprintf("initialization conflicts with %d existing file(s)", len(e.Paths))
}

func ValidatePlan(root string, p Plan) error {
	if _, err := os.Stat(filepath.Join(root, ".aof", "manifest.json")); err == nil {
		return ErrAlreadyInitialized
	}
	var conflicts []string
	for _, f := range p.Files {
		if !safePath(f.Path) {
			return fmt.Errorf("unsafe generated path: %q", f.Path)
		}
		target := filepath.Join(root, filepath.Clean(f.Path))
		rel, err := filepath.Rel(root, target)
		if err != nil || rel == ".." || filepath.IsAbs(rel) || len(rel) >= 3 && rel[:3] == "../" {
			return fmt.Errorf("generated path escapes project root: %q", f.Path)
		}
		if _, err := os.Stat(target); err == nil {
			conflicts = append(conflicts, f.Path)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("inspect %s: %w", f.Path, err)
		}
	}
	if len(conflicts) > 0 {
		return &ConflictError{Paths: conflicts}
	}
	return nil
}

func Apply(root string, p Plan) error {
	if err := ValidatePlan(root, p); err != nil {
		return err
	}
	created := make([]string, 0, len(p.Files))
	rollback := func() {
		for i := len(created) - 1; i >= 0; i-- {
			_ = os.Remove(created[i])
		}
	}
	for _, f := range p.Files {
		target := filepath.Join(root, filepath.Clean(f.Path))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			rollback()
			return fmt.Errorf("create directory for %s: %w", f.Path, err)
		}
		if err := os.WriteFile(target, f.Content, 0o644); err != nil {
			rollback()
			return fmt.Errorf("write %s: %w", f.Path, err)
		}
		created = append(created, target)
	}
	return nil
}
