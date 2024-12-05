package;

import haxe.Json;
import hx.ws.SocketImpl;
import hx.ws.WebSocketHandler;
import hx.ws.Types;
import logging.Logger;

using StringTools;

enum abstract Commands(String) to String {
	var UPDATE_JOBS;
	var CREATE_JOB;
	var IS_AUTHENTICATED;
	var LOGIN_STANDARD;
}

enum abstract Status(String) to String {
	var SUCCESS;
	var FAILURE;
}

class Routes extends WebSocketHandler {
	static final log: Logger = new Logger(Routes);
	var authenticated = false;

	public function new(s: SocketImpl) {
		super(s);

		onopen = () -> {
            trace(id + ". OPEN");
			DB.instance.all("jobs").then(result -> {
				var jbs = [for (job in result) job.field("CronJob")];
				send(Json.stringify({
					status: Status.SUCCESS,
					action: Commands.UPDATE_JOBS,
					jobs: jbs,
				}));
			});
        };

        onclose = () -> {
			authenticated = false;
        };

		onmessage = (message: MessageType) -> try {
			switch (message) {
				case StrMessage(content):
					var json = haxe.Json.parse(content);
					switch (json.action) {
						case Commands.LOGIN_STANDARD:
							// Check for needed fields
							if (json.username == null || json.password == null) {
								send(Json.stringify({
									status: Status.FAILURE,
									action: Commands.LOGIN_STANDARD,
									message: "Username or password is null",
								}));
								return;
							}
							DB.instance.validCredentials(json.username, json.password).then(_ -> {
								authenticated = true;
								send(Json.stringify({
									status: Status.SUCCESS,
									action: Commands.LOGIN_STANDARD,
									message: "Authenticated",
								}));
							}, e -> {
								send(Json.stringify({
									status: Status.FAILURE,
									action: Commands.LOGIN_STANDARD,
									message: "Invalid Credentials",
								}));
							});
						case IS_AUTHENTICATED:
							send(Json.stringify({
								status: Status.SUCCESS,
								action: Commands.IS_AUTHENTICATED,
								message: '$authenticated',
							}));
						case CREATE_JOB:
							//if (!authenticated) {
							//	send(Json.stringify({
							//		status: Status.FAILURE,
							//		action: Commands.CREATE_JOB,
							//		message: "Must be authenticated",
							//	}));
							//	return;
							//}

							if (json.job == null) {
								send(Json.stringify({
									status: Status.FAILURE,
									action: Commands.CREATE_JOB,
									message: "Job is null",
								}));
								return;
							}

							DB.instance.addJob(json.job).then(null, e -> {
								send(Json.stringify({
									status: Status.FAILURE,
									action: Commands.CREATE_JOB,
									message: "Could not add job to database"
								}));
							});
						default:
							send(Json.stringify({
								status: Status.FAILURE,
								message: "Undefined action",
							}));
					}
				default:
					send("Message must be a string");
			}
		} catch(e) {
			log.error(e.message);
			send(Json.stringify({
				status: Status.FAILURE,
				message: "Uncaught Error",
			}));
		};

        onerror = (error) -> {
            trace(id + ". ERROR: " + error);
        };
	}
}
