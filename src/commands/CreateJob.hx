package commands;

class CreateJob implements Command {
	public var requiresAuthentication = true;

	public function new() {}
	public function run(r: Routes, json: Dynamic) {
		if (json.job == null) {
			r.send(haxe.Json.stringify({
				status: Routes.Status.FAILURE,
				action: Routes.Commands.CREATE_JOB,
				message: "Job is null",
			}));
			return;
		}

		DB.instance.addJob(json.job).then(null, e -> {
			r.send(haxe.Json.stringify({
				status: Routes.Status.FAILURE,
				action: Routes.Commands.CREATE_JOB,
				message: "Could not add job to database"
			}));
		});
	}
}
