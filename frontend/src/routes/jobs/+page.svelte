<script>
import Job from "./Job.svelte";
import Add from "../Add.svelte";
import Modal from "../Modal.svelte";
import SvelteCronGen from "../vendor/Cron.svelte"
import toast from "svelte-hot-french-toast"
import { connect, websocketState, addJob } from "$lib/websocket-bridge.svelte";
import {toString} from "cronstrue";

let showModal = $state(false);
let cron = $state('0 * * * * ?');
let translatedCron = $state(toString(cron));

function add() {
	showModal = !showModal;
}

function createJob() {
	showModal = false;
	addJob(cron);
	toast.success("Job added");
}

connect();
</script>

<svelte:head>
	<title>Jobs - Betterbell</title>
</svelte:head>

<p>Each job represents a different timer that will ring the bell. Create one by pressing the blue add button in the bottom right.</p>
{#if websocketState.jobs.length > 0}
	<ul class="jobs">
		{#each websocketState.jobs as job }
			<Job title={job.title} toggled={job.toggled} />
		{/each}
	</ul>
{/if}
<Add on:click={add} />
<Modal bind:showModal>
	<SvelteCronGen
		bind:value={cron}
		bind:valueTranslated={translatedCron}
		language="en"
		showSeconds={true}
		showAdvanced={false}
		minutesStep={1}
	/>
	<button onclick={createJob} class="create-job button">Create Job</button>
</Modal>

<style>
	ul {
		list-style-type: none;
		padding: 0;
		margin: 0;
		border-radius: 8px;
		box-shadow: 0 0 0 1px RGB(0 0 0 / 3%),
                0 1px 3px 1px RGB(0 0 0 / 7%),
                0 2px 6px 2px RGB(0 0 0 / 3%);
	}
	.create-job {
		background-color: var(--accent-bg);
		color: var(--bg);
		padding: 10px;
		width: 100%;
	}
	.create-job:hover {
		background-color: var(--accent-bg-hover);
	}
</style>
