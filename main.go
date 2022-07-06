package main

import (
	"os"
	"os/signal"
	"path/filepath"

	"github.com/joycastle/casual-server-lib/config"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/service-monitor/monitor"
	"github.com/joycastle/service-monitor/types"
)

var (
	jobs []monitor.Job

	mdev  *types.Machine
	mdev1 *types.Machine
	mpre  *types.Machine
	mprod *types.Machine
)

func init() {
	mdev = types.NewMachine("dev", "ec2-3-96-176-114.ca-central-1.compute.amazonaws.com", "ec2-user", "./cert/mm-v2-dev.pem")
	mdev1 = types.NewMachine("dev1", "ec2-3-98-57-247.ca-central-1.compute.amazonaws.com", "ec2-user", "./cert/chengyuyue.pem")
	mpre = types.NewMachine("pre", "ec2-35-182-21-159.ca-central-1.compute.amazonaws.com", "ec2-user", "./cert/chengyuyue.pem")
	mprod = types.NewMachine("prod", "ec2-3-98-57-247.ca-central-1.compute.amazonaws.com", "ec2-user", "./cert/chengyuyue.pem")
}

func main() {
	configFile := filepath.Join("./conf")
	if err := config.InitConfig(configFile); err != nil {
		panic(err)
	}
	log.InitLogs(config.Logs)

	monitor.AddJob(monitor.NewDockerContainerStatus("matching-story-robot-service", 300, mdev, mdev1, mpre, mprod))

	monitor.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
}
