package;

import promises.Promise;
import db.RecordSet;
import haxe.crypto.SCrypt;
import haxe.crypto.BCrypt;
import haxe.crypto.Base64;
import haxe.io.Bytes;
import db.Record;
import db.IDatabase;
import db.DatabaseFactory;
import Query.*;

class DB {
	public static var instance: DB;
	var db: IDatabase;

	public function new() {
		db = DatabaseFactory.instance.createDatabase(DatabaseFactory.SQLITE, {
			filename: "bell.db"
		});	

		db.connect().then(result -> {
			return result.database.createTable("Jobs", [
				{name: "ID", type: Number, options: [PrimaryKey, NotNull, AutoIncrement]},
				{name: "CronJob", type: Text(50), options: [NotNull]},
				{name: "Toggled", type: Boolean, options: [NotNull]},
			]);
		}).then(result -> {
			return result.database.createTable("Users", [
				{name: "ID", type: Text(254), options: [PrimaryKey, NotNull]},
				{name: "Hash", type: Text(254), options: [NotNull]},
				{name: "Salt", type: Text(24), options: [NotNull]},
			]);
		}).then(result -> {
			return result.database.createTable("Sessions", [
				{name: "Token", type: Text(36), options: [PrimaryKey, NotNull]},
				{name: "User", type: Text(254), options: [NotNull]},
			]);
		}).then(_ -> {
			db.defineTableRelationship("Sessions.User", "Users.ID");
		});
	}
	
	public function all(table: String): promises.Promise<RecordSet> {
		return new promises.Promise((resolve, reject) -> {
			db.table("Jobs").then(result -> {
				return result.table.all();
			}).then(result -> {
				resolve(result.data);
			}, e -> {
				reject(e);
			});
		});
	}

	public function addUser(id: String, password: String): promises.Promise<Null<String>> {
		return new promises.Promise((resolve, reject) -> {
			db.table("Users").then(result -> {
				final h = hashPassword(password);
				final record = new Record();
				record.field("ID", id);
				record.field("Hash", h.hash);
				record.field("Salt", h.salt);
				return result.table.add(record);
			}).then(_ -> {
				resolve(null);
			}, e -> {
				reject(e.message);
			});
		});
	}

	public function addJob(cron: String, ?update: Bool, ?id: Int): promises.Promise<Null<String>> {
		return new promises.Promise((resolve, reject) -> {
			db.table("Jobs").then(result -> {
				final record = new Record();
				if (id != null) {
					record.field("ID", id);
				}
				record.field("CronJob", cron);
				record.field("Toggled", true);
				if (update) {
					Bell.unschedule(id);
					return result.table.update(query($ID = id), record);
				}
				return result.table.add(record);
			}).then(r -> {
				Bell.schedule(cron, r.data.field("ID"));
				return instance.all("Jobs");
			}).then(result -> {
				sendUpdateJobs(result);
				resolve(null);
			}, e -> {
				reject(e.message);
			});
		});
	}

	public function deleteJob(id: Int) {
		return db.table("Jobs").then(result ->{
			Bell.unschedule(id);
			result.table.deleteAll(query($ID = id));
			return instance.all("Jobs");
		}).then(result -> {
			sendUpdateJobs(result);
		});
	}

	public function setJobToggled(id: Int, bool: Bool) {
		return db.table("Jobs").then(result -> {
			return result.table.findOne(query($ID = id));
		}).then(result -> {
			var record = new Record();
			record.field("Toggled", bool);
			Bell.unschedule(id);
			if (bool) {
				Bell.schedule(result.data.field("CronJob"), id);
			}
			return result.table.update(query($ID = id), record);
		});
	}

	function sendUpdateJobs(result: RecordSet) {
		final jbs = [for (job in result) {
			id: job.field("ID"),
			expression: job.field("CronJob"),
			toggled: job.field("Toggled") == 1,
		}];
		Main.server.sendAll(haxe.Json.stringify({
			status: Routes.Status.SUCCESS,
			action: Routes.Commands.UPDATE_JOBS,
			jobs: jbs,
		}));
	}

	public function validCredentials(id: String, password: String): promises.Promise<Bool> {
		return new promises.Promise((resolve, reject) -> {
			db.table("Users").then(result -> {
				return result.table.findOne(query($ID = id));
			}).then(result -> {
				if (result?.data == null) {
					reject(false);
					return;
				}
				if (hashPassword(password, result.data.field("Salt")).hash == result.data.field("Hash")) {
					resolve(true);
				} else {
					reject(false);
				}
			}, _ -> {
				reject(false);
			});
		});
	}

	public function enrollSession(id: String): Promise<String> {
		final uuid = hx.ws.Util.generateUUID();
		return new promises.Promise((resolve, reject) -> {
			db.table("Sessions").then(result -> {
				var record = new Record();
				record.field("Token", uuid);
				record.field("User", id);
				return result.table.add(record);
			}).then(_ -> {
				resolve(uuid);
			});
		});
	}

	public function validSession(id: String, token: String): Promise<Bool> {
		return new promises.Promise((resolve, reject) -> {
			db.table("Sessions").then(result -> {
				return result.table.findOne(query($User = id && $Token = token));
			}).then(result -> {
				resolve(result.data != null);
			});
		});
	}

	function hashPassword(password: String, ?salt: String): {hash: String, salt: String} {
		final scrypt = new SCrypt();
		final salt = salt ?? BCrypt.generateSalt();
		final hash = Base64.encode(scrypt.hash(Bytes.ofString(password), Bytes.ofString(salt), 16384, 8, 1, 64));
		return {
			hash: hash,
			salt: salt,
		};
	}
}
