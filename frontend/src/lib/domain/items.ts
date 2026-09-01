import {
	type ConsumptionStatus,
	type PropertyValueType,
	type PTMeasure,
	type PTPrice
} from '$lib/api';
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
	{ value: 'expiry', label: 'Datum isteka roka' },
	{ value: 'mass', label: 'Masa' },
	{ value: 'volume', label: 'Zapremina' }
];

export const PriceMultiplier: number = 10000;

// Mirrors dto.MeasureMultiplier on the backend: a measured amount is stored as an integer in the
// unit the user picked, multiplied by this.
export const MeasureMultiplier: number = 10000;

export const massUnits: string[] = ['mg', 'g', 'kg'];
export const volumeUnits: string[] = ['mL', 'L'];

// Factors into the canonical base unit of each dimension - milligram for mass, millilitre for
// volume - so amounts recorded in different units can be summed. Mirrors dto.MassUnitFactors and
// dto.VolumeUnitFactors; keep in sync with massUnits / volumeUnits above.
export const massUnitFactors: Record<string, number> = { mg: 1, g: 1_000, kg: 1_000_000 };
export const volumeUnitFactors: Record<string, number> = { mL: 1, L: 1_000 };

export function measureUnits(valueType: PropertyValueType) {
	return valueType === 'mass' ? massUnits : volumeUnits;
}

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

// Compares two property values by content. The structured types (price, mass, volume) are objects,
// so `===` only ever reports reference identity - which is how an edit to one of them can look
// unchanged and get dropped. Keys are sorted so that a value round-tripped through the API compares
// equal to one built locally regardless of key order.
export function samePropertyValue(left: {} | null | undefined, right: {} | null | undefined) {
	return stableStringify(left) === stableStringify(right);
}

function stableStringify(value: unknown): string {
	return JSON.stringify(value, (_key, nested) =>
		nested && typeof nested === 'object' && !Array.isArray(nested)
			? Object.fromEntries(Object.entries(nested).sort(([a], [b]) => a.localeCompare(b)))
			: nested
	);
}

export function displayJson(
	valueType: PropertyValueType | undefined,
	value: {} | null | undefined
) {
	if (value === null || value === undefined || value === '') return '—';

	try {
		switch (valueType) {
			case "price":
				const price = value as PTPrice
				return (price.amount / PriceMultiplier).toLocaleString('sr-RS', {
					style: 'currency',
					currency: price.currency,
					minimumFractionDigits: 2,
					maximumFractionDigits: 4
				});
			case "mass":
			case "volume": {
				const measure = value as PTMeasure;
				const amount = (measure.amount / MeasureMultiplier).toLocaleString('sr-RS', {
					maximumFractionDigits: 4
				});
				return `${amount}${measure.unit}`;
			}
			case "boolean":
				return value ? "Da" : "Ne";
			default:
				return typeof value === 'string' ? value : JSON.stringify(value);
		}
	} catch {
		return "???";
	}
}

export function defaultJsonValue(valueType: PropertyValueType, value: {} | null) {
	if (value !== null) return value;
	if (valueType === 'string') return "";
	if (valueType === 'number') return 0;
	if (valueType === 'price') return { amount: 0, currency: 'RSD' };
	if (valueType === 'mass') return { amount: 0, unit: 'g' };
	if (valueType === 'volume') return { amount: 0, unit: 'mL' };
	if (valueType === 'expiry') {
		let sample_date = new Date();
		sample_date.setDate(sample_date.getDate() + 7)
		return sample_date.toISOString().split('T')[0];
	}
	return false;
}
