package commands;

class DeleteJob implements Command {
	public var requiresAuthentication = true;

	public function new() {}
	public function run(r: Routes, json: Dynamic) {
		if (json.id == null) {
			r.send(haxe.Json.stringify({
				status: Routes.Status.FAILURE,
				action: Routes.Commands.DELETE_JOB,
				message: "ID is null",
			}));
			return;
		}
		DB.instance.deleteJob(json.id);
	}
}
