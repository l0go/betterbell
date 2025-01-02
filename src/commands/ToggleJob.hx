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
		entities.Job.findById(json.id).then(r -> {
			r.isToggled = value;
			r.update();
		});
	}
}
