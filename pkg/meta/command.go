package meta

import "context"

type CommandFunc func(ctx context.Context) error

type TimerCommandFunc struct {
	CommandFunc
	Cron string
}
