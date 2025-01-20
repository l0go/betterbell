package commands;

class PeerRing implements Command {
	public var requiresAuthentication = false;
	
	public function new() {}
	public function run(r: Routes, json: Dynamic) {
		entities.Peer.findByCredentials("BETTERBELL__SELF", json.token).then(peer -> {
			if (peer != null) {
				Bell.instance.ring();
				r.send(haxe.Json.stringify({
					status: Routes.Status.SUCCESS,
					action: Routes.Commands.PEER_RING,
				}));
			}
		});
	}
}
