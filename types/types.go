package types

import "fmt"

type Machine struct {
	NickName string
	HostName string
	User     string
	CertFile string
}

func NewMachine(nickname, hostname, user, certfile string) *Machine {
	return &Machine{
		NickName: nickname,
		HostName: hostname,
		User:     user,
		CertFile: certfile,
	}
}

func (m *Machine) GetLoginCmd() string {
	return fmt.Sprintf("ssh -tt -i \"%s\" %s@%s", m.CertFile, m.User, m.HostName)
}
