package session

var Commands map[string]Command

func InitializeCommand() {
	Commands = make(map[string]Command)
	Commands["status"] = Status{}
	Commands["gc"] = Gc{}
	Commands["velocity"] = VelocityCommand{}
	Commands["attack"] = AttackCommand{}
}
