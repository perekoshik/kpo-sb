package importers

import (
	"fmt"
	"os"

	"github.com/perekoshik/oop-go-concepts/internal/services/repository"
)

// Importer defines contract for domain data import implementations.
type Importer interface {
	Import(path string) (*repository.StorageData, error)
}

type templateImporter struct{}

func (templateImporter) readFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}
	return content, nil
}
