import { apiPost, apiGet, apiPut, saveTokens, clearTokens } from './client';
import type { AuthResponse, Profile } from '$lib/types';

export async function register(
	email: string,
	password: string,
	real_name: string,
	username: string
): Promise<AuthResponse> {
	const data = await apiPost<AuthResponse>('/api/v1/auth/register', {
		email,
		password,
		real_name,
		username
	});
	saveTokens(data.access_token, data.refresh_token);
	return data;
}

export async function login(email: string, password: string): Promise<AuthResponse> {
	const data = await apiPost<AuthResponse>('/api/v1/auth/login', { email, password });
	saveTokens(data.access_token, data.refresh_token);
	return data;
}

export async function refreshToken(): Promise<AuthResponse | null> {
	// Handled automatically by client.ts; exposed for manual use
	return null;
}

export async function logout(): Promise<void> {
	try {
		await apiPost('/api/v1/auth/logout');
	} catch {
		// Ignore errors, clear locally anyway
	}
	clearTokens();
}

export async function getMe(): Promise<Profile> {
	return apiGet<Profile>('/api/v1/auth/me');
}

export async function updateProfile(data: Partial<Profile>): Promise<Profile> {
	return apiPut<Profile>('/api/v1/me', data);
}
