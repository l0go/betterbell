<script>
	let { toggled = $bindable(false), /** @type {{} => void} */ onToggled } = $props();
	const toggleIsPressed = (e) => {
		e.stopPropagation();
		toggled = !toggled;
		if (onToggled != null) {
			onToggled();
		}
	};
</script>

<button aria-label={toggled ? "On" : "Off"} aria-pressed={toggled} onclick={toggleIsPressed}>
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

	button[aria-pressed="true"] {
		background-color: var(--accent-bg);
		.slider {
			transform: translateX(20px);
		}
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
		transition: transform 200ms;
		transition-timing-function: cubic-bezier(0.5, 1, 0.89, 1);
	}
</style>
