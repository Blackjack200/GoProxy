package session

var Commands map[string]*Command
var i = false

func InitializeCommand() {
	if !i {
		Commands = make(map[string]*Command)
		Register("status", Status{})
		Register("gc", Gc{})
		Register("velocity", VelocityCommand{})
		Register("attack", RateCommand{})
		Register("nofall", NoFallCommand{})
		Register("speed", SpeedCommand{})
		Register("fly", FlyCommand{})
		Register("noclip", NoClipCommand{})
		Register("killaura", KillAuraCommand{})
		Register("gamemode", GameModeCommand{})
		i = true
	}
}

func Register(name string, command Command) {
	Commands[name] = &command
}
