package;

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
			return result.database.createTable("jobs", [
				{name: "ID", type: Number, options: [PrimaryKey, NotNull, AutoIncrement]},
				{name: "CronJob", type: Text(50), options: [NotNull]},
			]);
		}).then(result -> {
			return result.database.createTable("users", [
				{name: "ID", type: Text(254), options: [PrimaryKey, NotNull]},
				{name: "Hash", type: Text(254), options: [NotNull]},
				{name: "Salt", type: Text(24), options: [NotNull]},
			]);
		});
	}
	
	public function all(table: String): promises.Promise<RecordSet> {
		return new promises.Promise((resolve, reject) -> {
			db.table("jobs").then(result -> {
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
			db.table("users").then(result -> {
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

	public function addJob(cron: String): promises.Promise<Null<String>> {
		return new promises.Promise((resolve, reject) -> {
			db.table("jobs").then(result -> {
				final record = new Record();
				record.field("CronJob", cron);
				return result.table.add(record);
			}).then(_ -> {
				return instance.all("jobs");
			}).then(result -> {
				var jbs = [for (job in result) job.field("CronJob")];
				Main.server.sendAll(haxe.Json.stringify({
					status: Routes.Status.SUCCESS,
					action: Routes.Commands.UPDATE_JOBS,
					jobs: jbs,
				}));
				Bell.schedule('$cron');
				resolve(null);
			}, e -> {
				reject(e.message);
			});
		});
	}

	public function validCredentials(id: String, password: String): promises.Promise<Bool> {
		return new promises.Promise((resolve, reject) -> {
			db.table("users").then(result -> {
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
