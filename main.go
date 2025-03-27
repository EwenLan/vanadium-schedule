package main

import (
	"github.com/EwenLan/vanadium-schedule/slog"
)

func main() {
	slog.SetupGlobalLogger()
	slog.SetDisableStandardLogOutput(false)
	hostTopo := HostTopo{}
	hostTopo.Load("hosttopo.json")
	slog.Debugf("Read hosttopo.json: %+v", hostTopo)
}
