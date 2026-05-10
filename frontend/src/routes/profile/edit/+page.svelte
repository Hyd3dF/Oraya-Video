<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { updateProfile } from '$lib/api/auth';
	import { ApiError } from '$lib/api/client';

	let displayName = $state('');
	let username = $state('');
	let realName = $state('');
	let bio = $state('');
	let loading = $state(false);
	let saving = $state(false);
	let error = $state('');
	let success = $state('');

	onMount(() => {
		if (!auth.isLoggedIn) {
			goto('/login');
			return;
		}
		const u = auth.user;
		if (u) {
			displayName = u.display_name || '';
			username = u.username || '';
			realName = u.real_name || '';
			bio = u.bio || '';
		}
		loading = false;
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		success = '';
		saving = true;

		try {
			const updated = await updateProfile({
				display_name: displayName || undefined,
				username: username || undefined,
				real_name: realName || undefined,
				bio: bio || undefined
			});
			// Update auth store
			auth.setAuth({
				access_token: '',
				refresh_token: '',
				expires_in: 0,
				token_type: 'bearer',
				user: updated
			} as any);
			success = 'Profile updated successfully';
		} catch (e) {
			if (e instanceof ApiError) {
				error = e.message;
			} else {
				error = 'Failed to update profile';
			}
		} finally {
			saving = false;
		}
	}
</script>

<div class="page-container">
	<div class="edit-card">
		<h1 class="page-title">Edit Profile</h1>

		{#if !auth.isLoggedIn}
			<p>Please log in to edit your profile.</p>
		{:else}
			<form onsubmit={handleSubmit}>
				<div class="form-group">
					<label for="displayName">Display Name</label>
					<input id="displayName" type="text" bind:value={displayName} placeholder="Display name" />
				</div>
				<div class="form-group">
					<label for="username">Username</label>
					<input id="username" type="text" bind:value={username} placeholder="Username" required />
				</div>
				<div class="form-group">
					<label for="realName">Full Name</label>
					<input id="realName" type="text" bind:value={realName} placeholder="Full name" required />
				</div>
				<div class="form-group">
					<label for="bio">Bio</label>
					<textarea id="bio" bind:value={bio} placeholder="Tell us about yourself..." rows="3"></textarea>
				</div>

				{#if error}
					<p class="form-error">{error}</p>
				{/if}
				{#if success}
					<p style="color:var(--success);font-size:13px">{success}</p>
				{/if}

				<div class="form-actions">
					<button type="submit" class="btn btn-primary" disabled={saving}>
						{#if saving}Saving...{:else}Save Changes{/if}
					</button>
					<a href="/profile" class="btn btn-ghost">Cancel</a>
				</div>
			</form>
		{/if}
	</div>
</div>

<style>
	.edit-card {
		max-width: 500px;
		margin: 0 auto;
		background: var(--bg-card);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		padding: 32px;
	}

	.edit-card form {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.form-actions {
		display: flex;
		gap: 12px;
		margin-top: 8px;
	}
</style>
