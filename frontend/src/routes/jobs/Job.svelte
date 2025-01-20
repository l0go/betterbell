<script>
	import Switch from "../Switch.svelte";
	import right from "$lib/icons/right-smaller-symbolic.svg";
	import { deleteJob, toggleJob, updateJob, websocketState } from "$lib/websocket-bridge.svelte";
	import Modal from "../Modal.svelte";
	import SvelteCronGen from "../vendor/Cron.svelte";
	import toast from "svelte-hot-french-toast";

	let { title = "Job Name", toggled = websocketState.jobs.get(id).toggled, id = 0 } = $props();
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

<button class="card" onclick={onPress}>
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
