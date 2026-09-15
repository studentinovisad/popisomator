import type { DashboardConsumptionBucket, DashboardMonths } from '$lib/api';
import { consumptionLabel } from '$lib/domain/items';

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

const monthFormatter = new Intl.DateTimeFormat('sr-RS', { month: 'short' });
const monthYearFormatter = new Intl.DateTimeFormat('sr-RS', { month: 'long', year: 'numeric' });
const dateFormatter = new Intl.DateTimeFormat('sr-RS', { dateStyle: 'medium' });

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

// Serbian counts three ways, and "stavka" happens to need all three.
export function itemCountLabel(count: number) {
	const lastTwo = count % 100;
	const last = count % 10;
	if (last === 1 && lastTwo !== 11) return `${count} stavka`;
	if (last >= 2 && last <= 4 && (lastTwo < 12 || lastTwo > 14)) return `${count} stavke`;
	return `${count} stavki`;
}
