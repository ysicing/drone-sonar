// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/ysicing/drone-sonar/plugin"
	"os"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	app := cli.NewApp()
	app.Name = "drone sonar"
	app.Usage = "Drone Sonar"
	app.Action = run
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "key",
			Aliases: []string{"k"},
			Usage:   "key",
			Value:   "",
			EnvVars: []string{"KEY", "PROJECTKEY", "PLUGIN_KEY", "PLUGIN_PROJECTKEY"},
		},
		&cli.StringFlag{
			Name:    "sources",
			Aliases: []string{"s"},
			Usage:   "sources",
			Value:   ".",
			EnvVars: []string{"SOURCES", "PLUGIN_SOURCES"},
		},
		&cli.StringFlag{
			Name:    "host",
			Usage:   "sonar host",
			Value:   "",
			EnvVars: []string{"HOST", "SONAR_HOST", "PLUGIN_HOST", "PLUGIN_SONAR_HOST"},
		},
		&cli.StringFlag{
			Name:    "login",
			Aliases: []string{"t"},
			Usage:   "sonar login token",
			Value:   "",
			EnvVars: []string{"LOGIN", "TOKEN", "PLUGIN_LOGIN", "PLUGIN_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "user",
			Usage:   "sonar user",
			Value:   "admin",
			EnvVars: []string{"USER", "PLUGIN_USER"},
		},
		&cli.StringFlag{
			Name:    "pass",
			Usage:   "sonar user password",
			Value:   "admin",
			EnvVars: []string{"PASS", "PLUGIN_PASS"},
		},
		&cli.BoolFlag{
			Name:    "usingProperties",
			Usage:   "use Properties",
			EnvVars: []string{"USINGPROPERTIES", "PLUGIN_USINGPROPERTIES"},
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "debug",
			EnvVars: []string{"DEBUG", "PLUGIN_DEBUG"},
		},
		&cli.StringFlag{
			Name:    "level",
			Usage:   "log level",
			Value:   "INFO",
			EnvVars: []string{"PLUGIN_LEVEL", "LEVEL"},
		},
		//&cli.StringFlag{
		//	Name:   "branch",
		//	Aliases: []string{"b"},
		//	Usage:  "Project branch",
		//	EnvVars: []string{"DRONE_BRANCH", "PLUGIN_BRANCH", "BRANCH"},
		//},
		&cli.StringFlag{
			Name:    "timeout",
			Usage:   "Web request timeout",
			Value:   "60",
			EnvVars: []string{"PLUGIN_TIMEOUT", "TIMEOUT"},
		},
		&cli.StringFlag{
			Name:    "inclusions",
			Aliases: []string{"ins"},
			Usage:   "code inclusions",
			EnvVars: []string{"PLUGIN_INCLUSIONS", "INCLUSIONS", "INS"},
		},
		&cli.StringFlag{
			Name:    "exclusions",
			Aliases: []string{"exs"},
			Usage:   "code exclusions",
			EnvVars: []string{"PLUGIN_EXCLUSIONS", "EXCLUSIONS", "EXS"},
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) error {
	p := plugin.Plugin{Config: plugin.Config{
		Key:   c.String("key"),
		Host:  c.String("host"),
		Token: c.String("login"),
		User:  c.String("user"),
		Pass:  c.String("pass"),
		// Branch:          c.String("branch"),
		Sources:    c.String("sources"),
		Timeout:    c.String("timeout"),
		Inclusions: c.String("inclusions"),
		Exclusions: c.String("exclusions"),
		Level:      c.String("level"),
		// BranchAnalysis:  true,
		UsingProperties: c.Bool("usingProperties"),
		Debug:           c.Bool("debug"),
	}}
	return p.Exec()
}
