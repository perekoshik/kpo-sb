package importers

import (
	"encoding/json"

	"github.com/perekoshik/oop-go-concepts/internal/services/repository"
)

// JSONImporter implements Template Method for JSON payloads.
type JSONImporter struct {
	templateImporter
}

// Import loads data from JSON file.
func (i *JSONImporter) Import(path string) (*repository.StorageData, error) {
	raw, err := i.readFile(path)
	if err != nil {
		return nil, err
	}
	var data repository.StorageData
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
