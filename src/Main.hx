package;

import hx.ws.Util;
import entities.User;
import haxe.Timer;
import logging.Logger;
import logging.LogManager;
import hx.ws.Log;
import hx.ws.WebSocketServer;

class Main {
	public static var server: WebSocketServer<Routes>;

	static function main() {
		final port = Std.parseInt(Sys.getEnv("BETTERBELL_PORT")) ?? 1928;
		trace(port);

		// core-haxe logging setup
		LogManager.instance.addAdaptor(new logging.adaptors.ConsoleLogAdaptor());
		final log = new Logger(Main);
		// websocket logging, eventually should be merged with the core-haxe one
		Log.mask = 0;
		DB.instance = new DB();
		Bell.instance = new Bell();

		User.create("l0go", "password");
		entities.Peer.findByAddress("BETTERBELL__SELF").then(peer -> {
			if (peer == null) {
				var peer = new entities.Peer();
				peer.address = "BETTERBELL__SELF";
				peer.accessToken = Util.generateUUID();
				trace(peer.accessToken);
				peer.add();
			}
		});

		// Reschedule all jobs
		entities.Job.findAll().then(result -> {
			for (job in result) {
				if (job.isToggled) {
					Bell.schedule(job.expression, job.jobId);
				}
			}
		});

		// Setup websocket server
        server = new WebSocketServer<Routes>("localhost", port, 10);
		var timer = new Timer(5 * 1000);
		timer.run = () -> {
			server.sendAll("Ack");
		};

		log.info("Websocket server up");
		server.start();
	}
}
