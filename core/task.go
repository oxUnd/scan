package core

import (
	"strconv"
)

func NewTask(id int, ip string, ports []int, useSSL bool) Task {
	return Task{
		ip:     ip,
		id:     id,
		ports:  ports,
		useSSL: useSSL,
	}
}

type Task struct {
	ip     string
	id     int
	ports  []int
	useSSL bool
}

func (t *Task) GetPorts() []int {
	return t.ports
}

func (t *Task) AddPort(port int) {
	t.ports = append(t.ports, port)
}

func (t *Task) GetProtocol() string {
	var protocol = "http"
	if t.useSSL {
		protocol = "https"
	}
	return protocol
}

func (t *Task) GetUrl(port int) string {
	return t.GetProtocol() + "://" + t.ip + ":" + strconv.Itoa(port)
}

func (t *Task) GetIp() string {
	return t.ip
}
