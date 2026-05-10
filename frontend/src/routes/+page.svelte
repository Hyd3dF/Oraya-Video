<script lang="ts">
	import { onMount } from 'svelte';
	import { getVideos } from '$lib/api/videos';
	import VideoCard from '$lib/components/VideoCard.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import ErrorMessage from '$lib/components/ErrorMessage.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import type { Video } from '$lib/types';

	let videos = $state<Video[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let offset = $state(0);
	let hasMore = $state(true);
	let loadMoreLoading = $state(false);

	async function loadVideos(reset = false) {
		if (reset) {
			offset = 0;
			loading = true;
		} else {
			loadMoreLoading = true;
		}
		error = null;

		try {
			const data = await getVideos(24, reset ? 0 : offset);
			const newVideos = data.videos || [];
			if (reset) {
				videos = newVideos;
			} else {
				videos = [...videos, ...newVideos];
			}
			hasMore = newVideos.length === 24;
			if (reset) offset = newVideos.length;
			else offset += newVideos.length;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load videos';
		} finally {
			loading = false;
			loadMoreLoading = false;
		}
	}

	onMount(() => {
		loadVideos(true);
	});
</script>

<div class="page-container">
	<h1 class="page-title">Videos</h1>

	{#if loading}
		<LoadingSpinner message="Loading videos..." />
	{:else if error}
		<ErrorMessage message={error} onretry={() => loadVideos(true)} />
	{:else if videos.length === 0}
		<EmptyState icon="video" title="No videos yet" message="Videos will appear here once they are uploaded." />
	{:else}
		<div class="video-grid">
			{#each videos as video (video.id)}
				<VideoCard {video} />
			{/each}
		</div>

		{#if hasMore}
			<div class="load-more">
				<button class="btn btn-secondary" onclick={() => loadVideos(false)} disabled={loadMoreLoading}>
					{#if loadMoreLoading}
						Loading more...
					{:else}
						Load more
					{/if}
				</button>
			</div>
		{/if}
	{/if}
</div>

<style>
	.video-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 24px 16px;
	}

	.load-more {
		display: flex;
		justify-content: center;
		margin-top: 32px;
	}

	@media (max-width: 640px) {
		.video-grid {
			grid-template-columns: 1fr;
			gap: 20px;
		}
	}
</style>
