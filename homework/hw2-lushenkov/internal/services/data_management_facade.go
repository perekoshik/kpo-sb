package services

import (
	"fmt"
	"os"

	"github.com/perekoshik/oop-go-concepts/internal/domain/visitor"
	"github.com/perekoshik/oop-go-concepts/internal/services/importers"
	"github.com/perekoshik/oop-go-concepts/internal/services/repository"
)

// DataManagementFacade coordinates import/export workflows.
type DataManagementFacade struct {
	repo repository.Repository
}

// NewDataManagementFacade wires repository for data workflows.
func NewDataManagementFacade(repo repository.Repository) *DataManagementFacade {
	return &DataManagementFacade{repo: repo}
}

// ImportData loads snapshot using provided importer and replaces storage.
func (f *DataManagementFacade) ImportData(importer importers.Importer, path string) error {
	data, err := importer.Import(path)
	if err != nil {
		return err
	}
	if data == nil {
		data = &repository.StorageData{}
	}
	return f.repo.ReplaceAll(*data)
}

// ExportData persists snapshot using visitor to transform data into file format.
func (f *DataManagementFacade) ExportData(format visitor.ExportFormat, path string) error {
	exportVisitor, err := visitor.NewFileExportVisitor(format)
	if err != nil {
		return err
	}

	accounts, err := f.repo.ListAccounts()
	if err != nil {
		return err
	}
	categories, err := f.repo.ListCategories()
	if err != nil {
		return err
	}
	operations, err := f.repo.ListOperations()
	if err != nil {
		return err
	}

	for idx := range accounts {
		// iterate by index to take address of slice element
		if err := (&accounts[idx]).Accept(exportVisitor); err != nil {
			return err
		}
	}
	for idx := range categories {
		if err := (&categories[idx]).Accept(exportVisitor); err != nil {
			return err
		}
	}
	for idx := range operations {
		if err := (&operations[idx]).Accept(exportVisitor); err != nil {
			return err
		}
	}

	payload, err := exportVisitor.Render()
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		return fmt.Errorf("write export file: %w", err)
	}
	return nil
}
