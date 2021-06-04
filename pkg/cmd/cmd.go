// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cmd

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// IsFileExist is
func IsFileExist(filepath string) bool {
	// if remote file is
	// ls -l | grep aa | wc -l
	fileName := path.Base(filepath) // aa
	fileDirName := path.Dir(filepath)
	fileCommand := fmt.Sprintf("ls -l %s | grep %s | wc -l", fileDirName, fileName)
	data := strings.Replace(CmdToString("/bin/sh", "-c", fileCommand), "\r", "", -1)
	data = strings.Replace(data, "\n", "", -1)
	count, err := strconv.Atoi(strings.TrimSpace(data))
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[os][%s]RemoteFileExist:%s", filepath, err)
		}
	}()
	if err != nil {
		panic(1)
	}
	return count != 0
}

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
func CmdToString(name string, arg ...string) string {
	logrus.Debugf("[os]exec cmd is : ", name, arg)
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	if err != nil {
		logrus.Errorf("[os]os call error: %v", err)
		return ""
	}
	return b.String()
}
