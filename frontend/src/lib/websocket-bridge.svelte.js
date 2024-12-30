import { browser } from "$app/environment";
import { toString } from "cronstrue";
import toast from "svelte-hot-french-toast";
import { SvelteMap } from "svelte/reactivity";
import Cookies from "js-cookie";

export const websocketState = $state({
	authenticated: false,
	/** @type {Map<number, any>} */
	jobs: new SvelteMap(),
});

/** @type {WebSocket} */
let ws;

export function connect() {
	if (browser) {
		ws = new WebSocket("ws://127.0.0.1:1928/");
		//const token = cookies.get("token");
		//C
		//
		ws.addEventListener("open", () => {
			const username = Cookies.get("username");
			const token = Cookies.get("token");
			if (username != undefined && token != undefined) {
				ws.send(
					JSON.stringify({
						action: "LOGIN_TOKEN",
						username: username,
						token: token,
					}),
				);
			}
		});

		ws.addEventListener("message", (message) => {
			if (message.data == "Ack") return;
			const data = JSON.parse(message.data);

			if (data.status == "FAILURE") {
				toast.error(data.message);
			}

			switch (data.action) {
				case "IS_AUTHENTICATED":
					if (data.value) {
						websocketState.authenticated = true;
						Cookies.set("token", data.token);
					}
					break;
				case "UPDATE_JOBS":
					websocketState.jobs.clear();
					data.jobs.forEach((/** @type {{expression: string, id: number, toggled: boolean}} */ job) => {
						websocketState.jobs.set(job.id, {
							expression: job.expression,
							toggled: job.toggled,
							title: toString(job.expression, {
								throwExceptionOnParseError: false,
							}),
						});
					});
					break;
			}
		});
	}
}

/**
 * @param {string} username
 * @param {string} password
 */
export function loginStandard(username, password) {
	Cookies.set("username", username);
	ws.send(
		JSON.stringify({
			action: "LOGIN_STANDARD",
			username: username,
			password: password,
		}),
	);
}

/**
 * @param {string} expression
 * @returns {string}
 */
export function convertExpression(expression) {
	let [seconds, minute, hour, dayOfMonth, month, dayOfWeek] = expression.split(" ");
	// Some quick manipulations so it works on the libcron cron parser
	if (parseInt(dayOfMonth) < 1 || dayOfMonth == "*") {
		dayOfWeek = "?";
	} else if (parseInt(dayOfWeek) < 0 || dayOfWeek == "*") {
		dayOfMonth = "?";
	}
	if (seconds.startsWith("0/")) {
		seconds = "0";
	}
	if (minute.startsWith("0/0")) {
		minute = "*";
	}

	return `${seconds} ${minute} ${hour} ${dayOfMonth} ${month} ${dayOfWeek}`;
}

/**
 * @param {string} expression
 */
export function addJob(expression) {
	ws.send(
		JSON.stringify({
			action: "CREATE_JOB",
			job: convertExpression(expression),
		}),
	);
}

/**
 * @param {number} id
 * @param {string} expression
 */
export function updateJob(id, expression) {
	ws.send(
		JSON.stringify({
			action: "UPDATE_JOB",
			id: id,
			job: convertExpression(expression),
		}),
	);
}

/**
 * @param {number} id
 */
export function deleteJob(id) {
	ws.send(
		JSON.stringify({
			action: "DELETE_JOB",
			id: id,
		}),
	);
}

/**
 * @param {number} id
 * @param {boolean} value
 * @returns {boolean}
 */
export function toggleJob(id, value) {
	websocketState.jobs.set(id, {...value, toggled: value});
	ws.send(
		JSON.stringify({
			action: "TOGGLE_JOB",
			id: id,
			value: value,
		}),
	);
}

export function ring() {
	ws.send(
		JSON.stringify({
			action: "RING",
		}),
	);
}
