import type { ItemRequestStatus } from '$lib/api';

export type ItemRequestStatusFilter = 'all' | ItemRequestStatus;

export const itemRequestStatusOptions: { value: ItemRequestStatus; label: string }[] = [
	{ value: 'requested', label: 'Na čekanju' },
	{ value: 'approved', label: 'Odobreno' }
];

export const itemRequestStatusFilterOptions: { value: ItemRequestStatusFilter; label: string }[] = [
	{ value: 'all', label: 'Sve' },
	...itemRequestStatusOptions
];

export function itemRequestStatusLabel(status: ItemRequestStatus) {
	return itemRequestStatusOptions.find((option) => option.value === status)?.label ?? status;
}

export function itemRequestStatusClass(status: ItemRequestStatus) {
	return status === 'approved' ? 'bg-success-soft text-success' : 'bg-warning-soft text-warning';
}

export function formatRequestDate(value: string) {
	return new Date(value).toLocaleString('sr-RS', {
		day: '2-digit',
		month: '2-digit',
		year: 'numeric',
		hour: '2-digit',
		minute: '2-digit'
	});
}
