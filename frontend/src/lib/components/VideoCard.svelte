<script lang="ts">
	import type { Video } from '$lib/types';

	let { video }: { video: Video } = $props();

	let thumbnailUrl = $derived(video.thumbnail_url);
	let durationDisplay = $derived(formatDuration(video.duration_seconds));
	let viewDisplay = $derived(formatCount(video.views_count));
	let likeDisplay = $derived(formatCount(video.likes_count));

	function formatDuration(seconds: number): string {
		if (!seconds || seconds <= 0) return '';
		const h = Math.floor(seconds / 3600);
		const m = Math.floor((seconds % 3600) / 60);
		const s = seconds % 60;
		if (h > 0) return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
		return `${m}:${String(s).padStart(2, '0')}`;
	}

	function formatCount(n: number): string {
		if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
		if (n >= 1_000) return `${(n / 1_000).toFixed(1)}K`;
		return String(n);
	}
</script>

<a href={`/video/${video.id}`} class="video-card">
	<div class="thumbnail">
		{#if thumbnailUrl}
			<img src={thumbnailUrl} alt={video.title} loading="lazy" />
		{:else}
			<div class="thumbnail-placeholder">
				<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.4">
					<polygon points="5 3 19 12 5 21 5 3" />
				</svg>
			</div>
		{/if}
		{#if durationDisplay}
			<span class="duration">{durationDisplay}</span>
		{/if}
		{#if video.status && video.status !== 'ready'}
			<span class="status-badge">{video.status}</span>
		{/if}
	</div>
	<div class="info">
		<h3 class="title">{video.title}</h3>
		<div class="meta">
			<span>{viewDisplay} views</span>
			{#if video.likes_count > 0}
				<span class="separator">·</span>
				<span>{likeDisplay} likes</span>
			{/if}
		</div>
	</div>
</a>

<style>
	.video-card {
		display: flex;
		flex-direction: column;
		border-radius: var(--radius-lg);
		overflow: hidden;
		transition: transform 0.2s, box-shadow 0.2s;
	}

	.video-card:hover {
		transform: translateY(-2px);
	}

	.thumbnail {
		position: relative;
		width: 100%;
		aspect-ratio: 16 / 9;
		background: var(--bg-hover);
		overflow: hidden;
	}

	.thumbnail img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.thumbnail-placeholder {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--text-muted);
		background: linear-gradient(135deg, var(--bg-hover), var(--bg-card));
	}

	.duration {
		position: absolute;
		bottom: 6px;
		right: 6px;
		background: rgba(0, 0, 0, 0.85);
		color: white;
		padding: 2px 6px;
		border-radius: var(--radius-sm);
		font-size: 12px;
		font-weight: 500;
	}

	.status-badge {
		position: absolute;
		top: 6px;
		left: 6px;
		background: rgba(245, 158, 11, 0.9);
		color: black;
		padding: 2px 8px;
		border-radius: var(--radius-sm);
		font-size: 11px;
		font-weight: 600;
		text-transform: uppercase;
	}

	.info {
		padding: 12px 4px;
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.title {
		font-size: 15px;
		font-weight: 600;
		line-height: 1.3;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.meta {
		display: flex;
		align-items: center;
		gap: 6px;
		font-size: 13px;
		color: var(--text-secondary);
	}

	.separator {
		font-weight: 700;
	}
</style>
