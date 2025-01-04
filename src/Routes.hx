package;

import haxe.Json;
import hx.ws.SocketImpl;
import hx.ws.WebSocketHandler;
import hx.ws.Types;
import logging.Logger;

using StringTools;

enum abstract Commands(String) to String {
	final UPDATE_JOBS;
	final CREATE_JOB;
	final UPDATE_JOB;
	final TOGGLE_JOB;
	final DELETE_JOB;
	final IS_AUTHENTICATED;
	final LOGIN_STANDARD;
	final LOGIN_TOKEN;
	final RING;
	final PEER_RING;
}

enum abstract Status(String) to String {
	final SUCCESS;
	final FAILURE;
}

class Routes extends WebSocketHandler {
	static final log: Logger = new Logger(Routes);
	public var authenticated = false;
	public var token: String;

	final commandRoutes: Map<Commands, Command> = [
		CREATE_JOB => new commands.CreateJob(),
		UPDATE_JOB => new commands.UpdateJob(),
		TOGGLE_JOB => new commands.ToggleJob(),
		DELETE_JOB => new commands.DeleteJob(),
		IS_AUTHENTICATED => new commands.IsAuthenticated(),
		LOGIN_STANDARD => new commands.LoginStandard(),
		LOGIN_TOKEN => new commands.LoginToken(),
		RING => new commands.Ring(),
		PEER_RING => new commands.PeerRing(),
	];

	public function new(s: SocketImpl) {
		super(s);

		onopen = () -> {
			log.info(id + ". OPEN");
			DB.instance.broadcastUpdateJobs(this);
		};

        onclose = () -> {
			authenticated = false;
        };

		onmessage = (message: MessageType) -> try {
			switch (message) {
				case StrMessage(content):
					final json = haxe.Json.parse(content);
					if (!commandRoutes.exists(json.action)) {
						send(Json.stringify({
							status: Status.FAILURE,
							message: "Undefined action",
						}));
						return;
					}
					
					if (commandRoutes[json.action].requiresAuthentication && !authenticated) {
						send(haxe.Json.stringify({
							status: Routes.Status.FAILURE,
							action: commandRoutes[json.action],
							message: "Must be authenticated for this action",
						}));
						return;
					}

					commandRoutes[json.action].run(this, json);
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
            log.error('$id: $error');
        };
	}
}
