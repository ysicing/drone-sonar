// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

//Cmd is exec on os ,no return
func Cmd(name string, arg ...string) {
	logrus.Debugf("[os]exec cmd is : %v %v", name, arg)
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		logrus.Errorf("[os]os call error: %v", err)
	}
}

//CmdToString is exec on os , return result
func CmdToString(args string) string {
	logrus.Debugf("[os]exec cmd is :%v", args)
	cmd := exec.Command("/bin/sh", "-c", args)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		logrus.Errorf("[os]os call error: %v", err)
		return ""
	}
	return b.String()
}
