package commands;

class CreatePeer implements Command {
	public var requiresAuthentication = true;

	public function new() {}
	public function run(r: Routes, json: Dynamic) {
		if (json.address == null || json.accessToken == null || json.address.length <= 0 || json.accessToken.length <= 0) {
			r.send(haxe.Json.stringify({
				status: Routes.Status.FAILURE,
				action: Routes.Commands.CREATE_PEER,
				message: "Address or Access Token is null",
			}));
			return;
		}

		final peer = new entities.Peer();
		peer.address = Std.string(json.address);
		peer.accessToken = Std.string(json.accessToken);
		peer.add();
		entities.Peer.broadcastUpdate();
	}
}
