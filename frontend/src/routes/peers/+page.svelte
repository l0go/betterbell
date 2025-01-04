<script>
	import Add from "../Add.svelte";
	import Modal from "../Modal.svelte";
	import toast from "svelte-hot-french-toast";
	import { websocketState } from "$lib/websocket-bridge.svelte";

	let ip = $state("");
	let token = $state("");
	let showModal = $state(false);

	function add() {
		showModal = true;
	}

	function addPeer() {
		showModal = false;
		toast.success("Job added");
	}
</script>

<svelte:head>
	<title>Peers - Betterbell</title>
</svelte:head>

<p>
	Lorem
</p>
{#if websocketState.jobs.size > 0}
	<ul class="jobs">
		{#each websocketState.jobs as [id, job]}
		{/each}
	</ul>
{/if}
<Add on:click={add} />
<Modal bind:showModal>
	<div class="form">
		<label for="uname"><b>IP Address</b></label>
		<input type="text" placeholder="192.168.*.*:port" name="ip" bind:value={ip} required />

		<label for="psw"><b>Private Token</b></label>
		<input type="password" placeholder="base64 encoded token" name="token" bind:value={token} required />
		<input
			type="submit"
			value="Add Peer"
			onclick={addPeer}
			class="button button-lg button-accent"
		/>
	</div>
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

	.form {
		padding-top: 16px;
	}
</style>
