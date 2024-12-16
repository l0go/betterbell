<script>
	import Header from "./Header.svelte";
	let { children, showModal = $bindable(false) } = $props();

	/** @type {HTMLDialogElement} */
	let dialog;
	$effect(() => {
		if (showModal) dialog.showModal();
		else dialog.close();
	});

	/** @param {MouseEvent} e */
	function onContentClick(e) {
		e.stopPropagation();
	}
</script>

<div role="presentation" onclick={() => dialog.close()}>
	<dialog bind:this={dialog} onclose={() => (showModal = false)}>
		<div onclick={onContentClick} role="presentation">
			<Header closeButton={true} onclose={() => dialog.close()} />
			<div class="modal-content" onclick={onContentClick} role="presentation">
				{@render children()}
			</div>
		</div>
	</dialog>
</div>

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
</style>
