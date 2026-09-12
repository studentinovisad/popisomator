package service

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

// auditActor reads who to attribute a change to out of ctx. middleware.RequireAuth stores the id
// under "userID", the same key every controller already reads, so nothing has to pass an actor down
// the call stack. An absent id is not an error: it means the caller is not a request, and the entry
// carries no actor at all, which the client words however it likes.
func auditActor(ctx context.Context) pgtype.Int8 {
	id, ok := ctx.Value("userID").(int64)
	if !ok {
		return pgtype.Int8{}
	}

	return pgtype.Int8{Int64: id, Valid: true}
}

// auditReference is how a field that points at another row is recorded: the id to link by, and the
// name to show without having to resolve it again. Comparable, so auditDiff can tell two apart.
type auditReference struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// auditTarget is one row for writeAuditBulk.
type auditTarget struct {
	ID      int64
	Label   string
	Changes []dto.AuditChange
	Context dto.AuditContext
}

// encodeAuditPayload renders the two JSONB columns. Both are NOT NULL in the schema, so an absent
// value still has to be written as its empty literal rather than left off.
func encodeAuditPayload(changes []dto.AuditChange, entryContext dto.AuditContext) (json.RawMessage, json.RawMessage, error) {
	if changes == nil {
		changes = []dto.AuditChange{}
	}

	encodedChanges, err := json.Marshal(changes)
	if err != nil {
		return nil, nil, err
	}

	encodedContext, err := json.Marshal(entryContext)
	if err != nil {
		return nil, nil, err
	}

	return encodedChanges, encodedContext, nil
}

// writeAudit records one change. It takes the Querier the caller is already holding, the same way
// validatePropertyValue does, so passing db.Queries.WithTx(tx) puts the entry in the transaction of
// the change it describes. A failure here fails the operation on purpose: inside a transaction the
// statement has aborted it anyway, and a change that quietly went unrecorded is worse than one that
// did not happen.
func writeAudit(
	ctx context.Context,
	q repository.Querier,
	action repository.AuditAction,
	targetType repository.AuditTargetType,
	targetID int64,
	targetLabel string,
	changes []dto.AuditChange,
	entryContext dto.AuditContext,
) error {
	encodedChanges, encodedContext, err := encodeAuditPayload(changes, entryContext)
	if err != nil {
		return err
	}

	return q.WriteAuditEntry(ctx, repository.WriteAuditEntryParams{
		ActorID:     auditActor(ctx),
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		TargetLabel: targetLabel,
		Changes:     encodedChanges,
		Context:     encodedContext,
	})
}

// writeAuditBulk records one action against several targets in a single statement, the way
// CreateNotifications inserts one row per recipient. Writing nothing is not an error: a bulk add of
// zero items legitimately has no targets.
func writeAuditBulk(
	ctx context.Context,
	q repository.Querier,
	action repository.AuditAction,
	targetType repository.AuditTargetType,
	targets []auditTarget,
) error {
	if len(targets) == 0 {
		return nil
	}

	targetIDs := make([]int64, len(targets))
	targetLabels := make([]string, len(targets))
	encodedChanges := make([]json.RawMessage, len(targets))
	encodedContexts := make([]json.RawMessage, len(targets))

	for index, target := range targets {
		changes, entryContext, err := encodeAuditPayload(target.Changes, target.Context)
		if err != nil {
			return err
		}

		targetIDs[index] = target.ID
		targetLabels[index] = target.Label
		encodedChanges[index] = changes
		encodedContexts[index] = entryContext
	}

	return q.WriteAuditEntries(ctx, repository.WriteAuditEntriesParams{
		ActorID:      auditActor(ctx),
		Action:       action,
		TargetType:   targetType,
		TargetIds:    targetIDs,
		TargetLabels: targetLabels,
		Changes:      encodedChanges,
		Contexts:     encodedContexts,
	})
}

