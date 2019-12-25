package main

import (
	"fmt"
	"os"
	"strings"
	"tank/cgroup"
	"tank/cgroup/subsystem"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"tank/container"
)

const usage = "Tank is a simple container runtime implement"

func main() {

	app := cli.NewApp()
	app.Name = "Tank"
	app.Usage = usage

	app.Commands = []cli.Command{initCommand, runCommand}

	app.Before = func(context *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user process in container",
	Action: func(context *cli.Context) error {
		var err error
		logrus.Infof("init come on")
		cmd := context.Args().Get(0)
		container.RunContainerInitProcess(cmd, nil)
		return err
	},
}

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container with namespace and cgroups limit tank run -ti [command]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory",
		},
	},
	/**
	这里是 run 命令执行的真正函数
	1、判断参数是否包含 command
	2、获取用户指定的 command
	3、调用 Run function 去准备启动容器
	*/
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("ti")
		resConf := &subsystem.ResourceConfig{
			MemoryLimit: context.String("m"),
		}
		Run(tty, cmdArray, resConf)
		return nil
	},
}

func Run(tty bool, cmdArray []string, resConf *subsystem.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Errorf("new parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}
	cgroupManager := cgroup.NewCgroupManager("tank-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(resConf)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(cmdArray, writePipe)

	parent.Wait()
	//os.Exit(-1)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
