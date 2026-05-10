<script lang="ts">
	import { login } from '$lib/api/auth';
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { friendlyApiMessage } from '$lib/api/client';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		if (!email.trim() || !password.trim()) {
			error = 'Email and password are required';
			return;
		}
		loading = true;
		try {
			const data = await login(email.trim(), password);
			auth.setAuth(data);
			goto('/');
		} catch (e) {
			error = friendlyApiMessage(e, 'Login failed. Please try again.');
		} finally {
			loading = false;
		}
	}
</script>

<div class="page-container auth-page">
	<div class="auth-card">
		<h1 class="page-title">Login</h1>
		<form onsubmit={handleSubmit}>
			<div class="form-group">
				<label for="email">Email</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					placeholder="your@email.com"
					autocomplete="email"
					required
				/>
			</div>
			<div class="form-group">
				<label for="password">Password</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					placeholder="Enter your password"
					autocomplete="current-password"
					required
				/>
			</div>
			{#if error}
				<p class="form-error">{error}</p>
			{/if}
			<button type="submit" class="btn btn-primary" style="width:100%;margin-top:8px" disabled={loading}>
				{#if loading}Logging in...{:else}Login{/if}
			</button>
		</form>
		<p class="auth-footer">
			Don't have an account? <a href="/register">Register</a>
		</p>
	</div>
</div>

<style>
	.auth-page {
		display: flex;
		justify-content: center;
		padding-top: 80px;
	}

	.auth-card {
		width: 100%;
		max-width: 400px;
		background: var(--bg-card);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		padding: 32px;
	}

	.auth-card form {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.auth-footer {
		text-align: center;
		margin-top: 16px;
		font-size: 14px;
		color: var(--text-secondary);
	}

	.auth-footer a {
		color: var(--accent);
		font-weight: 500;
	}

	.auth-footer a:hover {
		color: var(--accent-hover);
	}
</style>
