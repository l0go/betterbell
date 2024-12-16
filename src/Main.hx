package;

import haxe.Timer;
import logging.Logger;
import logging.LogManager;
import hx.ws.Log;
import hx.ws.WebSocketServer;

class Main {
	public static var server: WebSocketServer<Routes>;

	static function main() {
		Bell.instance = new Bell();

		// core-haxe logging setup
		LogManager.instance.addAdaptor(new logging.adaptors.ConsoleLogAdaptor());
		final log = new Logger(Main);
		// websocket logging, eventually should be merged with the core-haxe one
		Log.mask = 0;

		DB.instance = new DB();
		DB.instance.addUser("l0go", "password").then(_ -> {
			trace("Account created");
		}, e -> {
			trace(e);
		});

		// Reschedule all jobs
		DB.instance.all("jobs").then(result -> {
			for (job in result) {
				if (job.field("Toggled")) {
					Bell.schedule(job.field("CronJob"), job.field("ID"));
				}
			}
		});

		// Setup websocket server
        server = new WebSocketServer<Routes>("localhost", 1928, 10);
		var timer = new Timer(5 * 1000);
		timer.run = () -> {
			server.sendAll("Ack");
		};

		log.info("Websocket server up");
		server.start();
	}
}
