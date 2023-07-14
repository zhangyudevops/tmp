package service

type sCron struct{}

func Cron() *sCron {
	return &sCron{}
}

//var cron *gcron.Cron
//
//func NewCron() {
//	cron = gcron.New()
//}
//
//// cleanHarborImagesCronJob clean harbor images cron task
//func (s *sCron) cleanHarborImagesCronJob(ctx context.Context, pattern string) (err error) {
//	// set cron task
//	if _, err = cron.Add(ctx, pattern, func(ctx context.Context) {
//		if err = Clean().HarborImageClean(ctx); err != nil {
//			return
//		}
//	}, "CleanHarborImagesCronJob"); err != nil {
//		return
//	}
//
//	return
//}
//
//// addAllCronJobs add all cron jobs
//func (s *sCron) addAllCronJobs(ctx context.Context) (err error) {
//	// add clean harbor images cron job
//	pattern, _ := g.Config().Get(ctx, "cron.harbor")
//	if err = Cron().cleanHarborImagesCronJob(ctx, pattern.String()); err != nil {
//		return
//	}
//
//	return
//}
//
//func CronSetUp() {
//	NewCron()
//	if err := Cron().addAllCronJobs(context.Background()); err != nil {
//		g.Log().Error(context.Background(), err)
//	}
//	cron.Start()
//}
