<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { search } from '$lib/api/channels';
	import VideoCard from '$lib/components/VideoCard.svelte';
	import ChannelCard from '$lib/components/ChannelCard.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import ErrorMessage from '$lib/components/ErrorMessage.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import type { Video, Profile } from '$lib/types';

	let query = $state('');
	let videos = $state<Video[]>([]);
	let channels = $state<Profile[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	$effect(() => {
		const q = $page.url.searchParams.get('q') || '';
		if (q) {
			query = q;
			doSearch(q);
		}
	});

	async function doSearch(q: string) {
		loading = true;
		error = null;
		try {
			const results = await search(q);
			videos = results.videos || [];
			channels = results.channels || [];
		} catch (e) {
			error = e instanceof Error ? e.message : 'Search failed';
			videos = [];
			channels = [];
		} finally {
			loading = false;
		}
	}

	let hasResults = $derived(videos.length > 0 || channels.length > 0);
</script>

<div class="page-container">
	<h1 class="page-title">Search results for "{query}"</h1>

	{#if loading}
		<LoadingSpinner message="Searching..." />
	{:else if error}
		<ErrorMessage message={error} onretry={() => doSearch(query)} />
	{:else if !hasResults}
		<EmptyState icon="search" title="No results found" message="Try a different search term." />
	{:else}
		{#if channels.length > 0}
			<div class="search-section">
				<h2 class="section-title">Channels</h2>
				<div class="channels-list">
					{#each channels as channel (channel.id)}
						<ChannelCard {channel} />
					{/each}
				</div>
			</div>
		{/if}

		{#if videos.length > 0}
			<div class="search-section">
				<h2 class="section-title">Videos</h2>
				<div class="video-grid">
					{#each videos as video (video.id)}
						<VideoCard {video} />
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</div>

<style>
	.search-section {
		margin-bottom: 32px;
	}

	.section-title {
		font-size: 18px;
		font-weight: 600;
		margin-bottom: 16px;
	}

	.channels-list {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.video-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 24px 16px;
	}

	@media (max-width: 640px) {
		.video-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
