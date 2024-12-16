<script>
	import { goto } from "$app/navigation";
	import { loginStandard, websocketState } from "$lib/websocket-bridge.svelte";

	let username = $state("");
	let password = $state("");

	$effect(() => {
		if (websocketState.authenticated) {
			goto("/jobs");
		}
	});

	function onLogin() {
		loginStandard(username, password);
	}
</script>

<svelte:head>
	<title>Login - Betterbell</title>
</svelte:head>

<div class="center-content">
	<div class="form">
		<label for="uname"><b>Username</b></label>
		<input type="text" placeholder="Enter Username" name="username" bind:value={username} required />

		<label for="psw"><b>Password</b></label>
		<input type="password" placeholder="Enter Password" name="password" bind:value={password} required />
		<input
			disabled={!(username.trim().length > 0 && password.trim().length > 0)}
			type="submit"
			value="Login"
			onclick={onLogin}
			class="button button-lg button-accent"
		/>
	</div>
</div>

<style>
	.center-content {
		display: flex;
		width: 100%;
		height: 100%;
		align-items: center;
		justify-content: center;
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 16px;
		width: 35vw;
	}

	label {
		font-size: 16px;
	}

	input[type="text"],
	input[type="password"] {
		padding: 9px;
		border-radius: 6px;
		border: 1px solid var(--fg);
	}

	input:focus {
		outline: none;
		border: 2px solid var(--accent-bg) !important;
	}
</style>
