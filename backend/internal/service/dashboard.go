package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

// A ceiling rather than a page size: the window an item type sets keeps this list short on its own,
// but the expired end of it only ever grows, and nothing should be able to make one request return
// the whole inventory.
const dashboardExpiringItemsLimit = 200

// GetDashboard collects every widget in one read. months is expected to already be one of the
// offered ranges; typeID narrows all four figures to one item type, or is nil for the whole
// inventory.
func GetDashboard(ctx context.Context, months int32, typeID *int64) (dto.Dashboard, error) {
	typeFilter := pgtype.Int8{}
	if typeID != nil {
		typeFilter = pgtype.Int8{Int64: *typeID, Valid: true}
	}

	expiringRows, err := db.Queries.CountExpiringByMonth(ctx, repository.CountExpiringByMonthParams{
		Months: months,
		TypeID: typeFilter,
	})
	if err != nil {
		return dto.Dashboard{}, err
	}

	expiredBacklog, err := db.Queries.CountExpiredBacklog(ctx, typeFilter)
	if err != nil {
		return dto.Dashboard{}, err
	}

	consumptionRows, err := db.Queries.CountConsumptionByMonth(ctx, repository.CountConsumptionByMonthParams{
		Months: months,
		TypeID: typeFilter,
	})
	if err != nil {
		return dto.Dashboard{}, err
	}

	stockRows, err := db.Queries.ListStockGroups(ctx, typeFilter)
	if err != nil {
		return dto.Dashboard{}, err
	}

	expiringItemRows, err := db.Queries.ListExpiringItems(ctx, repository.ListExpiringItemsParams{
		LimitVal: dashboardExpiringItemsLimit,
		TypeID:   typeFilter,
	})
	if err != nil {
		return dto.Dashboard{}, err
	}

	dashboard := dto.Dashboard{
		Months:             months,
		TypeID:             typeID,
		ExpiredBacklog:     expiredBacklog,
		ExpiringByMonth:    make([]dto.MonthCount, len(expiringRows)),
		ConsumptionByMonth: make([]dto.ConsumptionBucket, len(consumptionRows)),
		ExpiringItems:      make([]dto.DashboardExpiringItem, len(expiringItemRows)),
		StockGroups:        make([]dto.DashboardStockGroup, len(stockRows)),
	}

	for index, row := range expiringRows {
		dashboard.ExpiringByMonth[index] = dto.MonthCount{
			Month: row.Month.Time.Format(time.DateOnly),
			Count: row.ItemCount,
		}
	}

	for index, row := range consumptionRows {
		dashboard.ConsumptionByMonth[index] = dto.ConsumptionBucket{
			Month:             row.Month.Time.Format(time.DateOnly),
			FullyConsumed:     row.FullyConsumed,
			PartiallyConsumed: row.PartiallyConsumed,
			Damaged:           row.Damaged,
		}
	}

	// Every row carries the same pre-cap total, so any one of them answers for all. No rows means
	// nothing matched, which the zero value already says.
	if len(expiringItemRows) > 0 {
		dashboard.ExpiringItemsTotal = expiringItemRows[0].TotalCount
	}

	for index, row := range expiringItemRows {
		dashboard.ExpiringItems[index] = dto.DashboardExpiringItem{
			ID:            row.ID,
			Name:          stockGroupLabel(row.GroupName, row.TypeName),
			TypeName:      row.TypeName,
			ExpiresOn:     row.ExpiresOn.Time.Format(time.DateOnly),
			DaysRemaining: row.DaysRemaining,
		}
	}

	for index, row := range stockRows {
		group := dto.DashboardStockGroup{
			TypeID:       row.TypeID,
			TypeName:     row.TypeName,
			Name:         stockGroupLabel(row.GroupName, row.TypeName),
			InStockCount: row.InStockCount,
			TotalCount:   row.TotalCount,
			Low:          isLowStock(row.InStockCount, row.LowStockCount.Int32, row.LowStockCount.Valid),
		}
		if row.LowStockCount.Valid {
			group.Threshold = &row.LowStockCount.Int32
		}
		dashboard.StockGroups[index] = group
	}

	return dashboard, nil
}
