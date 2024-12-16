<script>
	import Job from "./Job.svelte";
	import Add from "../Add.svelte";
	import Modal from "../Modal.svelte";
	import SvelteCronGen from "../vendor/Cron.svelte";
	import toast from "svelte-hot-french-toast";
	import { websocketState, addJob } from "$lib/websocket-bridge.svelte";

	let showModal = $state(false);
	let cron = $state("0 * * * * ?");

	function add() {
		showModal = true;
	}

	function createJob() {
		showModal = false;
		addJob(cron);
		toast.success("Job added");
	}
</script>

<svelte:head>
	<title>Jobs - Betterbell</title>
</svelte:head>

<p>
	Each job represents a different timer that will ring the bell. Create one by pressing the blue add button in the
	bottom right.
</p>
{#if websocketState.jobs.size > 0}
	<ul class="jobs">
		{#each websocketState.jobs as [id, job]}
			<Job title={job.title} toggled={job.toggled} {id} />
		{/each}
	</ul>
{/if}
<Add on:click={add} />
<Modal bind:showModal>
	<SvelteCronGen bind:value={cron} language="en" showSeconds={true} showAdvanced={false} minutesStep={1} />
	<button onclick={createJob} class="button-lg button-accent button">Create Job</button>
</Modal>

<style>
	ul {
		list-style-type: none;
		padding: 0;
		margin: 0;
		border-radius: 8px;
		box-shadow:
			0 0 0 1px RGB(0 0 0 / 3%),
			0 1px 3px 1px RGB(0 0 0 / 7%),
			0 2px 6px 2px RGB(0 0 0 / 3%);
	}
</style>
