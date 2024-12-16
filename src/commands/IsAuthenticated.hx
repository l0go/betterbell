package commands;

class IsAuthenticated implements Command {
	public var requiresAuthentication = false;
	
	public function new() {}
	public function run(r: Routes, json: Dynamic) {
		r.send(haxe.Json.stringify({
			status: Routes.Status.SUCCESS,
			action: Routes.Commands.IS_AUTHENTICATED,
			value: r.authenticated,
			token: if (r.authenticated) r.token else null,
		}));
	}
}
