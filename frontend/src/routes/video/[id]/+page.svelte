<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { getVideo, likeVideo, viewVideo, addVideoLink, deleteVideoLink } from '$lib/api/videos';
	import { getComments, addComment } from '$lib/api/comments';
	import { auth } from '$lib/stores/auth.svelte';
	import CommentItem from '$lib/components/CommentItem.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import ErrorMessage from '$lib/components/ErrorMessage.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import type { Video, VideoAsset, VideoLink, Comment } from '$lib/types';

	let videoId = $derived($page.params.id ?? '');
	let video = $state<Video | null>(null);
	let assets = $state<VideoAsset[]>([]);
	let links = $state<VideoLink[]>([]);
	let comments = $state<Comment[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	let likeCount = $state(0);
	let isLiked = $state(false);
	let likeLoading = $state(false);

	let commentText = $state('');
	let commentLoading = $state(false);
	let commentsLoading = $state(true);

	let newLinkTitle = $state('');
	let newLinkUrl = $state('');
	let linkLoading = $state(false);
	let showLinkForm = $state(false);

	let isOwner = $derived(auth.user?.id === video?.owner_id);

	onMount(async () => {
		if (!videoId) return;
		await loadVideo();
		// Record view
		if (video) {
			viewVideo(videoId).catch(() => {});
		}
	});

	async function loadVideo() {
		if (!videoId) return;
		loading = true;
		error = null;
		try {
			const data = await getVideo(videoId);
			video = data.video;
			assets = data.assets || [];
			links = data.links || [];
			likeCount = data.video.likes_count || 0;
			loadComments();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load video';
		} finally {
			loading = false;
		}
	}

	async function loadComments() {
		if (!videoId) return;
		commentsLoading = true;
		try {
			const data = await getComments(videoId);
			comments = data.comments || [];
		} catch {
			comments = [];
		} finally {
			commentsLoading = false;
		}
	}

	async function handleLike() {
		if (likeLoading || !auth.isLoggedIn || !videoId) return;
		likeLoading = true;
		try {
			const res = await likeVideo(videoId);
			isLiked = res.liked;
			likeCount += res.liked ? 1 : -1;
		} finally {
			likeLoading = false;
		}
	}

	async function handleAddComment(e: Event) {
		e.preventDefault();
		const text = commentText.trim();
		if (!text || commentLoading || !videoId) return;

		commentLoading = true;
		try {
			const newComment = await addComment(videoId, text);
			comments = [newComment, ...comments];
			commentText = '';
		} catch {
			// ignore
		} finally {
			commentLoading = false;
		}
	}

	function handleCommentDelete(commentId: string) {
		comments = comments.filter((c) => c.id !== commentId);
	}

	async function handleAddLink(e: Event) {
		e.preventDefault();
		if (!newLinkTitle.trim() || !newLinkUrl.trim() || linkLoading || !videoId) return;
		linkLoading = true;
		try {
			const newLink = await addVideoLink(videoId, {
				title: newLinkTitle.trim(),
				url: newLinkUrl.trim(),
				sort_order: links.length
			});
			links = [...links, newLink];
			newLinkTitle = '';
			newLinkUrl = '';
			showLinkForm = false;
		} catch {
			// ignore
		} finally {
			linkLoading = false;
		}
	}

	async function handleDeleteLink(linkId: string) {
		if (!videoId) return;
		try {
			await deleteVideoLink(videoId, linkId);
			links = links.filter((l) => l.id !== linkId);
		} catch {
			// ignore
		}
	}

	function formatDate(dateStr: string): string {
		const d = new Date(dateStr);
		return d.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
	}

	function formatCount(n: number): string {
		if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
		if (n >= 1_000) return `${(n / 1_000).toFixed(1)}K`;
		return String(n);
	}

	let hlsUrl = $derived(assets.find((a) => a.playlist_url)?.playlist_url || null);
</script>

<div class="page-container">
	{#if loading}
		<LoadingSpinner message="Loading video..." />
	{:else if error}
		<ErrorMessage message={error} onretry={loadVideo} />
	{:else if video}
		<div class="video-detail">
			<div class="player-area">
				{#if hlsUrl}
					<video controls class="video-player" src={hlsUrl} poster={video.thumbnail_url || ''}>
						<track kind="captions" />
						Your browser does not support the video tag.
					</video>
				{:else}
					<div class="player-placeholder">
						{#if video.thumbnail_url}
							<img src={video.thumbnail_url} alt={video.title} class="player-thumb" />
						{/if}
						<svg width="64" height="64" viewBox="0 0 24 24" fill="white" opacity="0.6">
							<polygon points="5 3 19 12 5 21 5 3" />
						</svg>
						{#if video.status && video.status !== 'ready'}
							<p class="player-status">
								{#if video.status === 'processing'}
									Video is still processing. It will be available shortly.
								{:else if video.status === 'failed'}
									Video processing failed. Please upload it again.
								{:else}
									Status: {video.status}
								{/if}
							</p>
						{:else}
							<p class="player-status">Video file is not ready yet.</p>
						{/if}
					</div>
				{/if}
			</div>

			<div class="video-info">
				<h1 class="video-title">{video.title}</h1>
				<div class="video-stats">
					<span>{formatCount(video.views_count)} views</span>
					<span class="separator">·</span>
					<span>{formatDate(video.created_at)}</span>
				</div>
				<div class="video-actions">
					<button class="btn btn-secondary btn-sm" onclick={handleLike} disabled={!auth.isLoggedIn || likeLoading}>
						<svg width="16" height="16" viewBox="0 0 24 24" fill={isLiked ? 'currentColor' : 'none'} stroke="currentColor" stroke-width="2">
							<path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3H14z" />
						</svg>
						{formatCount(likeCount)}
					</button>
				</div>

				{#if video.description}
					<div class="video-description">
						<p>{video.description}</p>
					</div>
				{/if}
			</div>

			<hr class="divider" />

			<!-- Links Section -->
			<div class="section">
				<div class="section-header">
					<h2>Links</h2>
					{#if isOwner}
						<button class="btn btn-ghost btn-sm" onclick={() => (showLinkForm = !showLinkForm)}>
							{#if showLinkForm}Cancel{:else}+ Add Link{/if}
						</button>
					{/if}
				</div>

				{#if showLinkForm}
					<form class="link-form" onsubmit={handleAddLink}>
						<input type="text" bind:value={newLinkTitle} placeholder="Link title" required />
						<input type="url" bind:value={newLinkUrl} placeholder="https://..." required />
						<button type="submit" class="btn btn-primary btn-sm" disabled={linkLoading}>Add</button>
					</form>
				{/if}

				{#if links.length > 0}
					<div class="links-list">
						{#each links as link (link.id)}
							<div class="link-item">
								<a href={link.url} target="_blank" rel="noopener noreferrer" class="link-title">
									{link.title}
								</a>
								{#if isOwner}
									<button class="btn btn-ghost btn-sm" onclick={() => handleDeleteLink(link.id)} aria-label="Delete link">
										<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
											<line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
										</svg>
									</button>
								{/if}
							</div>
						{/each}
					</div>
				{:else}
					<p class="text-muted">No links added yet.</p>
				{/if}
			</div>

			<hr class="divider" />

			<!-- Comments Section -->
			<div class="section">
				<h2>Comments ({comments.length})</h2>

				{#if auth.isLoggedIn}
					<form class="comment-form" onsubmit={handleAddComment}>
						<textarea bind:value={commentText} placeholder="Add a comment..." rows="2" required></textarea>
						<div class="comment-form-actions">
							<button type="submit" class="btn btn-primary btn-sm" disabled={commentLoading || !commentText.trim()}>
								{#if commentLoading}Posting...{:else}Comment{/if}
							</button>
						</div>
					</form>
				{:else}
					<p class="text-muted"><a href="/login">Login</a> to leave a comment.</p>
				{/if}

				{#if commentsLoading}
					<LoadingSpinner message="Loading comments..." />
				{:else if comments.length === 0}
					<EmptyState icon="comment" title="No comments yet" message="Be the first to comment." />
				{:else}
					<div class="comments-list">
						{#each comments as comment (comment.id)}
							<CommentItem {comment} ondelete={() => handleCommentDelete(comment.id)} />
						{/each}
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	.video-detail {
		max-width: 900px;
		margin: 0 auto;
	}

	.player-area {
		width: 100%;
		aspect-ratio: 16 / 9;
		background: black;
		border-radius: var(--radius-lg);
		overflow: hidden;
		margin-bottom: 20px;
	}

	.video-player {
		width: 100%;
		height: 100%;
	}

	.player-placeholder {
		width: 100%;
		height: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		position: relative;
		background: var(--bg-secondary);
	}

	.player-thumb {
		position: absolute;
		inset: 0;
		width: 100%;
		height: 100%;
		object-fit: cover;
		opacity: 0.4;
	}

	.player-status {
		margin-top: 12px;
		color: var(--warning);
		font-size: 14px;
		text-transform: capitalize;
	}

	.video-info {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.video-title {
		font-size: 20px;
		font-weight: 700;
		line-height: 1.3;
	}

	.video-stats {
		display: flex;
		align-items: center;
		gap: 6px;
		font-size: 13px;
		color: var(--text-muted);
	}

	.separator {
		font-weight: 700;
	}

	.video-actions {
		display: flex;
		gap: 8px;
	}

	.video-description {
		background: var(--bg-card);
		border-radius: var(--radius);
		padding: 16px;
		font-size: 14px;
		line-height: 1.6;
		white-space: pre-wrap;
	}

	.divider {
		border: none;
		border-top: 1px solid var(--border);
		margin: 24px 0;
	}

	.section {
		margin-bottom: 24px;
	}

	.section h2 {
		font-size: 16px;
		font-weight: 600;
		margin-bottom: 12px;
	}

	.section-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 12px;
	}

	.section-header h2 {
		margin-bottom: 0;
	}

	.link-form {
		display: flex;
		gap: 8px;
		margin-bottom: 12px;
	}

	.link-form input {
		flex: 1;
	}

	.links-list {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.link-item {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 8px 12px;
		background: var(--bg-card);
		border-radius: var(--radius);
		border: 1px solid var(--border);
	}

	.link-title {
		flex: 1;
		color: var(--accent);
		font-size: 14px;
		font-weight: 500;
	}

	.link-title:hover {
		color: var(--accent-hover);
	}

	.text-muted {
		font-size: 14px;
		color: var(--text-muted);
	}

	.text-muted a {
		color: var(--accent);
		font-weight: 500;
	}

	.comment-form {
		margin-bottom: 16px;
	}

	.comment-form textarea {
		width: 100%;
		resize: vertical;
	}

	.comment-form-actions {
		display: flex;
		justify-content: flex-end;
		margin-top: 8px;
	}

	@media (max-width: 640px) {
		.video-title {
			font-size: 17px;
		}
		.link-form {
			flex-direction: column;
		}
	}
</style>
