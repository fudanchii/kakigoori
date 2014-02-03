package handler

import (
	"fmt"
	"os/exec"

	"github.com/fudanchii/kakigoori/event"
)

type SpawnerConfig struct {
	Cmd  string   `mapstructure:"cmd"`
	Args []string `mapstructure:"args"`
}

func Spawner(intent *event.Intent, config event.Config) {
	var cfg SpawnerConfig
	config.Cast("spawner", &cfg)
	cmd := cfg.Cmd
	args := append([]string{}, cfg.Args...)
	args = append(args, intent.FileName, event.EventName[intent.EventId])
	spawnee := exec.Command(cmd, args...)
	output, err := spawnee.Output()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(output))
	}
}
