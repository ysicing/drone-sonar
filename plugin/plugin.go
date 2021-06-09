// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"github.com/sirupsen/logrus"
	"github.com/ysicing/drone-sonar/pkg/cmd"
	"os"
	"os/exec"
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
		Version         string
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

	args = append(args, "-Dsonar.projectVersion=" + p.git())

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

func (p Plugin) git() string {
	if len(p.Config.Version) != 0 {
		return p.Config.Version
	}
	gitres := cmd.CmdToString("git rev-parse --short HEAD")
	if len(gitres) != 0 {
		return gitres
	}
	return "unknow"
}
