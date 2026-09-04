import type { PropertyOption } from '$lib/api';

export function requiredTextError(value: string, label: string) {
	return value.trim() ? undefined : `Unesite ${label.toLocaleLowerCase()}.`;
}

export function emailError(value: string) {
	if (!value.trim()) return 'Unesite email adresu.';
	return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value) ? undefined : 'Unesite ispravnu email adresu.';
}

export function passwordError(value: string) {
	if (!value) return 'Unesite lozinku.';
	return value.length >= 8 && /[A-Z]/.test(value) && /[a-z]/.test(value) && /\d/.test(value)
		? undefined
		: 'Lozinka mora imati najmanje 8 znakova, veliko i malo slovo, i broj.';
}

function isValidDate(value: string) {
	if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) return false;

	const date = new Date(`${value}T00:00:00.000Z`);
	return !Number.isNaN(date.valueOf()) && date.toISOString().slice(0, 10) === value;
}

// Keep this aligned with backend/internal/dto/validate.go. The backend remains the authority;
// this gives the user immediate, field-specific feedback before the request is sent.
export function propertyValueError(property: Pick<PropertyOption, 'value_type'>, value: unknown) {
	switch (property.value_type) {
		case 'boolean':
			return typeof value === 'boolean' ? undefined : 'Unesite vrednost.';
		case 'string':
			return typeof value === 'string' && value.trim() ? undefined : 'Unesite vrednost.';
		case 'expiry':
			if (typeof value !== 'string' || !value.trim()) return 'Unesite datum.';
			return isValidDate(value) ? undefined : 'Unesite ispravan datum.';
		case 'number':
			return typeof value === 'number' && Number.isFinite(value)
				? undefined
				: 'Unesite ispravan broj.';
	}

	if (typeof value !== 'object' || value === null) return 'Unesite vrednost i jedinicu.';
	const structuredValue = value as { amount?: unknown; currency?: unknown; unit?: unknown };
	if (
		typeof structuredValue.amount !== 'number' ||
		!Number.isSafeInteger(structuredValue.amount) ||
		structuredValue.amount < 0
	) {
		return 'Unesite nenegativnu vrednost.';
	}
	if (property.value_type === 'price') {
		return typeof structuredValue.currency === 'string' && structuredValue.currency
			? undefined
			: 'Odaberite valutu.';
	}
	return typeof structuredValue.unit === 'string' && structuredValue.unit
		? undefined
		: 'Odaberite jedinicu mere.';
}
