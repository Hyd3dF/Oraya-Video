<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { getVideos } from '$lib/api/videos';
	import VideoCard from '$lib/components/VideoCard.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import type { Video } from '$lib/types';

	let videos = $state<Video[]>([]);
	let videosLoading = $state(true);

	let isOwner = $derived(auth.user?.id);

	onMount(async () => {
		if (!auth.initialized) {
			// Wait for auth to initialize
			const checkAuth = setInterval(() => {
				if (auth.initialized) {
					clearInterval(checkAuth);
					if (!auth.isLoggedIn) {
						goto('/login');
						return;
					}
					loadUserVideos();
				}
			}, 100);
			return;
		}
		if (!auth.isLoggedIn) {
			goto('/login');
			return;
		}
		loadUserVideos();
	});

	async function loadUserVideos() {
		videosLoading = true;
		try {
			// Fetch all videos and filter by owner - simple approach
			const data = await getVideos(100, 0);
			videos = (data.videos || []).filter((v) => v.owner_id === auth.user?.id);
		} catch {
			videos = [];
		} finally {
			videosLoading = false;
		}
	}

	let user = $derived(auth.user);
</script>

<div class="page-container">
	{#if !user}
		<LoadingSpinner />
	{:else}
		<div class="profile-header">
			<div class="avatar avatar-lg">
				{#if user.avatar_url}
					<img src={user.avatar_url} alt={user.username} />
				{:else}
					{user.username.charAt(0).toUpperCase()}
				{/if}
			</div>
			<div class="profile-info">
				<h1 class="profile-name">{user.display_name || user.real_name || user.username}</h1>
				<span class="profile-username">@{user.username}</span>
				{#if user.bio}
					<p class="profile-bio">{user.bio}</p>
				{/if}
				<div class="profile-meta">
					<span>{user.email}</span>
				</div>
				<a href="/profile/edit" class="btn btn-secondary btn-sm" style="margin-top:8px">Edit Profile</a>
			</div>
		</div>

		<h2 class="section-title">My Videos</h2>
		{#if videosLoading}
			<LoadingSpinner message="Loading your videos..." />
		{:else if videos.length === 0}
			<EmptyState icon="video" title="No videos yet" message="Upload your first video to get started." />
		{:else}
			<div class="video-grid">
				{#each videos as video (video.id)}
					<VideoCard {video} />
				{/each}
			</div>
		{/if}
	{/if}
</div>

<style>
	.profile-header {
		display: flex;
		gap: 24px;
		align-items: flex-start;
		padding: 24px;
		background: var(--bg-card);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		margin-bottom: 32px;
	}

	.profile-info {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.profile-name {
		font-size: 24px;
		font-weight: 700;
	}

	.profile-username {
		color: var(--text-muted);
		font-size: 14px;
	}

	.profile-bio {
		margin-top: 8px;
		font-size: 14px;
		color: var(--text-secondary);
		line-height: 1.5;
	}

	.profile-meta {
		margin-top: 12px;
		font-size: 13px;
		color: var(--text-muted);
	}

	.section-title {
		font-size: 18px;
		font-weight: 600;
		margin-bottom: 16px;
	}

	.video-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: 20px 14px;
	}

	@media (max-width: 640px) {
		.profile-header {
			flex-direction: column;
			align-items: center;
			text-align: center;
		}
		.video-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
