import { browser } from '$app/environment';
import { getMe, logout as apiLogout } from '$lib/api/auth';
import { saveTokens, clearTokens, tryRefreshToken } from '$lib/api/client';
import type { Profile, AuthResponse } from '$lib/types';

function createAuthStore() {
	let user = $state<Profile | null>(null);
	let loading = $state(true);
	let initialized = $state(false);
	let error = $state<string | null>(null);

	async function init() {
		if (!browser) return;
		if (initialized || !loading) return;
		const token = localStorage.getItem('access_token');
		if (!token) {
			loading = false;
			initialized = true;
			return;
		}
		try {
			const profile = await getMe();
			user = profile;
		} catch {
			const refreshed = await tryRefreshToken();
			if (refreshed) {
				try {
					user = await getMe();
				} catch {
					clearTokens();
					user = null;
				}
			} else {
				clearTokens();
				user = null;
			}
		} finally {
			loading = false;
			initialized = true;
		}
	}

	function setAuth(data: AuthResponse) {
		if (data.access_token && data.refresh_token) {
			saveTokens(data.access_token, data.refresh_token);
		}
		user = data.user;
		error = null;
	}

	async function logout() {
		await apiLogout();
		user = null;
		error = null;
	}

	function clearError() {
		error = null;
	}

	return {
		get user() {
			return user;
		},
		get loading() {
			return loading;
		},
		get initialized() {
			return initialized;
		},
		get isLoggedIn() {
			return !!user;
		},
		get error() {
			return error;
		},
		set error(val: string | null) {
			error = val;
		},
		init,
		setAuth,
		logout,
		clearError
	};
}

export const auth = createAuthStore();
