package app

import (
	"context"
	"microservices/pkg/meta"
)

type command struct {
	onceCommands  []meta.CommandFunc
	timerCommands []meta.CommandFunc
}

var c = &command{}

func RegisterOnce(cmd meta.CommandFunc) {
	c.onceCommands = append(c.onceCommands, cmd)
}

func RegisterTimer(cmd meta.CommandFunc) {
	c.timerCommands = append(c.timerCommands, cmd)
}

func (c *command) Run(ctx context.Context) error {
	for _, command := range c.onceCommands {
		if err := command(ctx); err != nil {
			return err
		}
	}
	return nil
}
