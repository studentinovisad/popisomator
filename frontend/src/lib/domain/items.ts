import { type ConsumptionStatus, type PropertyValueType, type PTPrice } from '$lib/api';
import * as api from '$lib/api';

export const consumptionOptions: { value: ConsumptionStatus; label: string }[] = [
	{ value: 'not_consumed', label: 'Nije potrošeno' },
	{ value: 'partially_consumed', label: 'Delimično potrošeno' },
	{ value: 'fully_consumed', label: 'Potrošeno' },
	{ value: 'damaged', label: 'Oštećeno' }
];

export const propertyValueTypeOptions: { value: PropertyValueType; label: string }[] = [
	{ value: 'string', label: 'Tekst' },
	{ value: 'number', label: 'Broj' },
	{ value: 'boolean', label: 'Da / ne' },
	{ value: 'price', label: 'Novčana vrednost' },
	{ value: 'expiry', label: 'Datum isteka roka' }
];

export const PriceMultiplier: number = 10000;

export function consumptionLabel(status: ConsumptionStatus) {
	return consumptionOptions.find((option) => option.value === status)?.label ?? status;
}

export function consumptionClass(status: ConsumptionStatus) {
	if (status === 'fully_consumed') return 'bg-soft text-muted';
	if (status === 'damaged') return 'bg-danger-soft text-danger';
	if (status === 'partially_consumed') return 'bg-warning-soft text-warning';
	return 'bg-success-soft text-success';
}

export function propertyValueTypeLabel(valueType: PropertyValueType) {
	return propertyValueTypeOptions.find((option) => option.value === valueType)?.label ?? valueType;
}

export function displayJson(
	valueType: PropertyValueType | undefined,
	value: {} | null | undefined
) {
	if (value === null || value === undefined || value === '') return '—';

	try {
		if (valueType == 'price') {
			const price = value as PTPrice
			return (price.amount / PriceMultiplier).toLocaleString('sr-RS', {
				style: 'currency',
				currency: price.currency,
				minimumFractionDigits: 2,
				maximumFractionDigits: 4
			});
		} else {
			return typeof value === 'string' ? value : JSON.stringify(value);
		}
	} catch {
		return value;
	}
}

export function defaultJsonValue(valueType: PropertyValueType, value: {} | null) {
	if (value !== null) return value;
	if (valueType === 'string') return "";
	if (valueType === 'number') return 0;
	if (valueType === 'price') return { amount: 0, currency: 'RSD' };
	if (valueType === 'expiry') {
		let sample_date = new Date();
		sample_date.setDate(sample_date.getDate() + 7)
		return sample_date.toISOString().split('T')[0];
	}
	return false;
}
