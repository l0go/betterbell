<script>
	import { createEventDispatcher } from "svelte";
	import plus from "$lib/icons/plus-symbolic.svg";
	const dispatch = createEventDispatcher();
	function click() {
		dispatch("click", "a");
	}
</script>

<button on:click={click}>
	<img src={plus} alt="Add" />
</button>

<style>
	button {
		display: flex;
		justify-content: center;
		align-items: center;

		border: none;
		border-radius: 48px;
		font-size: 24pt;
		line-height: 0;
		background-color: var(--accent-bg);

		box-shadow:
			0 0 0 1px RGB(0 0 0 / 3%),
			0 1px 3px 1px RGB(0 0 0 / 7%),
			0 2px 6px 2px RGB(0 0 0 / 3%);

		position: fixed;
		bottom: 12px;
		right: 12px;
		width: 48px;
		height: 48px;
	}
	button:hover {
		background-color: var(--accent-bg-hover);
	}
	button img {
		width: 24px;
		filter: brightness(0) invert(1);
	}
</style>
