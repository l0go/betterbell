<script>
	import Add from "../Add.svelte";
	import Modal from "../Modal.svelte";
	import toast from "svelte-hot-french-toast";
	import { websocketState, addPeer } from "$lib/websocket-bridge.svelte";
	import Peer from "./Peer.svelte";

	let ip = $state("");
	let token = $state("");
	let showModal = $state(false);

	function add() {
		showModal = true;
	}

	function createPeer() {
		if (ip.length <= 0 || token.length <= 0) {
			toast.success("Length of IP or Token must exceed 0");
			return;
		}
		showModal = false;
		addPeer(ip, token);
		toast.success("Job added");
	}
</script>

<svelte:head>
	<title>Peers - Betterbell</title>
</svelte:head>

{#if websocketState.jobs.size > 0}
	{#each websocketState.peers as [_, peer]}
		{#if peer.address == "BETTERBELL__SELF"}
			<p>This instance's private token is <span class="secret">{peer.token}</span> (hover to reveal)</p>
		{/if}
	{/each}

	<ul class="jobs">
		{#each websocketState.peers as [id, peer]}
			{#if peer.address != "BETTERBELL__SELF"}
				<Peer id={id} title={peer.address} />
			{/if}
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
			onclick={createPeer}
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

	.secret {
		background-color: black;
		color: black;
	}

	.secret:hover {
		background-color: white;
	}
</style>