// auditDiff appends a change only when the value actually moved, so an update that sets a field to
// what it already held produces no entry at all. Marshalling failures drop the change rather than
// failing the write: a diff that cannot be rendered is worth less than the record that something
// changed at all.
func auditDiff[T comparable](
	changes []dto.AuditChange,
	key, valueType string,
	old, updated T,
) []dto.AuditChange {
	if old == updated {
		return changes
	}

	encodedOld, err := json.Marshal(old)
	if err != nil {
		return changes
	}

	encodedNew, err := json.Marshal(updated)
	if err != nil {
		return changes
	}

	return append(changes, dto.AuditChange{
		Key:       key,
		ValueType: valueType,
		Old:       encodedOld,
		New:       encodedNew,
	})
}

// auditRawDiff is auditDiff for values that are already JSON - property values, which are stored as
// jsonb and have no comparable Go shape. Equality is byte-wise on the encoded form.
func auditRawDiff(
	changes []dto.AuditChange,
	key, valueType string,
	old, updated json.RawMessage,
) []dto.AuditChange {
	if string(old) == string(updated) {
		return changes
	}

	return append(changes, dto.AuditChange{
		Key:       key,
		ValueType: valueType,
		Old:       old,
		New:       updated,
	})
}

func ListAuditLog(ctx context.Context, req dto.ListAuditLogRequest) (dto.AuditLogPage, error) {
	if err := dto.Validate(req); err != nil {
		return dto.AuditLogPage{}, err
	}

	action := repository.NullAuditAction{}
	if req.Action != nil {
		action = repository.NullAuditAction{AuditAction: *req.Action, Valid: true}
	}

	targetType := repository.NullAuditTargetType{}
	if req.TargetType != nil {
		targetType = repository.NullAuditTargetType{AuditTargetType: *req.TargetType, Valid: true}
	}

	targetID := pgtype.Int8{}
	if req.TargetID != nil {
		targetID = pgtype.Int8{Int64: *req.TargetID, Valid: true}
	}

	actorID := pgtype.Int8{}
	if req.ActorID != nil {
		actorID = pgtype.Int8{Int64: *req.ActorID, Valid: true}
	}

	createdFrom := pgtype.Timestamptz{}
	if req.CreatedFrom != nil {
		createdFrom = pgtype.Timestamptz{Time: *req.CreatedFrom, Valid: true}
	}

	createdTo := pgtype.Timestamptz{}
	if req.CreatedTo != nil {
		createdTo = pgtype.Timestamptz{Time: *req.CreatedTo, Valid: true}
	}

	total, err := db.Queries.CountAuditLog(ctx, repository.CountAuditLogParams{
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		ActorID:     actorID,
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
	})
	if err != nil {
		return dto.AuditLogPage{}, err
	}

	entries, err := db.Queries.ListAuditLog(ctx, repository.ListAuditLogParams{
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		ActorID:     actorID,
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		LimitVal:    req.Limit,
		OffsetVal:   req.Offset,
	})
	if err != nil {
		return dto.AuditLogPage{}, err
	}

	// Every row already carries everything it needs to render, so a page costs exactly these two
	// queries however many entries it holds.
	entriesDTO := make([]dto.AuditEntry, len(entries))
	for index, entry := range entries {
		entryDTO, err := dto.ToAuditEntryDTO(entry)
		if err != nil {
			return dto.AuditLogPage{}, err
		}
		entriesDTO[index] = entryDTO
	}

	return dto.AuditLogPage{
		Items:  entriesDTO,
		Limit:  req.Limit,
		Offset: req.Offset,
		Total:  total,
	}, nil
}

func ListAuditLogActors(ctx context.Context) ([]dto.AuditActorOption, error) {
	actors, err := db.Queries.ListAuditLogActors(ctx)
	if err != nil {
		return nil, err
	}

	options := make([]dto.AuditActorOption, len(actors))
	for index, actor := range actors {
		options[index] = dto.AuditActorOption{ID: actor.ID, Name: actor.Name}
	}

	return options, nil
}

// auditNew records a value that had no prior state, for the create actions.
func auditNew(changes []dto.AuditChange, key, valueType string, value any) []dto.AuditChange {
	encoded, err := json.Marshal(value)
	if err != nil {
		return changes
	}

	return append(changes, dto.AuditChange{
		Key:       key,
		ValueType: valueType,
		New:       encoded,
	})
}
