<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { browser } from '$app/environment';
	import Navbar from '$lib/components/Navbar.svelte';
	import '../app.css';

	let { children } = $props();

	$effect(() => {
		if (browser && !auth.initialized) {
			auth.init();
		}
	});
</script>

<Navbar />

<main>
	{@render children()}
</main>
