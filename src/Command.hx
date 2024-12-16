package;

interface Command {
	public var requiresAuthentication: Bool;
	public function run(r: Routes, json: Dynamic): Void;
}
