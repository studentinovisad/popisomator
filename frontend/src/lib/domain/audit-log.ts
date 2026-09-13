import type { Pathname } from '$app/types';
import type {
	AuditAction,
	AuditChange,
	AuditEntry,
	AuditTargetType,
	PropertyValue,
	PropertyValueType
} from '$lib/api';
import { consumptionLabel, displayJson } from '$lib/domain/items';

export type AuditIconName = 'create' | 'edit' | 'consume' | 'delete' | 'request' | 'catalog';
export type AuditActionFilter = 'all' | AuditAction;
export type AuditTargetTypeFilter = 'all' | AuditTargetType;

// The backend records identifiers, never wording, so every Serbian string in the audit log lives
// here. Each action is one short phrase for a filter dropdown and one full sentence for a row.
const actionLabels: Record<AuditAction, string> = {
	item_create: 'Dodavanje stavke',
	item_update: 'Izmena stavke',
	item_consume: 'Promena stanja stavke',
	item_delete: 'Brisanje stavke',
	item_property_add: 'Dodavanje svojstva stavci',
	item_property_update: 'Izmena svojstva stavke',
	item_property_remove: 'Uklanjanje svojstva sa stavke',
	item_request_create: 'Zahtev za stavku',
	item_request_approve: 'Odobravanje zahteva',
	item_request_delete: 'Odbijanje zahteva',
	item_type_create: 'Dodavanje tipa stavke',
	item_type_update: 'Izmena tipa stavke',
	item_type_delete: 'Brisanje tipa stavke',
	item_type_property_add: 'Dodavanje svojstva tipu',
	item_type_property_update: 'Izmena svojstva tipa',
	item_type_property_remove: 'Uklanjanje svojstva sa tipa',
	item_type_property_reorder: 'Promena redosleda svojstava',
	property_create: 'Dodavanje svojstva',
	property_update: 'Izmena svojstva',
	property_delete: 'Brisanje svojstva'
};

const targetTypeLabels: Record<AuditTargetType, string> = {
	item: 'Stavka',
	item_type: 'Tip stavke',
	property: 'Svojstvo'
};

// The accusative, for sentences that take the target as an object - "Otvori stavku", "Istorija za
// stavku". Lowercasing the nominative would leave "Otvori stavka", which is not Serbian.
const targetTypeAccusatives: Record<AuditTargetType, string> = {
	item: 'stavku',
	item_type: 'tip stavke',
	property: 'svojstvo'
};

// What to call each column an entry can diff. Property changes are not here: they carry the
// property's own recorded name instead, which is data rather than a translatable string.
const changeKeyLabels: Record<string, string> = {
	name: 'Naziv',
	description: 'Opis',
	value_type: 'Tip vrednosti',
	default_value: 'Podrazumevana vrednost',
	derived_name_format: 'Format izvedenog naziva',
	expiring_soon_days: 'Prag isteka roka (dana)',
	visibility: 'Vidljivost',
	consumption: 'Stanje',
	type_id: 'Tip stavke',
	order: 'Redosled'
};

const visibilityLabels: Record<string, string> = {
	overview: 'Pregled',
	details: 'Detalji'
};

// Which entity each action is filed under, mirroring the target_type the backend writes with it.
// Spelled out rather than matched on prefixes: item_type_property_add and item_property_add differ
// by one word in the middle, and property_create shares its prefix with neither.
const actionTargetTypes: Record<AuditAction, AuditTargetType> = {
	item_create: 'item',
	item_update: 'item',
	item_consume: 'item',
	item_delete: 'item',
	item_property_add: 'item',
	item_property_update: 'item',
	item_property_remove: 'item',
	// Requests are filed under the item they concern - item_requests is keyed by a pair of columns,
	// which a single target_id cannot hold.
	item_request_create: 'item',
	item_request_approve: 'item',
	item_request_delete: 'item',
	item_type_create: 'item_type',
	item_type_update: 'item_type',
	item_type_delete: 'item_type',
	item_type_property_add: 'item_type',
	item_type_property_update: 'item_type',
	item_type_property_remove: 'item_type',
	item_type_property_reorder: 'item_type',
	property_create: 'property',
	property_update: 'property',
	property_delete: 'property'
};

export function auditActionTargetType(action: AuditAction) {
	return actionTargetTypes[action];
}

export const auditActionFilterOptions: { value: AuditActionFilter; label: string }[] = [
	{ value: 'all', label: 'Sve radnje' },
	...(Object.keys(actionLabels) as AuditAction[]).map((action) => ({
		value: action as AuditActionFilter,
		label: actionLabels[action]
	}))
];

