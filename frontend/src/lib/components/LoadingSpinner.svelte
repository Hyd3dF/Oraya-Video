<script lang="ts">
	let { message = 'Loading...' }: { message?: string } = $props();
</script>

<div class="loading">
	<div class="spinner"></div>
	<p>{message}</p>
</div>

<style>
	.loading {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 16px;
		padding: 48px 16px;
		color: var(--text-secondary);
	}

	.spinner {
		width: 36px;
		height: 36px;
		border: 3px solid var(--border);
		border-top-color: var(--accent);
		border-radius: 50%;
		animation: spin 0.7s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.loading p {
		font-size: 14px;
	}
</style>
