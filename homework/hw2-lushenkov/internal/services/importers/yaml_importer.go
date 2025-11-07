package importers

import (
	"strconv"
	"strings"
	"time"

	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
	"github.com/perekoshik/oop-go-concepts/internal/services/repository"
)

// YAMLImporter imports data serialized by visitor.renderYAML.
type YAMLImporter struct {
	templateImporter
}

// Import parses simple YAML-like snapshot.
func (i *YAMLImporter) Import(path string) (*repository.StorageData, error) {
	raw, err := i.readFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(raw), "\n")
	section := ""
	current := map[string]string(nil)
	data := repository.StorageData{}

	flush := func() {
		if current == nil || section == "" {
			return
		}
		switch section {
		case "accounts":
			balance, _ := strconv.ParseFloat(current["balance"], 64)
			data.Accounts = append(data.Accounts, models.BankAccount{
				ID:      current["id"],
				Name:    unquote(current["name"]),
				Balance: balance,
			})
		case "categories":
			data.Categories = append(data.Categories, models.Category{
				ID:   current["id"],
				Name: unquote(current["name"]),
				Type: models.CategoryType(current["type"]),
			})
		case "operations":
			amount, _ := strconv.ParseFloat(current["amount"], 64)
			date, _ := time.Parse(time.RFC3339, current["date"])
			data.Operations = append(data.Operations, models.Operation{
				ID:            current["id"],
				Type:          models.OperationType(current["type"]),
				BankAccountID: current["bank_account_id"],
				CategoryID:    current["category_id"],
				Amount:        amount,
				Date:          date,
				Description:   unquote(current["description"]),
			})
		}
		current = nil
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !strings.HasPrefix(trimmed, "-") && strings.HasSuffix(trimmed, ":") {
			flush()
			section = strings.TrimSuffix(trimmed, ":")
			continue
		}
		if strings.HasPrefix(trimmed, "-") {
			flush()
			current = make(map[string]string)
			remainder := strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
			if remainder != "" {
				key, value := splitKeyValue(remainder)
				current[key] = value
			}
			continue
		}
		if current != nil {
			key, value := splitKeyValue(trimmed)
			current[key] = value
		}
	}
	flush()
	return &data, nil
}

func splitKeyValue(line string) (string, string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return strings.TrimSpace(line), ""
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	return key, value
}

func unquote(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
		return value[1 : len(value)-1]
	}
	return value
}
