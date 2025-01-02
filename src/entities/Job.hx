package entities;

import entities.IEntity;

@:exposeId
class Job implements IEntity {
	@:size(50) public var expression: String;
	public var isToggled: Bool;
}
