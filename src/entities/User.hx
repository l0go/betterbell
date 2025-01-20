package entities;

import haxe.io.Bytes;
import haxe.crypto.Base64;
import haxe.crypto.BCrypt;
import haxe.crypto.SCrypt;
import promises.Promise;
import entities.IEntity;

@:exposeId
class User implements IEntity {
	@:size(254) public var username: String;
	@:size(254) public var hash: String;
	@:size(24) public var salt: String;

	public static function findByUsername(username: String): Promise<User> {
		return find(Query.query($username == username));
	}

	public static function findByCredentials(username: String, password: String): Promise<User> {
		return findByUsername(username).then(r -> {
			if (r == null) {
				return null;
			}
			var hashword = hashPassword(password, r.salt);
			return find(Query.query($username == username && $hash == hashword.hash));
		});
	}

	public static function create(username: String, password: String): Promise<User> {
		final user = new User();
		final hash = User.hashPassword(password);
		user.username = username;
		user.hash = hash.hash;
		user.salt = hash.salt;
		return user.add();
	}

	public static function hashPassword(password: String, ?salt: String): {hash: String, salt: String} {
		final scrypt = new SCrypt();
		final salt = salt ?? BCrypt.generateSalt();
		final hash = Base64.encode(scrypt.hash(Bytes.ofString(password), Bytes.ofString(salt), 16384, 8, 1, 64));
		return {
			hash: hash,
			salt: salt,
		};
	}
}
