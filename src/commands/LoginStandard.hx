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
		DB.instance.validCredentials(json.username, json.password).then(_ -> {
			r.authenticated = true;
			return DB.instance.enrollSession(json.username);
		}, e -> {
			r.send(haxe.Json.stringify({
				status: Routes.Status.FAILURE,
				action: Routes.Commands.LOGIN_STANDARD,
				message: "Invalid Credentials",
			}));
			return null;
		}).then(token -> {
			r.token = token;
			r.send(haxe.Json.stringify({
				status: Routes.Status.SUCCESS,
				action: Routes.Commands.IS_AUTHENTICATED,
				value: r.authenticated,
				token: token,
			}));
		});
	}
}
