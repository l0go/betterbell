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

		entities.Job.findById(json.id).then(j -> {
			if (j == null) {
				r.send(haxe.Json.stringify({
					status: Routes.Status.FAILURE,
					action: Routes.Commands.UPDATE_JOB,
					message: "Job does not exist",
				}));
				return;
			}
			j.expression = json.job;
			j.update();
			DB.instance.broadcastUpdateJobs();
		});
	}
}
