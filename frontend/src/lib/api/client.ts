export class ApiError extends Error {
	constructor(
		message: string,
		readonly status: number
	) {
		super(message);
	}
}

const apiErrorMessages: Record<string, string> = {
	'invalid credentials': 'Email adresa ili lozinka nisu ispravni.',
	'User status not active': 'Vaš zahtev za registraciju još nije odobren.',
	'user with this email already exists': 'Korisnik sa ovom email adresom već postoji.',
	'item already approved to another user': 'Stavka je već odobrena drugom korisniku.',
	'not approved to consume item': 'Nemate odobrenje za korišćenje ove stavke.',
	'property is used by a derived name format': 'Svojstvo se koristi u formatu izvedenog naziva.',
	'invalid derived name format': 'Format izvedenog naziva nije ispravan.',
	'already exists': 'Podaci sa ovim vrednostima već postoje.',
	'invalid reference': 'Izabrana povezana vrednost ne postoji.',
	'not found': 'Traženi podatak nije pronađen.'
};

function userFacingErrorMessage(status: number, serverMessage: string | undefined) {
	if (serverMessage && apiErrorMessages[serverMessage]) return apiErrorMessages[serverMessage];

	switch (status) {
		case 400:
			return 'Poslati podaci nisu ispravni.';
		case 401:
			return 'Prijava je istekla. Prijavite se ponovo.';
		case 403:
			return 'Nemate dozvolu za ovu radnju.';
		case 404:
			return 'Traženi podatak nije pronađen.';
		case 409:
			return 'Podaci sa ovim vrednostima već postoje.';
		default:
			return 'Došlo je do greške na serveru.';
	}
}

export async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
	const response = await fetch(`/api${path}`, {
		credentials: 'include',
		...init
	});

	const contentType = response.headers.get('content-type');
	const body = contentType?.includes('application/json') ? await response.json() : undefined;

	if (!response.ok) {
		const serverMessage =
			typeof body === 'object' && body !== null && 'error' in body && typeof body.error === 'string'
				? body.error
				: undefined;
		throw new ApiError(userFacingErrorMessage(response.status, serverMessage), response.status);
	}

	return body as T;
}

export function jsonRequest(
	method: 'DELETE' | 'PATCH' | 'POST' | 'PUT',
	body: unknown
): RequestInit {
	return {
		method,
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	};
}