export const auditTargetTypeFilterOptions: { value: AuditTargetTypeFilter; label: string }[] = [
	{ value: 'all', label: 'Sve vrste' },
	...(Object.keys(targetTypeLabels) as AuditTargetType[]).map((targetType) => ({
		value: targetType as AuditTargetTypeFilter,
		label: targetTypeLabels[targetType]
	}))
];

// The actions worth offering once a kind of entity is chosen. Showing all twenty after picking
// "Svojstvo" means seventeen of them can only ever return nothing, which makes the filter work
// against whoever is reading.
export function auditActionFilterOptionsFor(targetType: AuditTargetTypeFilter) {
	if (targetType === 'all') return auditActionFilterOptions;

	const [allOption, ...actionOptions] = auditActionFilterOptions;
	return [
		allOption,
		...actionOptions.filter(
			(option) => actionTargetTypes[option.value as AuditAction] === targetType
		)
	];
}

// Whether an action can occur at all for a kind of entity, so a selection left over from a different
// one can be dropped rather than silently filtering everything away.
export function auditActionAllowedFor(
	action: AuditActionFilter,
	targetType: AuditTargetTypeFilter
) {
	if (action === 'all' || targetType === 'all') return true;
	return actionTargetTypes[action] === targetType;
}

export function auditActionLabel(action: AuditAction) {
	return actionLabels[action] ?? action;
}

export function auditTargetTypeLabel(targetType: AuditTargetType) {
	return targetTypeLabels[targetType] ?? targetType;
}

export function auditTargetTypeAccusative(targetType: AuditTargetType) {
	return targetTypeAccusatives[targetType] ?? targetType;
}

// Whoever the entry is filed under, named however it can be. An entry outlives both its target and
// its actor, so neither name is guaranteed to still resolve - the recorded one is all there is.
export function auditActorName(entry: AuditEntry) {
	// Both empty means nothing acted: the seeder, or a background job.
	if (!entry.actor_name) return 'Sistem';
	return entry.actor_name;
}

export function auditTargetName(entry: AuditEntry) {
	if (entry.target_label) return entry.target_label;

	switch (entry.target_type) {
		case 'item':
			return `Stavka #${entry.target_id}`;
		case 'item_type':
			return `Tip stavke #${entry.target_id}`;
		case 'property':
			return `Svojstvo #${entry.target_id}`;
	}
}

function subjectName(entry: AuditEntry) {
	const { subject_user_name, subject_user_id } = entry.context;
	if (subject_user_name) return subject_user_name;
	return subject_user_id ? `Korisnik #${subject_user_id}` : 'korisnik';
}

function propertyName(entry: AuditEntry) {
	return entry.context.property_name ?? 'svojstvo';
}

// The sentence a row leads with. Built here rather than on the backend so the wording stays a
// frontend concern, the same way notificationTitle does it.
export function auditEntryTitle(entry: AuditEntry): string {
	const actor = auditActorName(entry);
	const target = auditTargetName(entry);
	const property = propertyName(entry);
	const subject = subjectName(entry);

	switch (entry.action) {
		case 'item_create':
			return `„${actor}” je dodao stavku „${target}”`;
		case 'item_update':
			return `„${actor}” je izmenio stavku „${target}”`;
		case 'item_consume':
			return `„${actor}” je promenio stanje stavke „${target}”`;
		case 'item_delete':
			return `„${actor}” je obrisao stavku „${target}”`;
		case 'item_property_add':
			return `„${actor}” je dodao svojstvo „${property}” stavci „${target}”`;
		case 'item_property_update':
			return `„${actor}” je izmenio svojstvo „${property}” stavke „${target}”`;
		case 'item_property_remove':
			return `„${actor}” je uklonio svojstvo „${property}” sa stavke „${target}”`;
		case 'item_request_create':
			// Comparing the names rather than the ids: an actor who has since been deleted has no id
			// left to match on, and naming the same person twice in one sentence reads as a bug.
			return subject === actor
				? `„${actor}” je zatražio stavku „${target}”`
				: `„${actor}” je zatražio stavku „${target}” za korisnika „${subject}”`;
		case 'item_request_approve':
			return `„${actor}” je odobrio zahtev korisnika „${subject}” za stavku „${target}”`;
		case 'item_request_delete':
			// The status it held is what separates turning down a pending request from taking back one
			// that was already approved.
			return entry.context.request_status === 'approved'
				? `„${actor}” je uklonio odobrenje korisnika „${subject}” za stavku „${target}”`
				: `„${actor}” je odbio zahtev korisnika „${subject}” za stavku „${target}”`;
		case 'item_type_create':
			return `„${actor}” je napravio tip stavke „${target}”`;
		case 'item_type_update':
			return `„${actor}” je izmenio tip stavke „${target}”`;
		case 'item_type_delete':
			return `„${actor}” je obrisao tip stavke „${target}”`;
		case 'item_type_property_add':
			return `„${actor}” je dodao svojstvo „${property}” tipu „${target}”`;
		case 'item_type_property_update':
			return `„${actor}” je izmenio svojstvo „${property}” tipa „${target}”`;
		case 'item_type_property_remove':
			return `„${actor}” je uklonio svojstvo „${property}” sa tipa „${target}”`;
		case 'item_type_property_reorder':
			return `„${actor}” je promenio redosled svojstava tipa „${target}”`;
		case 'property_create':
			return `„${actor}” je napravio svojstvo „${target}”`;
		case 'property_update':
			return `„${actor}” je izmenio svojstvo „${target}”`;
		case 'property_delete':
			return `„${actor}” je obrisao svojstvo „${target}”`;
		default:
			return `„${actor}” je izmenio „${target}”`;
	}
}

