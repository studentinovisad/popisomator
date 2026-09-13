import { request } from '$lib/api/client';
import type {
	AuditActorOption,
	AuditEntry,
	AuditLogPage,
	ListAuditLogParams
} from '$lib/api/types';

export const auditLogApi = {
	// target_type and target_id only mean anything together: ids are unique within a kind of entity,
	// so an id on its own would mix an item's history with the item type sharing its number. The
	// backend rejects a lone target_id, so both are sent or neither.
	listAuditLog: ({
		limit = 20,
		offset = 0,
		action,
		targetType,
		targetID,
		actorID,
		createdFrom,
		createdTo
	}: ListAuditLogParams = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		if (action) query.set('action', action);
		if (targetType) query.set('target_type', targetType);
		if (targetType && targetID) query.set('target_id', String(targetID));
		if (actorID) query.set('actor_id', String(actorID));
		if (createdFrom) query.set('created_from', createdFrom);
		if (createdTo) query.set('created_to', createdTo);
		return request<AuditLogPage>(`/audit-log?${query}`);
	},
	getAuditLogEntry: (id: number) => request<AuditEntry>(`/audit-log/${id}`),
	listAuditLogActors: () => request<AuditActorOption[]>('/audit-log/actors')
};
