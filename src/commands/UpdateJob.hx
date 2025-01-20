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

		entities.Job.findById(json.id).then(job -> {
			if (job == null) {
				r.send(haxe.Json.stringify({
					status: Routes.Status.FAILURE,
					action: Routes.Commands.UPDATE_JOB,
					message: "Job does not exist",
				}));
				return;
			}
			Bell.unschedule(job.jobId);
			job.expression = json.job;
			job.update();
			Bell.schedule(job.expression, job.jobId);
			entities.Job.broadcastUpdate();
		});
	}
}
