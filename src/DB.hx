package;

import db.DatabaseFactory;
import entities.EntityManager;

class DB {
	public static var instance: DB;

	public function new() {
		EntityManager.instance.database = DatabaseFactory.instance.createDatabase(DatabaseFactory.SQLITE, {
			filename: "bell.db"
		});	
	}

	public function broadcastUpdateJobs(?router: Null<Routes>) {
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
