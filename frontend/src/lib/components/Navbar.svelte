<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	let searchQuery = $state('');
	let mobileMenuOpen = $state(false);

	function handleSearch(e: Event) {
		e.preventDefault();
		const q = searchQuery.trim();
		if (q) {
			mobileMenuOpen = false;
			goto(`/search?q=${encodeURIComponent(q)}`);
		}
	}

	function handleLogout() {
		auth.logout();
		mobileMenuOpen = false;
		goto('/');
	}

	let currentPath = $derived($page.url.pathname);
</script>

<nav class="navbar">
	<div class="navbar-inner">
		<div class="navbar-left">
			<button class="mobile-menu-btn" onclick={() => (mobileMenuOpen = !mobileMenuOpen)}>
				<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					{#if mobileMenuOpen}
						<line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
					{:else}
						<line x1="4" y1="6" x2="20" y2="6" /><line x1="4" y1="12" x2="20" y2="12" /><line x1="4" y1="18" x2="20" y2="18" />
					{/if}
				</svg>
			</button>
			<a href="/" class="logo">oroya</a>
		</div>

		<form class="search-bar" onsubmit={handleSearch}>
			<input
				type="text"
				placeholder="Search videos and channels..."
				bind:value={searchQuery}
				class="search-input"
			/>
			<button type="submit" class="search-btn" aria-label="Search">
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<circle cx="11" cy="11" r="8" /><line x1="21" y1="21" x2="16.65" y2="16.65" />
				</svg>
			</button>
		</form>

		<div class="navbar-right" class:mobile-open={mobileMenuOpen}>
			<div class="nav-links">
				<a href="/" class:active={currentPath === '/'}>Home</a>
				{#if auth.isLoggedIn}
					<a href="/upload" class:active={currentPath === '/upload'}>Upload</a>
				{/if}
			</div>
			<div class="nav-auth">
				{#if auth.loading}
					<div class="skeleton" style="width: 32px; height: 32px; border-radius: 50%"></div>
				{:else if auth.isLoggedIn}
					<a href="/profile" class="nav-user" class:active={currentPath.startsWith('/profile')}>
						<div class="avatar" style="width:32px;height:32px;font-size:14px">
							{#if auth.user?.avatar_url}
								<img src={auth.user.avatar_url} alt={auth.user.username} />
							{:else}
								{auth.user?.username?.charAt(0)?.toUpperCase() || '?'}
							{/if}
						</div>
					</a>
					<button class="btn btn-ghost btn-sm logout-btn" onclick={handleLogout}>Logout</button>
				{:else}
					<a href="/login" class="btn btn-ghost btn-sm" class:active={currentPath === '/login'}>Login</a>
					<a href="/register" class="btn btn-primary btn-sm">Register</a>
				{/if}
			</div>
		</div>
	</div>
</nav>

<style>
	.navbar {
		position: sticky;
		top: 0;
		z-index: 100;
		background: var(--bg-primary);
		border-bottom: 1px solid var(--border);
		height: var(--navbar-height);
	}

	.navbar-inner {
		max-width: var(--max-width);
		margin: 0 auto;
		height: 100%;
		display: flex;
		align-items: center;
		gap: 16px;
		padding: 0 16px;
	}

	.navbar-left {
		display: flex;
		align-items: center;
		gap: 12px;
		flex-shrink: 0;
	}

	.mobile-menu-btn {
		display: none;
		padding: 4px;
		color: var(--text-primary);
	}

	.logo {
		font-size: 22px;
		font-weight: 700;
		color: var(--text-primary);
		letter-spacing: -0.5px;
	}

	.search-bar {
		flex: 1;
		max-width: 600px;
		display: flex;
		position: relative;
	}

	.search-input {
		width: 100%;
		padding: 8px 44px 8px 16px;
		background: var(--bg-secondary);
		border: 1px solid var(--border);
		border-radius: 999px;
		font-size: 14px;
	}

	.search-input:focus {
		border-color: var(--accent);
		background: var(--bg-primary);
	}

	.search-btn {
		position: absolute;
		right: 4px;
		top: 50%;
		transform: translateY(-50%);
		padding: 6px;
		color: var(--text-secondary);
		border-radius: 50%;
	}

	.search-btn:hover {
		color: var(--text-primary);
	}

	.navbar-right {
		display: flex;
		align-items: center;
		gap: 8px;
		flex-shrink: 0;
	}

	.nav-links {
		display: flex;
		gap: 4px;
	}

	.nav-links a,
	.logout-btn {
		padding: 6px 12px;
		border-radius: var(--radius);
		font-size: 14px;
		color: var(--text-secondary);
		transition: background 0.15s, color 0.15s;
	}

	.nav-links a:hover,
	.logout-btn:hover {
		background: var(--bg-hover);
		color: var(--text-primary);
	}

	.nav-links a.active {
		color: var(--text-primary);
		background: var(--bg-hover);
	}

	.nav-auth {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	@media (max-width: 768px) {
		.mobile-menu-btn {
			display: block;
		}

		.navbar-right {
			position: fixed;
			top: var(--navbar-height);
			left: 0;
			right: 0;
			bottom: 0;
			background: var(--bg-primary);
			flex-direction: column;
			align-items: stretch;
			padding: 16px;
			gap: 8px;
			display: none;
			z-index: 99;
		}

		.navbar-right.mobile-open {
			display: flex;
		}

		.nav-links {
			flex-direction: column;
		}

		.nav-links a {
			padding: 12px 16px;
			font-size: 16px;
		}

		.nav-auth {
			flex-direction: column;
			align-items: stretch;
			gap: 8px;
		}

		.nav-auth .btn {
			justify-content: center;
		}
	}

	@media (max-width: 480px) {
		.search-bar {
			max-width: none;
		}
		.logo {
			font-size: 18px;
		}
	}
</style>
