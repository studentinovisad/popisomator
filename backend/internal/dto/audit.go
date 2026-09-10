package dto

import (
	"encoding/json"
	"time"

	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

// Value types an AuditChange can carry beyond the seven property value types, which are reused
// verbatim so the frontend can render a recorded value with the same helper it uses for a live one.
const (
	AuditValueTypeConsumption = "consumption"
	AuditValueTypeVisibility  = "visibility"
	// AuditValueTypeReference is an {id, name} pair, for a field that points at another row.
	AuditValueTypeReference = "reference"
	// AuditValueTypeOrder is an array of names, for a reordering.
	AuditValueTypeOrder = "order"
	AuditValueTypeText  = "text"
	AuditValueTypeCount = "count"
)

// AuditChange is one field that moved. Old and New stay raw JSON so a property value reaches the
// frontend in the shape its renderers already expect - a price, a mass - instead of being flattened
// to a string here and parsed back there.
//
// Nothing here is worded for a reader. Key and ValueType are stable identifiers the client turns
// into whatever language it renders in; the backend stays language-neutral.
type AuditChange struct {
	// Key names the column that changed ("name", "consumption", "type_id"), or "property" for a
	// value carried on an item or item type. The client maps it to a heading.
	Key string `json:"key"`
	// Label is the recorded name of the user-defined property this change concerns - data out of the
	// properties table, not a translatable string. Empty for the fixed columns above, which Key
	// already identifies. Only worth setting when one entry spans several properties and the name
	// cannot be read off context instead, as when an item is created with its initial values.
	Label string `json:"label,omitempty"`
	// ValueType tells the client how to render Old and New.
	ValueType string          `json:"value_type,omitempty"`
	Old       json.RawMessage `json:"old,omitempty"`
	New       json.RawMessage `json:"new,omitempty"`
}

// AuditProperty is a property as it read when the entry was written - the list an item type was
// created with. It names the property rather than carrying its value: what a value held over time is
// the business of the diffs on that target's own entries.
type AuditProperty struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ValueType string `json:"value_type"`
}

// AuditContext is everything an entry needs that is not a field diff. Every field is optional; which
// ones are filled depends on the action.
type AuditContext struct {
	TypeID   *int64 `json:"type_id,omitempty"`
	TypeName string `json:"type_name,omitempty"`
	// BatchSize is how many items one bulk add produced. Each item still gets its own entry, so its
	// own timeline is complete; this is what lets a feed say the add was part of a batch.
	BatchSize int32 `json:"batch_size,omitempty"`

	PropertyID   *int64 `json:"property_id,omitempty"`
	PropertyName string `json:"property_name,omitempty"`

	// The user a request concerns, who is not always the actor: an admin can request on someone
	// else's behalf, and a supersede is recorded against the person who lost their place.
	SubjectUserID   *int64 `json:"subject_user_id,omitempty"`
	SubjectUserName string `json:"subject_user_name,omitempty"`
	Reason          string `json:"reason,omitempty"`
	// RequestStatus is what the request was before it was deleted, which is what separates rejecting
	// a pending request from revoking an approved one.
	RequestStatus string `json:"request_status,omitempty"`
	// Cascaded marks a request that was not turned down on its own but went away with the item it
	// stood against. Same action either way; this is what lets the two be worded differently.
	Cascaded bool `json:"cascaded,omitempty"`

	// Properties is the list an item type was created with.
	Properties []AuditProperty `json:"properties,omitempty"`
}

type AuditEntry struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	// ActorID is null once that user is deleted, while ActorName still holds who it was. Both empty
	// means there was no user behind the change at all - the seeder, or a background job - which the
	// client names in its own language rather than being handed a word from here.
	ActorID     *int64                     `json:"actor_id"`
	ActorName   string                     `json:"actor_name"`
	Action      repository.AuditAction     `json:"action"`
	TargetType  repository.AuditTargetType `json:"target_type"`
	TargetID    int64                      `json:"target_id"`
	TargetLabel string                     `json:"target_label"`
	Changes     []AuditChange              `json:"changes"`
	Context     AuditContext               `json:"context"`
}

