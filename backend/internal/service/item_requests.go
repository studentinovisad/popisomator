package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

func CreateItemRequest(ctx context.Context, req dto.ItemRequestCreateRequest) (dto.ItemRequest, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemRequest{}, err
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return dto.ItemRequest{}, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	if _, err := queriesTx.LockItemForRequest(ctx, req.ItemID); err != nil {
		return dto.ItemRequest{}, err
	}

	hasApprovedRequest, err := queriesTx.HasApprovedItemRequest(ctx, req.ItemID)
	if err != nil {
		return dto.ItemRequest{}, err
	}
	if hasApprovedRequest {
		return dto.ItemRequest{}, ErrItemReservedByApproval
	}

	itemRequest, err := queriesTx.CreateItemRequest(ctx, repository.CreateItemRequestParams{
		UserID: req.UserID,
		ItemID: req.ItemID,
		Status: repository.RequestStatusRequested,
		Reason: req.Reason,
	})
	if err != nil {
		return dto.ItemRequest{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return dto.ItemRequest{}, err
	}

	itemRequestDTO := dto.ToItemRequestDTO(itemRequest)

	return itemRequestDTO, nil
}

func GetItemRequest(ctx context.Context, req dto.ItemRequestIdentifierRequest) (dto.ItemRequest, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemRequest{}, err
	}

	itemRequest, err := db.Queries.GetItemRequest(ctx, repository.GetItemRequestParams{
		UserID: req.UserID,
		ItemID: req.ItemID,
	})
	if err != nil {
		return dto.ItemRequest{}, err
	}

	itemRequestDTO := dto.ToItemRequestDTO(itemRequest)

	return itemRequestDTO, nil
}

func CheckItemApproval(ctx context.Context, userID, itemID int64) (bool, error) {
	itemRequests, err := db.Queries.CheckItemsForRequests(ctx, []int64{itemID})
	if err != nil {
		return false, err
	}

	for _, itemRequest := range itemRequests {
		if itemRequest.UserID == userID {
			if itemRequest.Status == "approved" {
				return true, nil
			} else {
				return false, nil
			}
		}
	}

	return false, nil
}

func ApproveItemRequest(ctx context.Context, req dto.ItemRequestIdentifierRequest) (dto.ItemRequest, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemRequest{}, err
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return dto.ItemRequest{}, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	if _, err := queriesTx.LockItemForRequest(ctx, req.ItemID); err != nil {
		return dto.ItemRequest{}, err
	}

	itemRequest, err := queriesTx.ApproveItemRequest(ctx, repository.ApproveItemRequestParams{
		UserID: req.UserID,
		ItemID: req.ItemID,
	})
	if err != nil {
		return dto.ItemRequest{}, err
	}

	if _, err := queriesTx.DeleteNonApprovedItemRequests(ctx, req.ItemID); err != nil {
		return dto.ItemRequest{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.ItemRequest{}, err
	}

	itemRequestDTO := dto.ToItemRequestDTO(itemRequest)

	return itemRequestDTO, nil
}

func ListItemRequests(ctx context.Context, req dto.ItemRequestsListRequest) (dto.ItemRequestsPage, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemRequestsPage{}, err
	}

	userID := pgtype.Int8{}
	if req.UserID != nil {
		userID = pgtype.Int8{Int64: *req.UserID, Valid: true}
	}

	status := repository.NullRequestStatus{}
	if req.Status != nil {
		status = repository.NullRequestStatus{
			RequestStatus: repository.RequestStatus(*req.Status),
			Valid:         true,
		}
	}

	total, err := db.Queries.CountItemRequests(ctx, repository.CountItemRequestsParams{
		UserID: userID,
		Status: status,
	})
	if err != nil {
		return dto.ItemRequestsPage{}, err
	}

	requests, err := db.Queries.ListItemRequests(ctx, repository.ListItemRequestsParams{
		LimitVal:  req.Limit,
		OffsetVal: req.Offset,
		UserID:    userID,
		Status:    status,
	})
	if err != nil {
		return dto.ItemRequestsPage{}, err
	}

	requestsDTO := make([]dto.ItemRequestSummary, len(requests))
	for i, request := range requests {
		requestsDTO[i] = dto.ToItemRequestSummaryDTO(request)
	}

	return dto.ItemRequestsPage{
		Items:  requestsDTO,
		Limit:  req.Limit,
		Offset: req.Offset,
		Total:  total,
	}, nil
}

func ListItemRequestUsers(ctx context.Context) ([]dto.ItemRequestUserOption, error) {
	users, err := db.Queries.ListItemRequestUsers(ctx)
	if err != nil {
		return nil, err
	}

	options := make([]dto.ItemRequestUserOption, len(users))
	for i, user := range users {
		options[i] = dto.ItemRequestUserOption{ID: user.ID, Name: user.FullName}
	}

	return options, nil
}

func GetItemRequestPreparationReport(ctx context.Context, userID int64) (dto.ItemRequestPreparationReport, error) {
	user, err := db.Queries.GetUserByID(ctx, userID)
	if err != nil {
		return dto.ItemRequestPreparationReport{}, err
	}

	requests, err := db.Queries.ListItemPreparationRequests(ctx, userID)
	if err != nil {
		return dto.ItemRequestPreparationReport{}, err
	}

	items := make([]dto.ItemRequestPreparationItem, len(requests))
	itemIndexes := make(map[int64]int, len(requests))
	itemIDs := make([]int64, len(requests))
	for index, request := range requests {
		items[index] = dto.ItemRequestPreparationItem{
			ID:                request.ItemID,
			Name:              request.ItemName,
			TypeName:          request.ItemTypeName,
			DerivedNameFormat: request.DerivedNameFormat.String,
			Consumption:       request.Consumption,
			Reason:            request.Reason,
			RequestedAt:       request.CreatedAt.Time,
			Properties:        make([]dto.ItemRequestPreparationProperty, 0),
		}
		itemIndexes[request.ItemID] = index
		itemIDs[index] = request.ItemID
	}

	if len(itemIDs) > 0 {
		properties, err := db.Queries.GetItemProperties(ctx, itemIDs)
		if err != nil {
			return dto.ItemRequestPreparationReport{}, err
		}

		for _, property := range properties {
			itemIndex := itemIndexes[property.ItemProperty.ItemID]
			items[itemIndex].Properties = append(items[itemIndex].Properties, dto.ItemRequestPreparationProperty{
				Name:       property.PropertyName,
				Value:      property.ItemProperty.PropertyValue,
				Visibility: property.Visibility,
				Position:   property.Position,
			})
		}

	}

	return dto.ItemRequestPreparationReport{
		User:  dto.ItemRequestUserOption{ID: user.ID, Name: user.FullName},
		Items: items,
	}, nil
}

func DeleteItemRequest(ctx context.Context, req dto.ItemRequestIdentifierRequest) error {
	if err := dto.Validate(req); err != nil {
		return err
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	if _, err := queriesTx.LockItemForRequest(ctx, req.ItemID); err != nil {
		return err
	}

	rowsAffected, err := queriesTx.DeleteItemRequest(ctx, repository.DeleteItemRequestParams{
		UserID: req.UserID,
		ItemID: req.ItemID,
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
