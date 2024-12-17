<svelte:options runes={true} />

<script>
	import Header from "./Header.svelte";
	import Sidebar from "./Sidebar.svelte";
	import "../app.css";
	import toast, { Toaster } from "svelte-hot-french-toast";
	import { connect, websocketState } from "$lib/websocket-bridge.svelte";
	import { goto } from "$app/navigation";
	import { browser } from "$app/environment";
	let { children } = $props();

	connect();
	if (browser && !websocketState.authenticated) {
		goto("/login");
	}
</script>

<div class="app">
	<main>
		{#if websocketState.authenticated}
			<Sidebar />
		{/if}
		<div class="right">
			<Header />
			<div class="content">
				{@render children()}
			</div>
		</div>
	</main>
	<Toaster style="position:relative; z-index: 100000000;" />
</div>

<style>
	main {
		display: flex;
		width: 100%;
		height: 100vh;
		flex-direction: row;
	}

	.right {
		width: 100%;
	}
	.content {
		overflow-y: auto;
		height: calc(100% - 47px);
		padding-left: 32px;
		padding-right: 32px;
	}
</style>
