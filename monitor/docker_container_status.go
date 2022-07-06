package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/service-monitor/agent"
	"github.com/joycastle/service-monitor/types"
)

type DockerContainerStatus struct {
	Machines []*types.Machine
	Name     string
	frency   int64
}

func NewDockerContainerStatus(containersName string, frency int64, ms ...*types.Machine) Job {
	return &DockerContainerStatus{Machines: ms, Name: containersName, frency: frency}
}

func (dcs *DockerContainerStatus) Frequency() int64 {
	return dcs.frency
}

func (dcs *DockerContainerStatus) Monitor() {
	for _, machine := range dcs.Machines {
		isFatal := 0
		ContainerInfo := ""
		for {
			cmd := MergeCmd(machine.GetLoginCmd(), dcs.cmdStatus())
			ret, err := Exec(cmd)
			if err != nil {
				log.Get("error").Fatal("DockerContainerStatus", "node:", machine.NickName, "cmd:", cmd, err)
				continue
			}
			info, err := dcs.checkAlarm(ret[0])
			ContainerInfo = info
			if err != nil {
				log.Get("run").Fatal("DockerContainerStatus", "node:", machine.NickName, "container:", dcs.Name, "info:", info, "Fatal", "times:", isFatal)
				isFatal++
				if isFatal >= 2 {
					break
				}
				time.Sleep(time.Second * 5)
				continue
			}

			log.Get("run").Info("DockerContainerStatus", "node:", machine.NickName, "container:", dcs.Name, "info:", info, "OK")
			isFatal = 0
			break
		}

		if isFatal >= 2 {
			cmd := MergeCmd(machine.GetLoginCmd(), dcs.cmdDockerLogs())
			ret, err := Exec(cmd)
			if err != nil {
				log.Get("run").Info("DockerContainerStatus", "node:", machine.NickName, "container:", dcs.Name, "info:", ContainerInfo, "OK")
				continue
			}

			msg := fmt.Sprintf("docker容器运行异常,请及时查看!\n容器名称: %s\n环境: %s\n时间: %s\n容器状态: %s\n容器日志:\n%v",
				dcs.Name, machine.NickName, time.Now().Format("2006-01-02 15:04:05"), ContainerInfo, ret)
			if err := agent.FeiShuSendToServiceRD(msg); err != nil {
				log.Get("alert").Fatal("DockerContainerStatus", "send alert fatal", err)
			} else {
				log.Get("alert").Fatal("DockerContainerStatus", "send alert success")
			}
		}
	}
}

func (dcs *DockerContainerStatus) cmdStatus() string {
	return fmt.Sprintf("docker ps -a --filter name=%s --format \"{{.Status}}\"", dcs.Name)
}

func (dcs *DockerContainerStatus) cmdDockerLogs() string {
	return fmt.Sprintf("docker logs -n 10 %s", dcs.Name)
}

func (dcs *DockerContainerStatus) checkAlarm(ret string) (string, error) {
	if strings.Contains(ret, "Exit") || strings.Contains(ret, "Restart") {
		return ret, fmt.Errorf("Program exception")
	}
	return ret, nil
}
