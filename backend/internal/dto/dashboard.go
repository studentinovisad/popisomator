package dto

// MonthCount is one bar of a monthly chart. Month is the first day of the bucket rather than
// YYYY-MM so the client can hand it straight to a date formatter.
type MonthCount struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

// ConsumptionBucket is one month of consumption, pivoted by the state each item was left in so the
// client can chart the three as stacked series without regrouping.
type ConsumptionBucket struct {
	Month             string `json:"month"`
	FullyConsumed     int64  `json:"fully_consumed"`
	PartiallyConsumed int64  `json:"partially_consumed"`
	Damaged           int64  `json:"damaged"`
}

// DashboardStockGroup is a StockGroup that names its own type, because a group name only identifies
// a stock line within one type and the dashboard counts across all of them.
type DashboardStockGroup struct {
	TypeID       int64               `json:"type_id"`
	TypeName     string              `json:"type_name"`
	Name         string              `json:"name"`
	InStockCount int64               `json:"in_stock_count"`
	TotalCount   int64               `json:"total_count"`
	Threshold    *int32              `json:"low_stock_count"`
	Low          bool                `json:"low"`
	Totals       []ItemPropertyTotal `json:"totals"`
}

type QuantityBucket struct {
	Month  string              `json:"month"`
	Totals []ItemPropertyTotal `json:"totals"`
}

type GroupConsumptionQuantity struct {
	Name         string              `json:"name"`
	Buckets      []QuantityBucket    `json:"buckets"`
	PeriodTotals []ItemPropertyTotal `json:"period_totals"`
}

type TypeConsumptionQuantity struct {
	TypeID   int64                      `json:"type_id"`
	TypeName string                     `json:"type_name"`
	Groups   []GroupConsumptionQuantity `json:"groups"`
}

type MostConsumedGroup struct {
	TypeID        int64               `json:"type_id"`
	TypeName      string              `json:"type_name"`
	Name          string              `json:"name"`
	ConsumedCount int64               `json:"consumed_count"`
	Totals        []ItemPropertyTotal `json:"totals"`
}

// DashboardExpiringItem is one item worth acting on before it is wasted. DaysRemaining is negative
// once the date has passed, which is the whole difference between a warning and a write-off.
type DashboardExpiringItem struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	TypeName      string `json:"type_name"`
	ExpiresOn     string `json:"expires_on"`
	DaysRemaining int32  `json:"days_remaining"`
}

// Dashboard is every widget's data in one response. The two ranges read in opposite directions:
// ExpiringByMonth covers the next Months months, ConsumptionByMonth the last Months.
//
// StockGroups is every group, not a shortlist. The shortages are the ones marked Low, which the
// client filters out for its own widget rather than being sent them twice.
//
// ExpiringItems is the only capped list here, so ExpiringItemsTotal says how many there were before
// the cap. Without it the client can only count what it was handed, and would report the ceiling as
// though it were the answer.
type Dashboard struct {
	Months              int32                     `json:"months"`
	TypeID              *int64                    `json:"type_id"`
	ExpiringByMonth     []MonthCount              `json:"expiring_by_month"`
	ExpiredBacklog      int64                     `json:"expired_backlog"`
	ExpiringItems       []DashboardExpiringItem   `json:"expiring_items"`
	ExpiringItemsTotal  int64                     `json:"expiring_items_total"`
	ConsumptionByMonth  []ConsumptionBucket       `json:"consumption_by_month"`
	StockGroups         []DashboardStockGroup     `json:"stock_groups"`
	ConsumptionQuantity []TypeConsumptionQuantity `json:"consumption_quantity"`
	MostConsumed        []MostConsumedGroup       `json:"most_consumed"`
}
