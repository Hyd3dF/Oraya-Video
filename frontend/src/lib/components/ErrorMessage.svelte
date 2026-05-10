<script lang="ts">
	let {
		message = 'Something went wrong',
		onretry = undefined
	}: {
		message?: string;
		onretry?: (() => void) | undefined;
	} = $props();
</script>

<div class="error">
	<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="var(--danger)" stroke-width="1.5">
		<circle cx="12" cy="12" r="10" />
		<line x1="15" y1="9" x2="9" y2="15" />
		<line x1="9" y1="9" x2="15" y2="15" />
	</svg>
	<p>{message}</p>
	{#if onretry}
		<button class="btn btn-secondary" onclick={onretry}>Try again</button>
	{/if}
</div>

<style>
	.error {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 16px;
		padding: 48px 16px;
		text-align: center;
	}

	.error p {
		font-size: 15px;
		color: var(--text-secondary);
		max-width: 400px;
	}
</style>
