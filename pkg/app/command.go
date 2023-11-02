package app

import (
	"context"
	"github.com/robfig/cron"
	"microservices/pkg/log"
	"microservices/pkg/meta"
)

type command struct {
	onceCommands  []meta.CommandFunc
	timerCommands []meta.TimerCommandFunc
}

var c = &command{}

func RegisterOnce(cmd meta.CommandFunc) {
	c.onceCommands = append(c.onceCommands, cmd)
}

func RegisterTimer(cmd meta.CommandFunc, cron string) {
	c.timerCommands = append(c.timerCommands, meta.TimerCommandFunc{
		CommandFunc: cmd,
		Cron:        cron,
	})
}

func (c *command) Run(ctx context.Context) error {
	for _, command := range c.onceCommands {
		go func(command meta.CommandFunc) {
			defer func() {
				if err := recover(); err != nil {
					log.Info(ctx, "command-error", map[string]any{
						"err": err,
					})
				}
			}()
			if err := command(ctx); err != nil {
				panic(err)
			}
		}(command)
	}
	c.runTimerFunc(ctx)
	return nil
}

func (c *command) runTimerFunc(ctx context.Context) {
	cronjob := cron.New()
	for _, timerCommand := range c.timerCommands {
		if err := cronjob.AddFunc(timerCommand.Cron, func() {
			defer func() {
				if err := recover(); err != nil {
					log.Info(ctx, "timer-command-error", map[string]any{
						"err": err,
					})
				}
			}()
			if err := timerCommand.CommandFunc(ctx); err != nil {
				panic(err)
			}
		}); err != nil {
			panic(err)
		}
	}
	cronjob.Start()
}
