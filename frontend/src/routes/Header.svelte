<script>
import { page } from "$app/stores";
import cross_small from "$lib/icons/cross-small-symbolic.svg"; 
let n = $derived($page?.url?.pathname);
let pagename = $state("");
$effect(() => {
	pagename = $page?.url?.pathname?.charAt(1).toUpperCase() + $page?.url?.pathname?.slice(2);
});
let { closeButton = false, onclose} = $props();
</script>

<header>
	{#if closeButton}
	<div />
	{/if}
	<p class="title">{pagename}</p>
	{#if closeButton}
		<div id="button-box">
			<button onclick="{onclose}" id="close-button"><img src={cross_small} alt="x" /></button>
		</div>
	{/if}
</header>

<style>
	header {
		display: flex;
		position: relative;
		width: 100%;
		align-items: center; 
		min-height: 47px;
		padding: 6px 7px 7px 7px;
		background-color: var(--bg);
		box-shadow: 0 1px color-mix(in srgb, var(--header-shadow) 50%, transparent),
                0 2px 4px color-mix(in srgb, var(--header-shadow) 50%, transparent);
	}
	header > * {
		flex: 1;
	}
	.title {
		text-align: center;
	}
	#button-box {
		display: flex;
		justify-content: end;
	}
	#close-button {
		border: none;
		display: grid;
  align-items: center;
		border-radius: 30px;
		width: 24px;
		height: 24px;
	}
</style>