type AuditLogPage struct {
	Items  []AuditEntry `json:"items"`
	Limit  int32        `json:"limit"`
	Offset int32        `json:"offset"`
	Total  int64        `json:"total"`
}

// AuditActorOption is someone who has made a recorded change, for the actor filter.
type AuditActorOption struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ListAuditLogRequest struct {
	Limit      int32
	Offset     int32
	Action     *repository.AuditAction
	TargetType *repository.AuditTargetType
	// TargetID only means something alongside TargetType; the controller rejects it on its own.
	TargetID    *int64 `validate:"omitempty,gt=0"`
	ActorID     *int64 `validate:"omitempty,gt=0"`
	CreatedFrom *time.Time
	CreatedTo   *time.Time
}

// ToAuditEntryDTO decodes the two JSONB columns into their typed shape. A row that fails to decode
// is a bug in whatever wrote it, so the error is returned rather than swallowed into a blank entry.
func ToAuditEntryDTO(entry repository.AuditLog) (AuditEntry, error) {
	dtoEntry := AuditEntry{
		ID:          entry.ID,
		CreatedAt:   entry.CreatedAt.Time,
		ActorName:   entry.ActorName,
		Action:      entry.Action,
		TargetType:  entry.TargetType,
		TargetID:    entry.TargetID,
		TargetLabel: entry.TargetLabel,
		Changes:     make([]AuditChange, 0),
	}

	if entry.ActorID.Valid {
		actorID := entry.ActorID.Int64
		dtoEntry.ActorID = &actorID
	}

	if len(entry.Changes) > 0 {
		if err := json.Unmarshal(entry.Changes, &dtoEntry.Changes); err != nil {
			return AuditEntry{}, err
		}
	}

	if len(entry.Context) > 0 {
		if err := json.Unmarshal(entry.Context, &dtoEntry.Context); err != nil {
			return AuditEntry{}, err
		}
	}

	return dtoEntry, nil
}

// auditActions is the closed set of recorded actions, so a filter naming an unknown one is rejected
// with a 400 instead of reaching Postgres and failing the enum cast as a 500.
var auditActions = map[repository.AuditAction]struct{}{
	repository.AuditActionItemCreate:              {},
	repository.AuditActionItemUpdate:              {},
	repository.AuditActionItemConsume:             {},
	repository.AuditActionItemDelete:              {},
	repository.AuditActionItemPropertyAdd:         {},
	repository.AuditActionItemPropertyUpdate:      {},
	repository.AuditActionItemPropertyRemove:      {},
	repository.AuditActionItemRequestCreate:       {},
	repository.AuditActionItemRequestApprove:      {},
	repository.AuditActionItemRequestDelete:       {},
	repository.AuditActionItemRequestSupersede:    {},
	repository.AuditActionItemTypeCreate:          {},
	repository.AuditActionItemTypeUpdate:          {},
	repository.AuditActionItemTypeDelete:          {},
	repository.AuditActionItemTypePropertyAdd:     {},
	repository.AuditActionItemTypePropertyUpdate:  {},
	repository.AuditActionItemTypePropertyRemove:  {},
	repository.AuditActionItemTypePropertyReorder: {},
	repository.AuditActionPropertyCreate:          {},
	repository.AuditActionPropertyUpdate:          {},
	repository.AuditActionPropertyDelete:          {},
}

func ParseAuditAction(value string) (repository.AuditAction, bool) {
	action := repository.AuditAction(value)
	_, ok := auditActions[action]

	return action, ok
}

func ParseAuditTargetType(value string) (repository.AuditTargetType, bool) {
	switch targetType := repository.AuditTargetType(value); targetType {
	case repository.AuditTargetTypeItem,
		repository.AuditTargetTypeItemType,
		repository.AuditTargetTypeProperty:
		return targetType, true
	default:
		return "", false
	}
}
