<script lang="ts">
	import { register } from '$lib/api/auth';
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { friendlyApiMessage } from '$lib/api/client';

	let email = $state('');
	let username = $state('');
	let realName = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';

		if (!email.trim() || !username.trim() || !realName.trim()) {
			error = 'All fields are required';
			return;
		}

		if (password.length < 8) {
			error = 'Password must be at least 8 characters';
			return;
		}

		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		loading = true;
		try {
			const data = await register(email.trim(), password, realName.trim(), username.trim());
			auth.setAuth(data);
			goto('/');
		} catch (e) {
			error = friendlyApiMessage(e, 'Registration failed. Please try again.');
		} finally {
			loading = false;
		}
	}
</script>

<div class="page-container auth-page">
	<div class="auth-card">
		<h1 class="page-title">Register</h1>
		<form onsubmit={handleSubmit}>
			<div class="form-group">
				<label for="email">Email</label>
				<input id="email" type="email" bind:value={email} placeholder="your@email.com" autocomplete="email" required />
			</div>
			<div class="form-group">
				<label for="username">Username</label>
				<input id="username" type="text" bind:value={username} placeholder="username" autocomplete="username" required />
			</div>
			<div class="form-group">
				<label for="realName">Full Name</label>
				<input id="realName" type="text" bind:value={realName} placeholder="Your full name" required />
			</div>
			<div class="form-group">
				<label for="password">Password</label>
				<input id="password" type="password" bind:value={password} placeholder="Min 8 characters" autocomplete="new-password" required />
			</div>
			<div class="form-group">
				<label for="confirmPassword">Confirm Password</label>
				<input id="confirmPassword" type="password" bind:value={confirmPassword} placeholder="Repeat your password" autocomplete="new-password" required />
			</div>
			{#if error}
				<p class="form-error">{error}</p>
			{/if}
			<button type="submit" class="btn btn-primary" style="width:100%;margin-top:8px" disabled={loading}>
				{#if loading}Creating account...{:else}Register{/if}
			</button>
		</form>
		<p class="auth-footer">
			Already have an account? <a href="/login">Login</a>
		</p>
	</div>
</div>

<style>
	.auth-page {
		display: flex;
		justify-content: center;
		padding-top: 60px;
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
