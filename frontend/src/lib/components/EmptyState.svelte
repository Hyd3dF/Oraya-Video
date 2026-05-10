<script lang="ts">
	let {
		icon = 'empty',
		title = 'Nothing here yet',
		message = ''
	}: {
		icon?: string;
		title?: string;
		message?: string;
	} = $props();
</script>

<div class="empty">
	{#if icon === 'search'}
		<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.4">
			<circle cx="11" cy="11" r="8" /><line x1="21" y1="21" x2="16.65" y2="16.65" />
		</svg>
	{:else if icon === 'video'}
		<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.4">
			<rect x="2" y="2" width="20" height="20" rx="2.18" ry="2.18" />
			<line x1="7" y1="2" x2="7" y2="22" /><line x1="17" y1="2" x2="17" y2="22" />
			<line x1="2" y1="12" x2="22" y2="12" /><line x1="2" y1="7" x2="7" y2="7" /><line x1="2" y1="17" x2="7" y2="17" />
			<line x1="17" y1="7" x2="22" y2="7" /><line x1="17" y1="17" x2="22" y2="17" />
		</svg>
	{:else if icon === 'comment'}
		<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.4">
			<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
		</svg>
	{:else}
		<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.4">
			<path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" />
		</svg>
	{/if}
	<h3>{title}</h3>
	{#if message}
		<p>{message}</p>
	{/if}
</div>

<style>
	.empty {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 12px;
		padding: 48px 16px;
		text-align: center;
		color: var(--text-muted);
	}

	.empty h3 {
		font-size: 18px;
		font-weight: 600;
		color: var(--text-secondary);
	}

	.empty p {
		font-size: 14px;
		max-width: 400px;
	}
</style>
