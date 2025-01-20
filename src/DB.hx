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
}
