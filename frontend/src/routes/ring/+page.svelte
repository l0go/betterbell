<script>
	import bell_unpressed from "$lib/icons/bell-unpressed.svg";
	import bell_pressed from "$lib/icons/bell-pressed.svg";
	import { ring } from "$lib/websocket-bridge.svelte.js";

	let bell = bell_unpressed;
	function onClick() {
		ring();
	}

	function onDown() {
		bell = bell_pressed;
	}
	function onUp() {
		setTimeout(() => {
			bell = bell_unpressed;
		}, 120);
	}
</script>

<svelte:head>
	<title>Ring - Betterbell</title>
</svelte:head>

<div class="outer">
	<button onclick={onClick} onmousedown={onDown} onmouseup={onUp}>
		<img src={bell} alt="Button with the word RING written on top" />
	</button>
</div>

<style>
	.outer {
		display: flex;
		width: 100%;
		height: 100%;
		align-items: center;
		justify-content: center;
	}
	button {
		border: none;
		background: none;
	}
</style>
