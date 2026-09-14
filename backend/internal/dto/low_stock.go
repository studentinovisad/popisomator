package dto

// StockGroup is one line of stock: every item of a type that renders to the same derived name,
// counted together. InStockCount is what a manager can actually hand out today - untouched items
// nobody holds an approved request for - while TotalCount includes the consumed and the spoken-for,
// which is what keeps a group that has run dry visible instead of silently absent.
type StockGroup struct {
	Name         string `json:"name"`
	InStockCount int64  `json:"in_stock_count"`
	TotalCount   int64  `json:"total_count"`
	Low          bool   `json:"low"`
}

// ItemTypeStock is one item type's stock broken down by group. Threshold is nil when the type has no
// low stock warning configured, in which case no group is ever marked low.
type ItemTypeStock struct {
	TypeID    int64        `json:"type_id"`
	TypeName  string       `json:"type_name"`
	Threshold *int32       `json:"low_stock_count"`
	Groups    []StockGroup `json:"groups"`
}

// NotificationDescriptor_LowStock renders entirely from what was recorded when the warning fired.
// A group is a computed name with no row behind it, and the type may since have been renamed or had
// its threshold moved, so nothing here can be looked up after the fact.
type NotificationDescriptor_LowStock struct {
	TypeID    *int64 `json:"type_id"`
	TypeName  string `json:"type_name"`
	GroupName string `json:"group_name"`
	Threshold int32  `json:"threshold"`
	Observed  int32  `json:"observed"`
}
