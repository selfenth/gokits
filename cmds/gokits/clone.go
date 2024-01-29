package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xdeeper/gokits/cmder"
)

type Clone struct {
	Description string

	Project cmder.Arg[string] `usage:"要clone的项目，应该是ssh://或者http://、https://打头，.git结尾"`
}

func (a Clone) Run(ctx context.Context) error {
	if !strings.HasSuffix(a.Project.Value(), ".git") || (!strings.HasPrefix(a.Project.Value(), "https://") &&
		!strings.HasPrefix(a.Project.Value(), "http://") && !strings.HasPrefix(a.Project.Value(), "git@")) {
		return fmt.Errorf("the project url(%s) is not valid", a.Project.Value())
	}
	projRelPath := projectUrlToPath(a.Project.Value())
	projDir := filepath.Join(goPath, "src", projRelPath)
	parentDir := filepath.Dir(projDir)
	if d, err := os.Stat(projDir); err == nil && d.IsDir() {
		return fmt.Errorf("project dir already exists:%s", projDir)
	}
	if err := os.MkdirAll(parentDir, 0777); err != nil {
		return fmt.Errorf("failed create parent dir(%s) of project with error:%v", parentDir, err)
	}

	cloneCmd := exec.Command("git", "clone", a.Project.Value())
	cloneCmd.Dir = parentDir
	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("failed clone project:%s into dir:%s with error:%v", a.Project.Value(), projDir, err)
	}
	return nil
}

func projectUrlToPath(projUrl string) string {
	return strings.ReplaceAll(
		strings.TrimPrefix(
			strings.TrimPrefix(
				strings.TrimPrefix(
					strings.TrimSuffix(projUrl, ".git"),
					"http://",
				),
				"https://",
			),
			"git@",
		),
		":", "/",
	)
}
