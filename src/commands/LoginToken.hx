package commands;

class LoginToken implements Command {
	public var requiresAuthentication = false;

	public function new() {}
	public function run(r: Routes, json: Dynamic) {
		entities.User.findByUsername(json.username).then(user -> {
			if (user == null) {
				throw null;
			}
			return entities.Session.findByCredentials(user, json.token);
		}).then(result -> {
			if (result != null) {
				r.authenticated = true;
				r.token = json.token;
				r.send(haxe.Json.stringify({
					status: Routes.Status.SUCCESS,
					action: Routes.Commands.IS_AUTHENTICATED,
					value: r.authenticated,
					token: json.token,
				}));
				entities.Peer.broadcastUpdate(r);
			} else {
				r.send(haxe.Json.stringify({
					status: Routes.Status.FAILURE,
					action: Routes.Commands.LOGIN_TOKEN,
					message: "Invalid Credentials",
				}));
			}
		});
	}
}
