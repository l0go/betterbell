package commands;

class LoginToken implements Command {
	public var requiresAuthentication = false;

	public function new() {}
	public function run(r: Routes, json: Dynamic) {
		DB.instance.validSession(json.username, json.token).then(result -> {
			if (result) {
				r.authenticated = true;
				r.token = json.token;
				r.send(haxe.Json.stringify({
					status: Routes.Status.SUCCESS,
					action: Routes.Commands.IS_AUTHENTICATED,
					value: r.authenticated,
					token: json.token,
				}));
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
