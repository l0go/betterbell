package commands;

class ToggleJob implements Command {
	public var requiresAuthentication = true;

	public function new() {}
	public function run(r: Routes, json: Dynamic) {
		if (json.id == null || json.value == null) {
			r.send(haxe.Json.stringify({
				status: Routes.Status.FAILURE,
				action: Routes.Commands.TOGGLE_JOB,
				message: "ID or value is null",
			}));
			return;
		}
		final value = (json.value : String).toLowerCase() == "true";
		DB.instance.setJobToggled(json.id, value);
	}
}
