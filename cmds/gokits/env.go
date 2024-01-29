package main

import (
	"bytes"
	"os/exec"
	"strings"
)

var (
	goPath string
)

func init() {
	goPath = must1(GetGoEnv("GOPATH"))
}
func GetGoEnv(name string) (val string, err error) {
	cmd := exec.Command("go", "env", name)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err = cmd.Run(); err != nil {
		return
	}
	return strings.TrimSpace(buf.String()), nil
}

func must1[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
