<script>
	import Header from "./Header.svelte";
	import { stopPropagation } from 'svelte/legacy';
	let { children, showModal = $bindable(false) } = $props();
	let dialog;
	$effect(() => {
		if (showModal) dialog.showModal(); else dialog.close();
	});
</script>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-noninteractive-element-interactions -->
<dialog
	bind:this={dialog}
	onclose={() => (showModal = false)}
	onclick={() => dialog.close()}
>
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div onclick="{stopPropagation()}">
		<Header closeButton="{true}" onclose={() => dialog.close()} />
		<div class="modal-content" onclick="{stopPropagation()}">
			{@render children()}
		</div>
	</div>
</dialog>

<style>
	dialog {
		width: 38em;
		border-radius: 8px;
		border: none;
		padding: 0;
	}
	dialog::backdrop {
		background: rgba(0, 0, 0, 0.3);
	}
	.modal-content {
		padding-left: 1em;
		padding-right: 1em;
		padding-bottom: 1em;
	}
	dialog[open] {
		animation: zoom 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
	}
	@keyframes zoom {
		from {
			transform: scale(0.95);
		}
		to {
			transform: scale(1);
		}
	}
	dialog[open]::backdrop {
		animation: fade 0.2s ease-out;
	}
	@keyframes fade {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}
	button {
		display: block;
	}
</style>
