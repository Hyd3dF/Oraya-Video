<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { getChannel, subscribeToChannel } from '$lib/api/channels';
	import { auth } from '$lib/stores/auth.svelte';
	import VideoCard from '$lib/components/VideoCard.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import ErrorMessage from '$lib/components/ErrorMessage.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import type { Profile, Video } from '$lib/types';

	let channelId = $derived($page.params.id ?? '');
	let profile = $state<Profile | null>(null);
	let videos = $state<Video[]>([]);
	let subscriberCount = $state(0);
	let isSubscribed = $state(false);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let subLoading = $state(false);

	let isOwnChannel = $derived(auth.user?.id === channelId);

	onMount(() => {
		if (channelId) loadChannel();
	});

	async function loadChannel() {
		if (!channelId) return;
		loading = true;
		error = null;
		try {
			const data = await getChannel(channelId);
			profile = data.profile;
			videos = data.videos || [];
			subscriberCount = data.subscriber_count || 0;
			isSubscribed = data.is_subscribed || false;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load channel';
		} finally {
			loading = false;
		}
	}

	async function handleSubscribe() {
		if (subLoading || !auth.isLoggedIn || !channelId) return;
		subLoading = true;
		try {
			const res = await subscribeToChannel(channelId);
			isSubscribed = res.subscribed;
			subscriberCount += res.subscribed ? 1 : -1;
		} catch {
			// ignore
		} finally {
			subLoading = false;
		}
	}

	function formatCount(n: number): string {
		if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
		if (n >= 1_000) return `${(n / 1_000).toFixed(1)}K`;
		return String(n);
	}
</script>

<div class="page-container">
	{#if loading}
		<LoadingSpinner message="Loading channel..." />
	{:else if error}
		<ErrorMessage message={error} onretry={loadChannel} />
	{:else if profile}
		<div class="channel-header">
			<div class="avatar avatar-lg">
				{#if profile.avatar_url}
					<img src={profile.avatar_url} alt={profile.username} />
				{:else}
					{profile.username.charAt(0).toUpperCase()}
				{/if}
			</div>
			<div class="channel-info">
				<h1 class="channel-name">{profile.display_name || profile.real_name || profile.username}</h1>
				<span class="channel-username">@{profile.username}</span>
				<div class="channel-meta">
					<span>{formatCount(subscriberCount)} subscribers</span>
					<span class="separator">·</span>
					<span>{videos.length} videos</span>
				</div>
				{#if profile.bio}
					<p class="channel-bio">{profile.bio}</p>
				{/if}

				{#if !isOwnChannel && auth.isLoggedIn}
					<button class="btn btn-primary btn-sm" onclick={handleSubscribe} disabled={subLoading} style="margin-top:12px">
						{#if subLoading}
							...
						{:else if isSubscribed}
							Subscribed
						{:else}
							Subscribe
						{/if}
					</button>
				{:else if isOwnChannel}
					<a href="/profile" class="btn btn-secondary btn-sm" style="margin-top:12px">Edit Profile</a>
				{/if}
			</div>
		</div>

		<h2 class="section-title">Videos</h2>
		{#if videos.length === 0}
			<EmptyState icon="video" title="No videos yet" message="This channel hasn't uploaded any videos." />
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
	.channel-header {
		display: flex;
		gap: 24px;
		align-items: flex-start;
		padding: 24px;
		background: var(--bg-card);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		margin-bottom: 32px;
	}

	.channel-info {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.channel-name {
		font-size: 24px;
		font-weight: 700;
	}

	.channel-username {
		color: var(--text-muted);
		font-size: 14px;
	}

	.channel-meta {
		display: flex;
		align-items: center;
		gap: 6px;
		font-size: 13px;
		color: var(--text-secondary);
		margin-top: 4px;
	}

	.separator {
		font-weight: 700;
	}

	.channel-bio {
		margin-top: 8px;
		font-size: 14px;
		color: var(--text-secondary);
		line-height: 1.5;
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
		.channel-header {
			flex-direction: column;
			align-items: center;
			text-align: center;
		}
		.video-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
