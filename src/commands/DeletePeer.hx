package commands;

class DeletePeer implements Command {
	public var requiresAuthentication = true;

	public function new() {}
	public function run(r: Routes, json: Dynamic) {
		if (json.id == null) {
			r.send(haxe.Json.stringify({
				status: Routes.Status.FAILURE,
				action: Routes.Commands.DELETE_PEER,
				message: "ID is null",
			}));
			return;
		}
		entities.Peer.findById(json.id).then(peer -> {
			if (peer == null) {
				r.send(haxe.Json.stringify({
					status: Routes.Status.FAILURE,
					action: Routes.Commands.DELETE_PEER,
					message: "Peer does not exist",
				}));
			} else {
				peer.delete();
			}
		});
		entities.Peer.broadcastUpdate();
	}
}
