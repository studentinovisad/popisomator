export class ApiError extends Error {
	constructor(
		message: string,
		readonly status: number
	) {
		super(message);
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
		const message =
			typeof body === 'object' && body !== null && 'error' in body && typeof body.error === 'string'
				? body.error
				: 'Request failed';
		throw new ApiError(message, response.status);
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
