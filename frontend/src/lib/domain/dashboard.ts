import type {
	DashboardConsumptionBucket,
	DashboardGroupConsumptionQuantity,
	DashboardMonths,
	DashboardTypeConsumptionQuantity,
	ItemPropertyTotal,
	PropertyValueType,
	PTMeasure,
	PTPrice
} from '$lib/api';
import {
	consumptionLabel,
	displayJson,
	massUnitFactors,
	MeasureMultiplier,
	PriceMultiplier,
	volumeUnitFactors
} from '$lib/domain/items';

export const dashboardMonthRanges: { value: DashboardMonths; label: string }[] = [
	{ value: 3, label: 'Tromesečje' },
	{ value: 6, label: 'Pola godine' },
	{ value: 12, label: 'Godinu dana' }
];

export const defaultDashboardMonths: DashboardMonths = 6;

// The backend rejects anything outside the three ranges, so a value arriving from the URL is checked
// here rather than sent on to be refused.
export function parseDashboardMonths(value: string): DashboardMonths {
	const parsed = Number.parseInt(value, 10);
	const range = dashboardMonthRanges.find((option) => option.value === parsed);
	return range?.value ?? defaultDashboardMonths;
}

// The range as it reads inside a sentence, which is not the button's label: "u narednih Tromesečje"
// is not Serbian. Twelve months is said as a year, and three takes a different case from six.
export function monthRangeLabel(months: DashboardMonths) {
	if (months === 12) return 'godinu dana';
	return months === 3 ? '3 meseca' : `${months} meseci`;
}

// The three series of the consumption chart, coloured from the app's own tokens so the chart reads
// as part of the page. The wording of each status is shared, not restated here.
export const consumptionSeries: {
	key: keyof Omit<DashboardConsumptionBucket, 'month'>;
	color: string;
}[] = [
	{ key: 'fully_consumed', color: 'var(--color-brand)' },
	{ key: 'partially_consumed', color: 'var(--color-warning)' },
	{ key: 'damaged', color: 'var(--color-danger)' }
];

export function consumptionSeriesLabel(key: keyof Omit<DashboardConsumptionBucket, 'month'>) {
	return consumptionLabel(key);
}

const monthFormatter = new Intl.DateTimeFormat('sr-Latn-RS', { month: 'short' });
const monthYearFormatter = new Intl.DateTimeFormat('sr-Latn-RS', {
	month: 'long',
	year: 'numeric'
});
const dateFormatter = new Intl.DateTimeFormat('sr-Latn-RS', { dateStyle: 'medium' });

// Dates arrive as plain YYYY-MM-DD and are read as UTC, since a local reading would step one near
// midnight into the neighbouring day.
function parseMonth(month: string) {
	return new Date(`${month}T00:00:00Z`);
}

export function formatExpiryDate(date: string) {
	return dateFormatter.format(parseMonth(date));
}

// How long is left, or how long it has been. The count is already worked out against the server's
// clock, so this only has to word it.
export function expiryDelayLabel(daysRemaining: number) {
	if (daysRemaining < 0) return `istekao pre ${dayCountLabel(-daysRemaining)}`;
	if (daysRemaining === 0) return 'ističe danas';
	if (daysRemaining === 1) return 'ističe sutra';
	return `za ${dayCountLabel(daysRemaining)}`;
}

function dayCountLabel(days: number) {
	const isSingular = days % 10 === 1 && days % 100 !== 11;
	return isSingular ? `${days} dan` : `${days} dana`;
}

// Ticks are tight on space, so the year shows only where it changes rather than on all twelve.
export function formatMonthTick(month: string, previousMonth?: string) {
	const date = parseMonth(month);
	const label = monthFormatter.format(date).replace(/\.$/, '');
	if (previousMonth && parseMonth(previousMonth).getUTCFullYear() === date.getUTCFullYear()) {
		return label;
	}
	return `${label} ${date.getUTCFullYear()}`;
}

export function formatMonthLong(month: string) {
	return monthYearFormatter.format(parseMonth(month));
}

export function formatPropertyTotal(total: ItemPropertyTotal) {
	return displayJson(total.value_type, total.value);
}

export function formatPropertyTotals(totals: ItemPropertyTotal[]) {
	return totals.map((total) => `${total.property_name}: ${formatPropertyTotal(total)}`).join(' · ');
}

export type ConsumptionGroupRow = {
	key: string;
	label: string;
	periodTotal: ItemPropertyTotal;
	cells: (ItemPropertyTotal | undefined)[];
};

// A group almost always tracks a single property within one table (a liquid in litres, a solid in
// grams, or its cost), so the row is named after the substance. Only the rare group whose filtered
// properties still number more than one - two different measures on the same substance - needs the
// property named too, or the rows would look like duplicates.
export function consumptionGroupRows(
	group: DashboardGroupConsumptionQuantity,
	valueTypes: PropertyValueType[]
): ConsumptionGroupRow[] {
	const periodTotals = group.period_totals.filter((total) => valueTypes.includes(total.value_type));

	return periodTotals.map((periodTotal) => ({
		key: `${group.name}-${periodTotal.property_id}`,
		label: periodTotals.length > 1 ? `${group.name} (${periodTotal.property_name})` : group.name,
		periodTotal,
		cells: group.buckets.map((bucket) => {
			const total = bucket.totals.find(
				(candidate) => candidate.property_id === periodTotal.property_id
			);
			return total && total.value_count > 0 ? total : undefined;
		})
	}));
}

export type ConsumptionRowsFooter = {
	cells: ItemPropertyTotal[];
	total: ItemPropertyTotal;
};

