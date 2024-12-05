<svelte:options runes="{true}" />

<script>
import Header from "./Header.svelte";
import Sidebar from "./Sidebar.svelte";
import "../app.css";
import {Toaster} from 'svelte-hot-french-toast'
let { children } = $props();
</script>

<div class="app">
	<main>
		<Sidebar />
		<div class="right">
			<Header />
			<div class="content">
				{@render children()}
			</div>
		</div>
	</main>
	<Toaster style="position:relative; z-index: 100000000;"/>
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
		padding-bottom: 16px;
	}
</style>
