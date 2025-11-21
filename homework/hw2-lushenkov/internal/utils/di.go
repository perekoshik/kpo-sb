package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/perekoshik/oop-go-concepts/internal/domain/factory"
	"github.com/perekoshik/oop-go-concepts/internal/services"
	"github.com/perekoshik/oop-go-concepts/internal/services/repository"
)

// Container wires all application services together.
type Container struct {
	Repository      repository.Repository
	Factory         *factory.DomainFactory
	AccountFacade   *services.AccountFacade
	CategoryFacade  *services.CategoryFacade
	OperationFacade *services.OperationFacade
	AnalyticsFacade *services.AnalyticsFacade
	DataManager     *services.DataManagementFacade
}

// NewContainer builds dependencies and ensures storage file exists.
func NewContainer(root string) (*Container, error) {
	storagePath := filepath.Join(root, "data", "db.json")
	if err := ensureStorageFile(storagePath); err != nil {
		return nil, err
	}

	fileRepo := repository.NewFileRepository(storagePath)
	proxyRepo := repository.NewProxyRepository(fileRepo)
	idGenerator := factory.UUIDGenerator{}
	domainFactory := factory.NewDomainFactory(idGenerator)

	accountFacade := services.NewAccountFacade(proxyRepo, domainFactory)
	categoryFacade := services.NewCategoryFacade(proxyRepo, domainFactory)
	operationFacade := services.NewOperationFacade(proxyRepo, domainFactory)
	analyticsFacade := services.NewAnalyticsFacade(proxyRepo)
	dataManager := services.NewDataManagementFacade(proxyRepo)

	container := &Container{
		Repository:      proxyRepo,
		Factory:         domainFactory,
		AccountFacade:   accountFacade,
		CategoryFacade:  categoryFacade,
		OperationFacade: operationFacade,
		AnalyticsFacade: analyticsFacade,
		DataManager:     dataManager,
	}
	return container, nil
}

func ensureStorageFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return fmt.Errorf("create storage dir: %w", err)
			}
			file, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("create storage file: %w", err)
			}
			defer file.Close()
			if err := json.NewEncoder(file).Encode(repository.StorageData{}); err != nil {
				return fmt.Errorf("initialise storage: %w", err)
			}
			return nil
		}
		return fmt.Errorf("stat storage file: %w", err)
	}
	return nil
}
