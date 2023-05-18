package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
)

type sCron struct{}

func Cron() *sCron {
	return &sCron{}
}

var cron *gcron.Cron

func NewCron() {
	cron = gcron.New()
}

// cleanHarborImagesCronJob clean harbor images cron task
func (s *sCron) cleanHarborImagesCronJob(ctx context.Context) (err error) {
	// set cron task
	if _, err = cron.Add(ctx, "@every 3m", func(ctx context.Context) {
		if err = Clean().HarborImageClean(ctx); err != nil {
			g.Log().Error(ctx, err)
			return
		}
		g.Log().Info(ctx, "clean harbor images cron job success")
	}, "CleanHarborImagesCronJob"); err != nil {
		g.Log().Error(ctx, err)
		return
	}

	return
}

func CronSetUp() {
	NewCron()
	_ = Cron().cleanHarborImagesCronJob(context.Background())
	cron.Start("CleanHarborImagesCronJob")
}
