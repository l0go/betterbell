<script>
import { page } from "$app/stores";
import menu_large from "$lib/icons/menu-large-symbolic.svg";
import bell_outline from "$lib/icons/bell-outline-symbolic.svg";
import chain_link from "$lib/icons/chain-link-symbolic.svg";
</script>

<nav>
	<div class="sidebar-top">
		<a class="button toggle-sidebar"><img src={menu_large} width="14px" alt="Open Sidebar" /></a>
		<p class="title">Betterbell</p>
	</div>
	<div class="sidebar-content">
		<a class="button" aria-current={$page.url.pathname === '/jobs' ? 'page' : undefined} href="/jobs">
			<img src={bell_outline} alt="Bell icon" />
			Jobs
		</a>
		<a class="button" aria-current={$page.url.pathname === '/peers' ? 'page' : undefined} href="/peers">
			<img src={chain_link} alt="Chain icon" />
			Peers
		</a>
	</div>
</nav>

<style>
	nav {
		z-index: 10000;
		width: 300px;
		background-color: var(--header);
		border-right: 1px solid color-mix(in srgb, var(--header-shadow) 50%, transparent);
	}

	.sidebar-content {
		display: flex;
		flex-direction: column;
		gap: 4px;
		padding: 8px;
	}

	.sidebar-content .button {
		justify-content: start;
		gap: 8px;
		height: 36px;
		font-weight: 700;
		padding-left: 12px;
		padding-right: 12px;
	}

	.toggle-sidebar {
		width: 34px;
		height: 34px;
	}

	.sidebar-top {
		display: grid;
		grid-template-columns: 32px 1fr 32px;
		min-height: 47px;
		padding: 6px 7px 7px 7px;
		box-shadow: 0 1px color-mix(in srgb, var(--header-shadow) 50%, transparent),
                0 2px 4px color-mix(in srgb, var(--header-shadow) 50%, transparent);
	}
	.sidebar-top .title {
		display: flex;
		justify-content: center; /* Align horizontal */
		align-items: center;
	}
</style>
