import { browser } from '$app/environment';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

class ApiError extends Error {
	status: number;
	code: string;

	constructor(status: number, code: string, message: string) {
		super(message);
		this.name = 'ApiError';
		this.status = status;
		this.code = code;
	}
}

function getToken(): string | null {
	if (!browser) return null;
	return localStorage.getItem('access_token');
}

function getRefreshToken(): string | null {
	if (!browser) return null;
	return localStorage.getItem('refresh_token');
}

export function saveTokens(access: string, refresh: string) {
	if (!browser) return;
	localStorage.setItem('access_token', access);
	localStorage.setItem('refresh_token', refresh);
}

export function clearTokens() {
	if (!browser) return;
	localStorage.removeItem('access_token');
	localStorage.removeItem('refresh_token');
}

let isRefreshing = false;
let refreshPromise: Promise<boolean> | null = null;

async function tryRefreshToken(): Promise<boolean> {
	if (isRefreshing && refreshPromise) return refreshPromise;

	isRefreshing = true;
	refreshPromise = (async () => {
		try {
			const refresh = getRefreshToken();
			if (!refresh) return false;

			const res = await fetch(`${PUBLIC_API_BASE_URL}/api/v1/auth/refresh`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ refresh_token: refresh })
			});

			if (!res.ok) {
				clearTokens();
				return false;
			}

			const data = await res.json();
			saveTokens(data.access_token, data.refresh_token);
			return true;
		} catch {
			clearTokens();
			return false;
		} finally {
			isRefreshing = false;
			refreshPromise = null;
		}
	})();

	return refreshPromise;
}

type RequestOptions = {
	method?: string;
	body?: unknown;
	headers?: Record<string, string>;
	rawBody?: boolean;
};

export async function apiRequest<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
	const { method = 'GET', body, headers: extraHeaders, rawBody = false } = options;

	const headers: Record<string, string> = {
		...extraHeaders
	};

	if (!rawBody) {
		headers['Content-Type'] = 'application/json';
	}

	const token = getToken();
	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	const res = await fetch(`${PUBLIC_API_BASE_URL}${endpoint}`, {
		method,
		headers,
		body: body ? (rawBody ? body as BodyInit : JSON.stringify(body)) : undefined
	});

	if (res.status === 401 && token) {
		const refreshed = await tryRefreshToken();
		if (refreshed) {
			const newToken = getToken();
			if (newToken) {
				headers['Authorization'] = `Bearer ${newToken}`;
			}
			const retryRes = await fetch(`${PUBLIC_API_BASE_URL}${endpoint}`, {
				method,
				headers,
				body: body ? (rawBody ? body as BodyInit : JSON.stringify(body)) : undefined
			});
			if (!retryRes.ok) {
				const err = await retryRes.json().catch(() => ({ error: 'unknown', message: 'Request failed' }));
				throw new ApiError(retryRes.status, err.error || 'unknown', err.message || 'Request failed');
			}
			if (retryRes.status === 204) return undefined as T;
			return retryRes.json();
		}
	}

	if (!res.ok) {
		const err = await res.json().catch(() => ({ error: 'network_error', message: 'Network error' }));
		throw new ApiError(res.status, err.error || 'unknown', err.message || 'Request failed');
	}

	if (res.status === 204) return undefined as T;
	return res.json();
}

export function friendlyApiMessage(error: unknown, fallback = 'Something went wrong. Please try again.'): string {
	if (!(error instanceof ApiError)) return fallback;
	if (error.code === 'username_taken') {
		return 'This username is already taken.';
	}
	if (error.code === 'invalid_input' || error.code === 'bad_request') {
		return 'Please check the form and try again.';
	}
	if (error.status === 401 || error.status === 403) {
		return 'We could not complete this action. Please check your account details and try again.';
	}
	if (error.status === 404) {
		return 'The requested item could not be found.';
	}
	if (error.status >= 500) {
		return 'The server is temporarily unavailable. Please try again shortly.';
	}
	return fallback;
}

export function apiGet<T>(endpoint: string): Promise<T> {
	return apiRequest<T>(endpoint);
}

export function apiPost<T>(endpoint: string, body?: unknown): Promise<T> {
	return apiRequest<T>(endpoint, { method: 'POST', body });
}

export function apiPut<T>(endpoint: string, body?: unknown): Promise<T> {
	return apiRequest<T>(endpoint, { method: 'PUT', body });
}

export function apiDelete<T>(endpoint: string): Promise<T> {
	return apiRequest<T>(endpoint, { method: 'DELETE' });
}

export function apiUpload(endpoint: string, body: BodyInit): Promise<Response> {
	return fetch(`${PUBLIC_API_BASE_URL}${endpoint}`, {
		method: 'PUT',
		body
	});
}

export { ApiError };
