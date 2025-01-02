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

		final job = new entities.Job();
		job.expression = json.job;
		job.isToggled = true;
		job.add();
		DB.instance.broadcastUpdateJobs();
	}
}
