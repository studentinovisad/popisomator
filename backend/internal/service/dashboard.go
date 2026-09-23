package service

import (
	"context"
	"sort"
	"strconv"
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

const dashboardMostConsumedLimit = 20

// GetDashboard collects every widget in one read. months is expected to already be one of the
// offered ranges; typeID narrows all four figures to one item type, or is nil for the whole
// inventory.
func GetDashboard(ctx context.Context, months int32, typeID *int64) (dto.Dashboard, error) {
	typeFilter := pgtype.Int8{}
	if typeID != nil {
		typeFilter = pgtype.Int8{Int64: *typeID, Valid: true}
	}
	unitValueTypes, unitNames, unitFactors := dto.MeasureUnitFactorRows()

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

	stockQuantityRows, err := db.Queries.SumStockQuantities(ctx, repository.SumStockQuantitiesParams{
		UnitValueTypes: unitValueTypes,
		UnitNames:      unitNames,
		UnitFactors:    unitFactors,
		TypeID:         typeFilter,
	})
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

	consumptionQuantityRows, err := db.Queries.SumConsumptionQuantityByMonth(ctx, repository.SumConsumptionQuantityByMonthParams{
		Months:         months,
		TypeID:         typeFilter,
		UnitValueTypes: unitValueTypes,
		UnitNames:      unitNames,
		UnitFactors:    unitFactors,
	})
	if err != nil {
		return dto.Dashboard{}, err
	}

	mostConsumedRows, err := db.Queries.ListMostConsumedGroups(ctx, repository.ListMostConsumedGroupsParams{
		Months:         months,
		TypeID:         typeFilter,
		LimitVal:       dashboardMostConsumedLimit,
		UnitValueTypes: unitValueTypes,
		UnitNames:      unitNames,
		UnitFactors:    unitFactors,
	})
	if err != nil {
		return dto.Dashboard{}, err
	}

	dashboard := dto.Dashboard{
		Months:              months,
		TypeID:              typeID,
		ExpiredBacklog:      expiredBacklog,
		ExpiringByMonth:     make([]dto.MonthCount, len(expiringRows)),
		ConsumptionByMonth:  make([]dto.ConsumptionBucket, len(consumptionRows)),
		ExpiringItems:       make([]dto.DashboardExpiringItem, len(expiringItemRows)),
		StockGroups:         make([]dto.DashboardStockGroup, len(stockRows)),
		ConsumptionQuantity: buildConsumptionQuantity(consumptionQuantityRows),
		MostConsumed:        buildMostConsumedGroups(mostConsumedRows),
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

	stockTotalsByGroup := make(map[dashboardGroupKey][]dto.PropertyTotalRow, len(stockQuantityRows))
	for _, row := range stockQuantityRows {
		key := dashboardGroupKey{typeID: row.TypeID, groupName: row.GroupName}
		stockTotalsByGroup[key] = append(stockTotalsByGroup[key], dto.PropertyTotalRow{
			PropertyID:   row.PropertyID,
			PropertyName: row.PropertyName,
			ValueType:    row.ValueType,
			TotalAmount:  row.TotalAmount,
			ValueCount:   row.ValueCount,
		})
	}

	for index, row := range stockRows {
		group := dto.DashboardStockGroup{
			TypeID:       row.TypeID,
			TypeName:     row.TypeName,
			Name:         stockGroupLabel(row.GroupName, row.TypeName),
			InStockCount: row.InStockCount,
			TotalCount:   row.TotalCount,
			Low:          isLowStock(row.InStockCount, row.LowStockCount.Int32, row.LowStockCount.Valid),
			Totals:       dto.BuildPropertyTotals(stockTotalsByGroup[dashboardGroupKey{typeID: row.TypeID, groupName: row.GroupName}]),
		}
		if row.LowStockCount.Valid {
			group.Threshold = &row.LowStockCount.Int32
		}
		dashboard.StockGroups[index] = group
	}

	return dashboard, nil
}

type dashboardGroupKey struct {
	typeID    int64
	groupName string
}

type dashboardPropertyKey struct {
	typeID     int64
	propertyID int64
}

func sumPropertyTotalRows(rows []dto.PropertyTotalRow) dto.PropertyTotalRow {
	sum := dto.PropertyTotalRow{
		PropertyID:   rows[0].PropertyID,
		PropertyName: rows[0].PropertyName,
		ValueType:    rows[0].ValueType,
		Currency:     rows[0].Currency,
	}

	var totalAmount, valueCount int64
	for _, row := range rows {
		amount, _ := strconv.ParseInt(row.TotalAmount, 10, 64)
		totalAmount += amount
		valueCount += row.ValueCount
	}

	sum.TotalAmount = strconv.FormatInt(totalAmount, 10)
	sum.ValueCount = valueCount
	return sum
}

func buildConsumptionQuantity(rows []repository.SumConsumptionQuantityByMonthRow) []dto.TypeConsumptionQuantity {
	typeOrder := make([]int64, 0)
	typeNames := make(map[int64]string)
	monthOrderByType := make(map[int64][]string)
	monthRowsByType := make(map[int64]map[string][]dto.PropertyTotalRow)
	periodRowsByKey := make(map[dashboardPropertyKey][]dto.PropertyTotalRow)

	for _, row := range rows {
		if _, seen := typeNames[row.TypeID]; !seen {
			typeOrder = append(typeOrder, row.TypeID)
			typeNames[row.TypeID] = row.TypeName
			monthRowsByType[row.TypeID] = make(map[string][]dto.PropertyTotalRow)
		}

		month := row.Month.Time.Format(time.DateOnly)
		if _, seen := monthRowsByType[row.TypeID][month]; !seen {
			monthOrderByType[row.TypeID] = append(monthOrderByType[row.TypeID], month)
		}

		totalRow := dto.PropertyTotalRow{
			PropertyID:   row.PropertyID,
			PropertyName: row.PropertyName,
			ValueType:    row.ValueType,
			Currency:     row.Currency,
			TotalAmount:  row.TotalAmount,
			ValueCount:   row.ValueCount,
		}
		monthRowsByType[row.TypeID][month] = append(monthRowsByType[row.TypeID][month], totalRow)

		periodKey := dashboardPropertyKey{typeID: row.TypeID, propertyID: row.PropertyID}
		periodRowsByKey[periodKey] = append(periodRowsByKey[periodKey], totalRow)
	}

	result := make([]dto.TypeConsumptionQuantity, 0, len(typeOrder))
	for _, typeID := range typeOrder {
		months := monthOrderByType[typeID]
		buckets := make([]dto.QuantityBucket, 0, len(months))
		for _, month := range months {
			buckets = append(buckets, dto.QuantityBucket{
				Month:  month,
				Totals: dto.BuildPropertyTotals(monthRowsByType[typeID][month]),
			})
		}

		periodRows := make([]dto.PropertyTotalRow, 0)
		for key, keyRows := range periodRowsByKey {
			if key.typeID == typeID {
				periodRows = append(periodRows, sumPropertyTotalRows(keyRows))
			}
		}
		sort.Slice(periodRows, func(i, j int) bool { return periodRows[i].PropertyID < periodRows[j].PropertyID })

		result = append(result, dto.TypeConsumptionQuantity{
			TypeID:       typeID,
			TypeName:     typeNames[typeID],
			Buckets:      buckets,
			PeriodTotals: dto.BuildPropertyTotals(periodRows),
		})
	}

	return result
}

func buildMostConsumedGroups(rows []repository.ListMostConsumedGroupsRow) []dto.MostConsumedGroup {
	groups := make([]dto.MostConsumedGroup, 0)
	groupIndexByKey := make(map[dashboardGroupKey]int)

	for _, row := range rows {
		key := dashboardGroupKey{typeID: row.TypeID, groupName: row.GroupName}
		groupIndex, exists := groupIndexByKey[key]
		if !exists {
			groupIndex = len(groups)
			groupIndexByKey[key] = groupIndex
			groups = append(groups, dto.MostConsumedGroup{
				TypeID:        row.TypeID,
				TypeName:      row.TypeName,
				Name:          stockGroupLabel(row.GroupName, row.TypeName),
				ConsumedCount: row.ConsumedCount,
			})
		}

		if !row.PropertyID.Valid {
			continue
		}

		groups[groupIndex].Totals = append(groups[groupIndex].Totals, dto.BuildPropertyTotals([]dto.PropertyTotalRow{{
			PropertyID:   row.PropertyID.Int64,
			PropertyName: row.PropertyName.String,
			ValueType:    row.ValueType.String,
			Currency:     row.Currency,
			TotalAmount:  row.TotalAmount.String,
			ValueCount:   row.ValueCount.Int64,
		}})...)
	}

	return groups
}
