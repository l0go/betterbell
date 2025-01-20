<script>
	import { deletePeer } from "$lib/websocket-bridge.svelte";
	import toast from "svelte-hot-french-toast";

	let { title = "Peer", id = 0 } = $props();
	function onDeletePeer() {
		deletePeer(id);
		toast.success("Deleted Peer");
	}
</script>

<div class="card">
	<p class="title">{title}</p>
	<div>
		<button class="button button-lg button-accent" onclick={onDeletePeer}>Delete</button>
	</div>
</div>
