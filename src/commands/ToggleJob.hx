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
		entities.Job.findById(json.id).then(job -> {
			if (job == null) {
				r.send(haxe.Json.stringify({
					status: Routes.Status.FAILURE,
					action: Routes.Commands.TOGGLE_JOB,
					message: "Job does not exist",
				}));
			}
			job.isToggled = value;
			if (value) {
				Bell.schedule(job.expression, job.jobId);
			} else {
				Bell.unschedule(job.jobId);
			}
			return job.update();
		});
	}
}
