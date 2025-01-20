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

	public static function broadcastUpdate(?router: Null<Routes>) {
		entities.Peer.findAll().then(result -> {
			final peers = [for (peer in result) {
				id: peer.peerId,
				address: peer.address,
				token: peer.address == "BETTERBELL__SELF" ? peer.accessToken : null,
			}];

			final resp = haxe.Json.stringify({
				status: Routes.Status.SUCCESS,
				action: Routes.Commands.UPDATE_PEERS,
				peers: peers,
			});

			if (router == null) {
				for (handler in Main.server.handlers) {
					if (handler.authenticated) {
						handler.send(resp);
					}
				}
			} else {
				router.send(resp);
			}
		});
	}

}
