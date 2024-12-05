<script>
import Switch from "../Switch.svelte";
import right from "$lib/icons/right-smaller-symbolic.svg";
let { title = "Job Name", toggled = false } = $props();
</script>

<li class="job">
	<p class="title">{title}</p>
	<div>
		<Switch toggled={toggled} />
		<img src={right} alt="Right arrow" />
	</div>
</li>

<style>
	.job {
		display: flex;
		justify-content: space-between;
		align-items: center;
		background-color: var(--bg);
		padding: 12px;
		padding-top: 16px;
		padding-bottom: 16px;
		border-bottom: 1px solid var(--header-shadow);
		transition: 200ms ease-out;
	}
	.job:hover {
		background-color: var(--header);
	}

	div {
		display: flex;
		justify-content: space-between;
		gap: 10px;
	}
	:global(.job:nth-child(1)) {
		border-top-left-radius: 8px;
		border-top-right-radius: 8px;
	}
	:global(.job:last-child) {
		border-bottom: none;
		border-bottom-left-radius: 8px;
		border-bottom-right-radius: 8px;
	}
</style>
