<script lang="ts">
	import type { Comment } from '$lib/types';
	import { likeComment, deleteComment } from '$lib/api/comments';
	import { auth } from '$lib/stores/auth.svelte';

	let {
		comment,
		ondelete = () => {}
	}: {
		comment: Comment;
		ondelete?: () => void;
	} = $props();

	let likesCount = $state(comment.likes_count);
	let isLiked = $state(false);
	let isDeleting = $state(false);
	let likeLoading = $state(false);

	let isOwner = $derived(auth.user?.id === comment.user_id);

	function formatDate(dateStr: string): string {
		const d = new Date(dateStr);
		const now = new Date();
		const diff = now.getTime() - d.getTime();
		const mins = Math.floor(diff / 60000);
		const hrs = Math.floor(diff / 3600000);
		const days = Math.floor(diff / 86400000);
		if (mins < 1) return 'Just now';
		if (mins < 60) return `${mins}m ago`;
		if (hrs < 24) return `${hrs}h ago`;
		if (days < 7) return `${days}d ago`;
		return d.toLocaleDateString();
	}

	async function handleLike() {
		if (likeLoading || !auth.isLoggedIn) return;
		likeLoading = true;
		try {
			const res = await likeComment(comment.id);
			isLiked = res.liked;
			likesCount += res.liked ? 1 : -1;
		} catch {
			// ignore
		} finally {
			likeLoading = false;
		}
	}

	async function handleDelete() {
		if (isDeleting) return;
		isDeleting = true;
		try {
			await deleteComment(comment.id);
			ondelete();
		} catch {
			isDeleting = false;
		}
	}
</script>

<div class="comment">
	<div class="comment-avatar">
		<div class="avatar" style="width:32px;height:32px;font-size:14px">
			{#if comment.user?.avatar_url}
				<img src={comment.user.avatar_url} alt={comment.user.username} />
			{:else}
				{comment.user?.username?.charAt(0)?.toUpperCase() || '?'}
			{/if}
		</div>
	</div>
	<div class="comment-body">
		<div class="comment-header">
			<span class="comment-author">{comment.user?.display_name || comment.user?.username || 'Unknown'}</span>
			<span class="comment-time">{formatDate(comment.created_at)}</span>
		</div>
		<p class="comment-content">{comment.content}</p>
		<div class="comment-actions">
			<button class="action-btn" onclick={handleLike} disabled={!auth.isLoggedIn || likeLoading}>
				<svg width="16" height="16" viewBox="0 0 24 24" fill={isLiked ? 'currentColor' : 'none'} stroke="currentColor" stroke-width="2">
					<path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3H14z" />
				</svg>
				{#if likesCount > 0}
					<span>{likesCount}</span>
				{/if}
			</button>
			{#if isOwner}
				<button class="action-btn danger" onclick={handleDelete} disabled={isDeleting}>
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<polyline points="3 6 5 6 21 6" />
						<path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
					</svg>
					{#if isDeleting}...{:else}Delete{/if}
				</button>
			{/if}
		</div>
	</div>
</div>

<style>
	.comment {
		display: flex;
		gap: 12px;
		padding: 12px 0;
	}

	.comment + .comment {
		border-top: 1px solid var(--border);
	}

	.comment-body {
		flex: 1;
		min-width: 0;
	}

	.comment-header {
		display: flex;
		align-items: center;
		gap: 8px;
		margin-bottom: 4px;
	}

	.comment-author {
		font-size: 13px;
		font-weight: 600;
	}

	.comment-time {
		font-size: 12px;
		color: var(--text-muted);
	}

	.comment-content {
		font-size: 14px;
		line-height: 1.5;
		color: var(--text-primary);
		white-space: pre-wrap;
		word-break: break-word;
	}

	.comment-actions {
		display: flex;
		align-items: center;
		gap: 8px;
		margin-top: 8px;
	}

	.action-btn {
		display: inline-flex;
		align-items: center;
		gap: 4px;
		padding: 4px 8px;
		border-radius: var(--radius);
		font-size: 12px;
		color: var(--text-secondary);
		transition: background 0.15s, color 0.15s;
	}

	.action-btn:hover:not(:disabled) {
		background: var(--bg-hover);
		color: var(--text-primary);
	}

	.action-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.action-btn.danger:hover:not(:disabled) {
		color: var(--danger);
	}
</style>
