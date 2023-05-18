package harbor

import (
	"context"
	apiv1 "pack/api/pkg/v1"
	"pack/internal/service"
)

type cHarbor struct{}

func Harbor() *cHarbor {
	return &cHarbor{}
}

func (c *cHarbor) ListRepos(ctx context.Context, req *apiv1.ListReposReq) (res *apiv1.ListReposRes, err error) {
	ret, err := service.Harbor().ListRepos(ctx, req.ProjectName, req.Sort, req.PageSize, req.Page)
	if err != nil {
		return
	}

	res = &apiv1.ListReposRes{
		Repos: ret,
	}

	return
}

func (c *cHarbor) ListArtifacts(ctx context.Context, req *apiv1.ListArtifactsReq) (res *apiv1.ListArtifactsRes, err error) {
	ret, err := service.Harbor().ListDigest(ctx, req.ProjectName, req.RepoName, req.Sort, req.PageSize, req.Page)
	if err != nil {
		return
	}

	res = &apiv1.ListArtifactsRes{
		Artifacts: ret,
	}

	return
}

func (c *cHarbor) DeleteArtifact(ctx context.Context, req *apiv1.DeleteArtifactReq) (res *apiv1.DeleteArtifactRes, err error) {
	err = service.Harbor().DeleteTag(ctx, req.ProjectName, req.RepoName, req.Digest)
	if err != nil {
		return
	}

	res = &apiv1.DeleteArtifactRes{}

	return
}

func (c *cHarbor) ListProject(ctx context.Context, req *apiv1.ListProjectReq) (res *apiv1.ListProjectRes, err error) {
	ret, err := service.Harbor().ListProjects(ctx)
	if err != nil {
		return
	}

	res = &apiv1.ListProjectRes{
		Projects: ret,
	}

	return
}
