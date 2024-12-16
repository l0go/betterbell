<script>
	import Switch from "../Switch.svelte";
	import right from "$lib/icons/right-smaller-symbolic.svg";
	import { deleteJob, toggleJob, updateJob, websocketState } from "$lib/websocket-bridge.svelte";
	import Modal from "../Modal.svelte";
	import SvelteCronGen from "../vendor/Cron.svelte";
	import toast from "svelte-hot-french-toast";

	let { title = "Job Name", toggled = false, id = 0 } = $props();
	let showModal = $state(false);
	let expression = $state(websocketState.jobs.get(id).expression);

	function onToggled() {
		toggleJob(id, toggled);
	}

	function onPress() {
		showModal = true;
	}

	function onUpdateJob() {
		updateJob(id, expression);
		showModal = false;
		toast.success("Updated job");
	}

	function onDeleteJob() {
		deleteJob(id);
		showModal = false;
		toast.success("Deleted job");
	}
</script>

<button class="job" onclick={onPress}>
	<p class="title">{title}</p>
	<div>
		<Switch bind:toggled {onToggled} />
		<img src={right} alt="Right arrow" />
	</div>
</button>

<Modal bind:showModal>
	<SvelteCronGen
		bind:value={expression}
		bind:valueTranslated={title}
		language="en"
		showSeconds={true}
		showAdvanced={false}
		minutesStep={1}
	/>
	<div style="display: flex;">
		<button onclick={onDeleteJob} class="button-lg button">Delete Job</button>
		<button onclick={onUpdateJob} class="button-lg button-accent button">Update Job</button>
	</div>
</Modal>

<style>
	.job {
		display: flex;
		width: 100%;
		border: none;
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
