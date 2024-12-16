package commands;

class UpdateJob implements Command {
	public var requiresAuthentication = true;

	public function new() {}
	public function run(r: Routes, json: Dynamic) {
		if (json.job == null) {
			r.send(haxe.Json.stringify({
				status: Routes.Status.FAILURE,
				action: Routes.Commands.UPDATE_JOB,
				message: "Job is null",
			}));
			return;
		}

		DB.instance.addJob(json.job, true, json.id).then(null, e -> {
			r.send(haxe.Json.stringify({
				status: Routes.Status.FAILURE,
				action: Routes.Commands.UPDATE_JOB,
				message: e,
			}));
		});
	}
}
