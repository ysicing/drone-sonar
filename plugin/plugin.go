// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ysicing/drone-sonar/pkg/cmd"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
)

type (
	Config struct {
		Key   string
		Host  string
		Token string
		User  string
		Pass  string
		// Branch         string
		Sources         string
		Timeout         string
		Inclusions      string
		Exclusions      string
		Level           string
		UsingProperties bool
		Debug           bool
	}
	Plugin struct {
		Config Config
	}
)

func (p Plugin) getProjectKey() string {
	return strings.Replace(p.Config.Key, "/", ":", -1)
}

func (p Plugin) Exec() error {
	if err := p.Check(); err != nil {
		return err
	}
	args := []string{
		"-Dsonar.host.url=" + p.Config.Host,
		"-Dsonar.login=" + p.Config.Token,
		"-Dsonar.sourceEncoding=UTF-8",
	}

	if !p.Config.UsingProperties {
		argsParameter := []string{
			"-Dsonar.projectKey=" + p.getProjectKey(),
			"-Dsonar.sources=" + p.Config.Sources,
			"-Dsonar.ws.timeout=" + p.Config.Timeout,
			"-Dsonar.inclusions=" + p.Config.Inclusions,
			"-Dsonar.exclusions=" + p.Config.Exclusions,
			"-Dsonar.log.level=" + p.Config.Level,
			"-Dsonar.scm.provider=git",
		}
		args = append(args, argsParameter...)
	} else {
		args = append(args, "")
	}

	//if len(p.Config.Branch) != 0 {
	//	args = append(args, "-Dsonar.branch.name=" + p.Config.Branch)
	//}
	if p.Config.Debug {
		debugargs := []string{
			"-Dsonar.showProfiling=true",
			"-X",
		}
		args = append(args, debugargs...)
	}
	cmd := exec.Command("sonar-scanner", args...)
	logrus.Debugf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logrus.Debugf("==> Code Analysis Result:\n")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (p Plugin) Check() error {
	logrus.Debugf("==> Check sonar scan: ")
	return nil
}

func downloadFile(location string) (filePATH string) {
	if _, isUrl := isUrl(location); isUrl {
		absPATH := "/tmp/drone-sonar/" + path.Base(location)
		if !cmd.IsFileExist(absPATH) {
			//generator download cmd
			dwnCmd := downloadCmd(location)
			//os exec download command
			cmd.Cmd("/bin/sh", "-c", "mkdir -p /tmp/drone-sonar && cd /tmp/drone-sonar && "+dwnCmd)
		}
		location = absPATH
	}
	return location
}

func downloadCmd(url string) string {
	//only http
	u, isHttp := isUrl(url)
	var c = ""
	if isHttp {
		param := ""
		if u.Scheme == "https" {
			param = "--no-check-certificate"
		}
		c = fmt.Sprintf(" wget -c %s %s", param, url)
	}
	return c
}

func isUrl(u string) (url.URL, bool) {
	if uu, err := url.Parse(u); err == nil && uu != nil && uu.Host != "" {
		return *uu, true
	}
	return url.URL{}, false
}