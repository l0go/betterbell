package commands;

class LoginStandard implements Command {
	public var requiresAuthentication = false;

	public function new() {}
	public function run(r: Routes, json: Dynamic) {
		// Check for needed fields
		if (json.username == null || json.password == null) {
			r.send(haxe.Json.stringify({
				status: Routes.Status.FAILURE,
				action: Routes.Commands.LOGIN_STANDARD,
				message: "Username or password is null",
			}));
			return;
		}
		entities.User.findByCredentials(json.username, json.password).then(u -> {
			if (u == null) {
				r.send(haxe.Json.stringify({
					status: Routes.Status.FAILURE,
					action: Routes.Commands.LOGIN_STANDARD,
					message: "Invalid Credentials",
				}));
				throw null;
			}
			return entities.Session.enroll(u);
		}).then(session -> {
			r.authenticated = true;
			r.send(haxe.Json.stringify({
				status: Routes.Status.SUCCESS,
				action: Routes.Commands.IS_AUTHENTICATED,
				value: r.authenticated,
				token: session.accessToken,
			}));
			entities.Peer.broadcastUpdate(r);
		});
	}
}
