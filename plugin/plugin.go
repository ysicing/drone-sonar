// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ysicing/drone-sonar/pkg/build"
	"github.com/ysicing/drone-sonar/pkg/cmd"
	"github.com/ysicing/sonarapi"
)

type (
	Config struct {
		Key             string
		Host            string
		Token           string
		User            string
		Pass            string
		Branch          string
		PV              string
		Sources         string
		Timeout         string
		Inclusions      string
		Exclusions      string
		Level           string
		UsingProperties bool
		Debug           bool
		Lang            string
		ExtSonarArgs    []string
	}
	Plugin struct {
		Config Config
	}
)

func (p *Plugin) getProjectKey() string {
	return strings.Replace(p.Config.Key, "/", ":", -1)
}

func (p *Plugin) Exec() error {
	if err := Check(p); err != nil {
		logrus.Errorf("check err: %v", err)
		return err
	}
	args := []string{
		"-Dsonar.host.url=" + p.Config.Host,
		"-Dsonar.login=" + p.Config.Token,
	}

	p.preCompile()
	if p.Config.Lang != "no" {
		args = append(args, fmt.Sprintf("-Dsonar.language=%v", p.Config.Lang))
	}

	gitbranch := p.gitbranch()
	gitsha := p.gitsha()

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
	}

	args = append(args, analysis()...)

	args = append(args, "-Dsonar.branch.name="+gitbranch)
	args = append(args, fmt.Sprintf("-Dsonar.projectVersion=%v-%v-%v", gitbranch, getToday(), gitsha))

	if p.Config.Debug {
		debugargs := []string{
			"-Dsonar.showProfiling=true",
			"-X",
		}
		args = append(args, debugargs...)
	}

	if len(p.Config.ExtSonarArgs) > 0 {
		for _, arg := range p.Config.ExtSonarArgs {
			if len(strings.Split(arg, "=")) == 2 {
				if strings.HasSuffix(arg, "-D") {
					args = append(args, arg)
				} else {
					args = append(args, fmt.Sprintf("-D%v", arg))
				}
			} else {
				logrus.Warnf("ext args: %v err, skip", arg)
			}
		}
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
	p.RevokeToken()
	return nil
}

func Check(p *Plugin) error {
	logrus.Debugf("==> Check sonar status: ")
	api, err := Api(p.Config.Host, p.Config.User, p.Config.Pass, p.getProjectKey())
	if err != nil {
		return err
	}
	checkres := api.Health()
	if checkres == nil {
		return fmt.Errorf("sonarapi err, not return anything")
	}
	logrus.Debugf("Sonar Version: %v, Status: %v",
		checkres.(*sonarapi.SystemStatusObject).Version,
		checkres.(*sonarapi.SystemStatusObject).Status)

	exist, err := api.CheckProject()
	if err != nil {
		return err
	}
	if exist {
		logrus.Debugf("==> Check key exist ")
	} else {
		logrus.Debugf("==> Check key not exist, will create project")
		err := api.CreateProject()
		if err != nil {
			return err
		}
		logrus.Debugf("==> Create project: %v", p.Config.Key)
	}
	if len(p.Config.Token) == 0 {
		logrus.Debugf("==> Generate temporary token: %v", "ci-token-***********")
		token, err := api.GenerateToken()
		if err != nil {
			return err
		}
		p.Config.Token = token
	}
	return nil
}

func (p *Plugin) gitsha() string {
	if len(p.Config.PV) != 0 {
		return p.Config.PV
	}
	gitres := cmd.CmdToString("git rev-parse --short HEAD")
	gitres = gitutil(gitres, "sha")
	logrus.Debugf("==> Check Git Commit Sha: %v", gitres)
	return gitres
}

func (p *Plugin) gitbranch() string {
	gitres := cmd.CmdToString("git symbolic-ref --short -q HEAD")
	gitres = gitutil(gitres, "branch")
	if len(p.Config.Branch) != 0 && p.Config.Branch != gitres {
		logrus.Warnf("==> Detect Git Branch: %v, Use Branch: %v", gitres, p.Config.Branch)
		return p.Config.Branch
	}
	logrus.Debugf("==> Check Git Branch: %v", gitres)
	return gitres
}

func (p *Plugin) preCompile() {
	p.Config.Lang = build.GetLangType(p.Config.Sources).String()
	logrus.Debugf("==> Pre Compile Detect LANGUAGE: %v", p.Config.Lang)
	switch p.Config.Lang {
	case "go":
		p.Config.Exclusions = "*.conf,*.yaml,*.ini,*.properties,*.json,*.xml,*.toml,**/*_test.go,**/vendor/**"
		// p.Config.ExtSonarArgs = append(p.Config.ExtSonarArgs, "-Dsonar.go.golangci-lint.reportPaths=report.xml")
		// p.runlint("golangci-lint run --out-format checkstyle ./... > report.xml")
	}
}

// func (p *Plugin) runlint(lintcmd string) {
// 	logrus.Debugf("==> Lint Code")
// 	cmd.CmdToStdout(lintcmd)
// 	logrus.Debugf("==> Lint Code Done")
// }

func gitutil(g, mode string) string {
	if len(g) == 0 {
		switch mode {
		case "branch":
			g = getEnv("CI_COMMIT_REF_NAME", "unknow")
		case "sha":
			g = getEnv("CI_COMMIT_SHA", "unknow")
		default:
			g = "unknow"
		}
	} else {
		reg := regexp.MustCompile(`\s+`)
		g = reg.ReplaceAllString(g, "")
	}
	return g
}

func (p *Plugin) RevokeToken() {
	api, err := Api(p.Config.Host, p.Config.User, p.Config.Pass, p.getProjectKey())
	if err != nil {
		logrus.Errorf("==> Revoke temporary token err")
		return
	}
	if err := api.RevokeToken(); err != nil {
		logrus.Errorf("==> Revoke temporary token err")
		return
	}
	logrus.Debugf("==> Revoke temporary token ci-token-*********** done.")
}

func getToday() string {
	return time.Now().Format("20060102")
}

func analysis() []string {
	return []string{"-Dsonar.sourceEncoding=UTF-8", "-Dsonar.analysis.runtype=gaeaapi"}
}

func getEnv(envstr string, fallback ...string) string {
	e := os.Getenv(envstr)
	if e == "" && len(fallback) > 0 {
		e = fallback[0]
	}
	logrus.Debugf("==> Try Detect Env: %v, value: %v", envstr, e)
	return e
}
