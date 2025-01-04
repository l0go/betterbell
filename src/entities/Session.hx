package entities;

import entities.IEntity;
import promises.Promise;

class Session implements IEntity {
	public var user: User;
	@:size(36) public var accessToken: String;

	public static function findByCredentials(user: entities.User, accessToken: String) {
		return find(Query.query($user == user && $accessToken == accessToken));
	}

	public static function enroll(user: entities.User): Promise<Session> {
		final session = new Session();
		session.user = user;
		session.accessToken = hx.ws.Util.generateUUID();
		return session.add();
	}
}
