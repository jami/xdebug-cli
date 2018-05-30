package dbgp

import (
	"os"
	"path/filepath"
)

// NewSession creates a new session instance
func NewSession() *Session {
	return &Session{
		State:         StateNone,
		TransactionID: 1,
		CurrentLine:   1,
	}
}

// SetTargetFiles return all possible execution files
func (s *Session) SetTargetFiles(rootFile string) {
	rootDir := filepath.Dir(rootFile)
	s.TargetFiles = fileWalker(rootDir)
}

// NextTransactionID increments and returns the trans id
func (s *Session) NextTransactionID() int {
	s.TransactionID++
	return s.TransactionID
}

func fileWalker(root string) []string {
	files := []string{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})
	return files
}
