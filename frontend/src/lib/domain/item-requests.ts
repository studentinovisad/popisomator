import type { ItemRequestStatus } from '$lib/api';

export type ItemRequestStatusFilter = 'all' | ItemRequestStatus;
export type ItemRequestDateFilter = 'all' | 'today' | 'last_7_days' | 'last_30_days' | 'this_month';

export const itemRequestStatusOptions: { value: ItemRequestStatus; label: string }[] = [
	{ value: 'requested', label: 'Na čekanju' },
	{ value: 'approved', label: 'Odobreno' }
];

export const itemRequestStatusFilterOptions: { value: ItemRequestStatusFilter; label: string }[] = [
	{ value: 'all', label: 'Sve' },
	...itemRequestStatusOptions
];

export const itemRequestDateFilterOptions: { value: ItemRequestDateFilter; label: string }[] = [
	{ value: 'all', label: 'Svi datumi' },
	{ value: 'today', label: 'Danas' },
	{ value: 'last_7_days', label: 'Poslednjih 7 dana' },
	{ value: 'last_30_days', label: 'Poslednjih 30 dana' },
	{ value: 'this_month', label: 'Ovaj mesec' }
];

export function itemRequestDateRange(filter: ItemRequestDateFilter, now = new Date()) {
	if (filter === 'all') return {};

	const createdFrom = new Date(now);
	createdFrom.setHours(0, 0, 0, 0);

	switch (filter) {
		case 'last_7_days':
			createdFrom.setDate(createdFrom.getDate() - 6);
			break;
		case 'last_30_days':
			createdFrom.setDate(createdFrom.getDate() - 29);
			break;
		case 'this_month':
			createdFrom.setDate(1);
	}

	return { createdFrom: createdFrom.toISOString() };
}

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
