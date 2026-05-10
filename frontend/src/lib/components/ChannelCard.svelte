<script lang="ts">
	import type { Profile } from '$lib/types';

	let { channel }: { channel: Profile } = $props();

	function formatCount(n: number): string {
		if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
		if (n >= 1_000) return `${(n / 1_000).toFixed(1)}K`;
		return String(n);
	}
</script>

<a href="/channel/{channel.id}" class="channel-card">
	<div class="channel-avatar">
		{#if channel.avatar_url}
			<img src={channel.avatar_url} alt={channel.username} />
		{:else}
			{channel.username.charAt(0).toUpperCase()}
		{/if}
	</div>
	<div class="channel-info">
		<span class="channel-name">{channel.display_name || channel.username}</span>
		<span class="channel-username">@{channel.username}</span>
		{#if channel.bio}
			<p class="channel-bio">{channel.bio}</p>
		{/if}
	</div>
</a>

<style>
	.channel-card {
		display: flex;
		align-items: center;
		gap: 16px;
		padding: 16px;
		border-radius: var(--radius-lg);
		background: var(--bg-card);
		border: 1px solid var(--border);
		transition: border-color 0.15s;
	}

	.channel-card:hover {
		border-color: var(--accent);
	}

	.channel-avatar {
		width: 64px;
		height: 64px;
		border-radius: 50%;
		background: var(--accent-dim);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 24px;
		font-weight: 700;
		color: white;
		overflow: hidden;
		flex-shrink: 0;
	}

	.channel-avatar img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.channel-info {
		display: flex;
		flex-direction: column;
		gap: 2px;
		min-width: 0;
	}

	.channel-name {
		font-size: 16px;
		font-weight: 600;
	}

	.channel-username {
		font-size: 13px;
		color: var(--text-muted);
	}

	.channel-bio {
		font-size: 13px;
		color: var(--text-secondary);
		margin-top: 4px;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
