import {browser} from '$app/environment';
import { toString } from "cronstrue";

export const websocketState = $state({
	/** @type {Array<any>} */
	jobs: [],
});

/** @type {WebSocket} */
let ws;

export function connect() {
	websocketState.jobs = [];
	if (browser) {
		ws = new WebSocket("ws://127.0.0.1:1928/");
		ws.addEventListener("message", message => {
			const data = JSON.parse(message.data);
			if (data.action = "UPDATE_JOBS") {
				data.jobs.forEach((/** @type {string} */ job) => {
					websocketState.jobs.push({
						expression: job,
						toggled: true,
						title: toString(job, {throwExceptionOnParseError: false}),
					});
				});
			}
		});
	}

	$effect(() => {
		console.log(websocketState.jobs);
	});
}

/**
 * @param {string} expression
 */
export function addJob(expression) {
	console.log(expression);
	const split = expression.split(" ");
	let [seconds, minute, hour, dayOfMonth, month, dayOfWeek] = split;
	if (parseInt(dayOfMonth) < 1 || dayOfMonth == "*") {
		dayOfWeek = "?";
	} else if (parseInt(dayOfWeek) < 0 || dayOfWeek == "*") {
		dayOfMonth = "?";
	}
	if (seconds.startsWith("0/")) {
		seconds = "0";
	}
	ws.send(JSON.stringify({
		action: "CREATE_JOB",
		job: `${seconds} ${minute} ${hour} ${dayOfMonth} ${month} ${dayOfWeek}`,
	}));
}
