package commands;

class Ring implements Command {
	public var requiresAuthentication = true;
	
	public function new() {}
	public function run(r: Routes, json: Dynamic) {
		Bell.instance.ring();
		r.send(haxe.Json.stringify({
			status: Routes.Status.SUCCESS,
			action: Routes.Commands.RING,
		}));
	}
}
