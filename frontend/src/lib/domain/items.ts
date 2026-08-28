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
	{ value: 'price', label: 'Novčana vrednost' }
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
	value: string | null | undefined
) {
	if (value === null || value === undefined || value === '') return '—';

	try {
		const parsed = JSON.parse(value);
		if (valueType == 'price') {
			let value = (parsed.amount / PriceMultiplier).toLocaleString('sr-RS', {
				style: 'currency',
				currency: parsed.currency,
				minimumFractionDigits: 4,
				maximumFractionDigits: 4
			});
			return `${value}`;
		} else {
			return typeof parsed === 'string' ? parsed : JSON.stringify(parsed);
		}
	} catch {
		return value;
	}
}

export function defaultJsonValue(valueType: PropertyValueType, value: string | null) {
	if (value !== null) return value;
	if (valueType === 'string') return '""';
	if (valueType === 'number') return '0';
	if (valueType === 'price') return JSON.stringify({ amount: 0, currency: 'RSD' });
	return 'false';
}
