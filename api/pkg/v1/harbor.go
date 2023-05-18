package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"pack/internal/model"
)

type ListReposReq struct {
	g.Meta      `path:"/harbor/list-repositories" method:"get"  tag:"harbor" summary:"get harbor repositories"`
	ProjectName string `json:"projectName"`
	Sort        string `json:"sort"`
	PageSize    int64  `json:"pageSize"`
	Page        int64  `json:"page"`
}

type ListReposRes struct {
	Repos []model.RepoInfo `json:"repos"`
}

type ListArtifactsReq struct {
	g.Meta      `path:"/harbor/list-artifacts" method:"get"  tag:"harbor" summary:"get artifact list"`
	ProjectName string `json:"projectName"`
	RepoName    string `json:"repoName"`
	Sort        string `json:"sort"`
	PageSize    int64  `json:"pageSize"`
	Page        int64  `json:"page"`
}

type ListArtifactsRes struct {
	Artifacts []model.Artifact `json:"artifacts"`
}

type DeleteArtifactReq struct {
	g.Meta      `path:"/harbor/delete-artifact" method:"delete"  tag:"harbor" summary:"delete artifact"`
	ProjectName string `json:"projectName"`
	RepoName    string `json:"repoName"`
	Digest      string `json:"digest"` // the digest of the image
}

type DeleteArtifactRes struct {
}

type ListProjectReq struct {
	g.Meta `path:"/harbor/list-projects" method:"get"  tag:"harbor" summary:"get project list"`
}

type ListProjectRes struct {
	Projects []model.ProjectInfo `json:"projects"`
}
