package inventory

type InventoryItem interface {
	InventoryNumber() int
	DisplayName() string
}
