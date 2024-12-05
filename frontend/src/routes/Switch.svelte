<script>
let { toggled = false} = $props();
const toggleIsPressed = () => {
	toggled = !toggled;
};
</script>

<button aria-pressed={toggled} onclick={toggleIsPressed}>
	<span class="slider"></span>
</button>

<style>
	button {
		border: none;
		padding: 3px;
		min-width: 46px;
		min-height: 23px;
		border-radius: 14px;
		background-color: var(--trough-color);
	}

	button:hover {
		cursor: pointer;
	}

	button[aria-pressed=true] {
		background-color: var(--accent-bg);
	}
	button[aria-pressed=true] .slider {
		transform: translateX(20px);
	}

	.slider {
		display: block;
		padding: 0;
		margin: 0;
		background-color: var(--bg);
		border-radius: 20px;
		width: 20px;
	    height: 20px;
		box-shadow: 0 2px 4px RGB(0 0 0 / 20%);
	}
</style>
