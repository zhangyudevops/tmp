package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
)

type sClean struct{}

func Clean() *sClean {
	return &sClean{}
}

func (s *sClean) HarborImageClean(ctx context.Context) error {
	// delete all tags from harbor except the latest 5
	projects, err := Harbor().ListProjects(ctx)
	if err != nil {
		return err
	}
	for _, project := range projects {
		// if current project is empty, skip
		if project.RepoCount > 0 {
			repos, err := Harbor().ListRepos(ctx, project.Name, "name", project.RepoCount, 1)
			if err != nil {
				return err
			}
			for _, repo := range repos {
				// if current repo artifacts greater than 5, next
				if repo.ArtifactCount > 5 {
					repoName := strings.Split(repo.Name, "/")[1]
					tags, err := Harbor().ListDigest(ctx, project.Name, repoName, "-push_time", repo.ArtifactCount, 1)
					if err != nil {
						return err
					}

					// delete all tags except the latest 5
					for _, tag := range tags[8:] {
						err = Harbor().DeleteTag(ctx, project.Name, repoName, tag.Digest)
						if err != nil {
							return err
						}
						ip, _ := Config().ParseConfig(ctx, "harbor.ip")
						g.Log().Debugf(ctx, "delete tag: %s", ip+"/"+project.Name+"/"+repoName+":"+tag.Tags[0].Name)
					}
				}
			}
		}
	}

	// create a new cron job for harbor gc
	schedule, err := Harbor().GetGcCron(ctx)
	if schedule == nil {
		err = Harbor().CreateGcCron(ctx)
		if err != nil {
			return err
		}
		g.Log().Info(ctx, "create harbor gc cron job success")
	} else {
		g.Log().Debug(ctx, "harbor gc cron job already exists")
	}

	return nil
}
