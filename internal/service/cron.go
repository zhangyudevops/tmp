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

// CleanHarborImagesCronJob clean harbor images cron task
func (s *sCron) CleanHarborImagesCronJob(ctx context.Context) (err error) {
	// set cron task
	if _, err = gcron.Add(ctx, "@every 3m", func(ctx context.Context) {
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