// A second line for a row: the reason behind a request, or a short read of what moved.
export function auditEntryDetail(entry: AuditEntry): string {
	if (entry.context.reason) return entry.context.reason;
	if (entry.changes.length === 0) return '';

	return entry.changes.map((change) => auditChangeLabel(change)).join(' · ');
}

export function auditEntryIcon(entry: AuditEntry): AuditIconName {
	if (entry.action.startsWith('item_request')) return 'request';
	if (entry.action.endsWith('_delete') || entry.action.endsWith('_remove')) return 'delete';
	if (entry.action === 'item_consume') return 'consume';
	if (entry.action.endsWith('_create') || entry.action.endsWith('_add')) return 'create';
	if (entry.target_type !== 'item') return 'catalog';
	return 'edit';
}

export function auditActionClass(action: AuditAction) {
	if (action.endsWith('_delete') || action.endsWith('_remove')) return 'bg-danger-soft text-danger';
	if (action.endsWith('_create') || action.endsWith('_add')) return 'bg-success-soft text-success';
	if (action === 'item_request_approve') return 'bg-brand-soft text-brand';
	return 'bg-soft text-muted';
}

// Where the target still lives, for entries whose target survives. A deletion has nothing left to
// open, so it gets no link.
export function auditTargetLink(entry: AuditEntry): Pathname | null {
	if (entry.action === 'item_delete' && entry.target_type === 'item') return null;
	if (entry.action === 'item_type_delete' && entry.target_type === 'item_type') return null;
	if (entry.action === 'property_delete' && entry.target_type === 'property') return null;

	switch (entry.target_type) {
		case 'item':
			return `/items/${entry.target_id}`;
		case 'item_type':
			return `/catalog/item-types/${entry.target_id}`;
		case 'property':
			return `/catalog/properties/${entry.target_id}`;
	}
}

// The target's own history, which outlives the target itself - so this is offered even for entries
// that record a deletion, and is the only link those get. Pathname does not model a query string,
// so this is cast, the same way notificationLink casts its filtered link.
export function auditHistoryLink(targetType: AuditTargetType, targetID: number): Pathname {
	return `/admin/audit-log?target_type=${targetType}&target_id=${targetID}` as Pathname;
}

export function auditChangeLabel(change: AuditChange) {
	if (change.key === 'property') return change.label ?? 'Svojstvo';
	return changeKeyLabels[change.key] ?? change.key;
}

// Renders one side of a diff. value_type is snapshotted on the change itself, so a recorded price or
// mass goes through the same helper a live one does and needs no lookup.
export function auditChangeValue(change: AuditChange, side: 'old' | 'new') {
	const value = side === 'old' ? change.old : change.new;
	if (value === null || value === undefined) return '—';

	switch (change.value_type) {
		case 'consumption':
			return consumptionLabel(value as Parameters<typeof consumptionLabel>[0]);
		case 'visibility':
			return visibilityLabels[String(value)] ?? String(value);
		case 'reference':
			return (value as { name?: string }).name ?? String(value);
		case 'order':
			return Array.isArray(value) ? value.join(' → ') : String(value);
		case 'text':
		case 'count':
			return String(value);
		default:
			// Anything left is one of the property value types, which displayJson already renders.
			return displayJson(change.value_type as PropertyValueType, value as PropertyValue);
	}
}
