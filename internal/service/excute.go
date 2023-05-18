package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"os/exec"
	"strings"
)

type sShell struct{}

func Shell() *sShell {
	return &sShell{}
}

// Exec execute shell script
// args: file script file absolute path
// args: args is shell args
func (s *sShell) Exec(ctx context.Context, file string, args []string) (error, []byte) {
	var fullArgs string
	index := len(args)
	if index > 0 {
		fullArgs = strings.Join(args, " ")
	}
	file = file + " " + fullArgs
	bytes, err := exec.Command("/bin/bash", "-c", file).Output()
	if err != nil {
		g.Log().Error(ctx, err)
		return err, nil
	}
	return nil, bytes
}
