package entities;

import entities.IEntity;
import promises.Promise;

class Peer implements IEntity {
	@:size(45) public var address: String;
	@:size(36) public var accessToken: String;

	public static function findByAddress(address: String) {
		return find(Query.query($address == address));
	}

	public static function findByCredentials(address: String, accessToken: String) {
		return find(Query.query($address == address && $accessToken == accessToken));
	}
}