export type ConsumptionRowsFooterGroup = {
	label: string;
	footer: ConsumptionRowsFooter;
};

function scaleToLargestUnit(baseAmount: number, unitFactors: Record<string, number>) {
	const sorted = Object.entries(unitFactors).sort(([, left], [, right]) => right - left);
	for (const [unit, factor] of sorted) {
		if (baseAmount >= MeasureMultiplier * factor && baseAmount % factor === 0) {
			return { amount: baseAmount / factor, unit };
		}
	}
	const [baseUnit, baseFactor] = sorted[sorted.length - 1];
	return { amount: baseAmount / baseFactor, unit: baseUnit };
}

function sumRowValues(
	totals: (ItemPropertyTotal | undefined)[],
	valueType: PropertyValueType,
	currency: string
): ItemPropertyTotal {
	const baseAmount = totals.reduce((sum, total) => sum + (total ? totalToBaseUnit(total) : 0), 0);

	const value =
		valueType === 'price'
			? { amount: baseAmount, currency }
			: scaleToLargestUnit(baseAmount, valueType === 'mass' ? massUnitFactors : volumeUnitFactors);

	return {
		property_id: -1,
		property_name: '',
		value_type: valueType,
		value,
		value_count: baseAmount > 0 ? 1 : 0
	};
}

export function sumRowsFooters(rows: ConsumptionGroupRow[]): ConsumptionRowsFooterGroup[] {
	const rowsByValueType = new Map<PropertyValueType, ConsumptionGroupRow[]>();
	for (const row of rows) {
		const grouped = rowsByValueType.get(row.periodTotal.value_type) ?? [];
		grouped.push(row);
		rowsByValueType.set(row.periodTotal.value_type, grouped);
	}

	const groups = [...rowsByValueType.entries()];

	return groups.map(([valueType, groupRows]) => {
		const currency =
			valueType === 'price' ? (groupRows[0].periodTotal.value as PTPrice).currency : '';
		const monthCount = groupRows[0].cells.length;
		const cells = Array.from({ length: monthCount }, (_, monthIndex) =>
			sumRowValues(
				groupRows.map((row) => row.cells[monthIndex]),
				valueType,
				currency
			)
		);
		const total = sumRowValues(
			groupRows.map((row) => row.periodTotal),
			valueType,
			currency
		);

		return {
			label: groups.length > 1 ? `Ukupno (${groupRows[0].periodTotal.property_name})` : 'Ukupno',
			footer: { cells, total }
		};
	});
}

function totalToBaseUnit(total: ItemPropertyTotal): number {
	if (total.value_type === 'price') return (total.value as PTPrice).amount;
	if (total.value_type === 'mass') {
		const measure = total.value as PTMeasure;
		return measure.amount * (massUnitFactors[measure.unit] ?? 1);
	}
	if (total.value_type === 'volume') {
		const measure = total.value as PTMeasure;
		return measure.amount * (volumeUnitFactors[measure.unit] ?? 1);
	}
	return 0;
}

export function sortGroupsByTotal(
	groups: DashboardGroupConsumptionQuantity[],
	valueTypes: PropertyValueType[]
): DashboardGroupConsumptionQuantity[] {
	const totalFor = (group: DashboardGroupConsumptionQuantity) =>
		group.period_totals
			.filter((total) => valueTypes.includes(total.value_type))
			.reduce((sum, total) => sum + totalToBaseUnit(total), 0);

	return [...groups].sort((left, right) => totalFor(right) - totalFor(left));
}

export type MonthlySpend = {
	month: string;
	amount: number;
	currency: string;
};

export function monthlySpendTotals(types: DashboardTypeConsumptionQuantity[]): MonthlySpend[] {
	const totalsByMonth = new Map<string, MonthlySpend>();
	for (const month of types[0]?.groups[0]?.buckets.map((bucket) => bucket.month) ?? []) {
		totalsByMonth.set(month, { month, amount: 0, currency: '' });
	}

	for (const type of types) {
		for (const group of type.groups) {
			for (const bucket of group.buckets) {
				for (const total of bucket.totals) {
					if (total.value_type !== 'price' || total.value_count === 0) continue;

					const price = total.value as PTPrice;
					const existing = totalsByMonth.get(bucket.month) ?? {
						month: bucket.month,
						amount: 0,
						currency: ''
					};
					existing.amount += price.amount / PriceMultiplier;
					if (!existing.currency) existing.currency = price.currency;
					totalsByMonth.set(bucket.month, existing);
				}
			}
		}
	}

	return [...totalsByMonth.values()].sort((left, right) => left.month.localeCompare(right.month));
}

export type SpendTotal = {
	amount: number;
	currency: string;
};

export function totalSpend(spend: MonthlySpend[]): SpendTotal {
	return {
		amount: spend.reduce((sum, entry) => sum + entry.amount, 0),
		currency: spend.find((entry) => entry.currency)?.currency ?? 'RSD'
	};
}

export function formatSpendTotal(total: SpendTotal) {
	return total.amount.toLocaleString('sr-Latn-RS', {
		style: 'currency',
		currency: total.currency,
		minimumFractionDigits: 2,
		maximumFractionDigits: 2
	});
}

// Serbian counts three ways, and "stavka" happens to need all three.
export function itemCountLabel(count: number) {
	const lastTwo = count % 100;
	const last = count % 10;
	if (last === 1 && lastTwo !== 11) return `${count} stavka`;
	if (last >= 2 && last <= 4 && (lastTwo < 12 || lastTwo > 14)) return `${count} stavke`;
	return `${count} stavki`;
}
