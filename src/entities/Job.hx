package entities;

import entities.IEntity;

@:exposeId
class Job implements IEntity {
	@:size(50) public var expression: String;
	public var isToggled: Bool;

	public static function broadcastUpdate(?router: Null<Routes>) {
		entities.Job.findAll().then(result -> {
			final jbs = [for (job in result) {
				id: job.jobId,
				expression: job.expression,
				toggled: job.isToggled,
			}];

			final resp = haxe.Json.stringify({
				status: Routes.Status.SUCCESS,
				action: Routes.Commands.UPDATE_JOBS,
				jobs: jbs,
			});

			if (router == null) {
				Main.server.sendAll(resp);
			} else {
				router.send(resp);
			}
		});
	}
}
