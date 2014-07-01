package core

import (
	"github.com/kr/pretty"
)

func NewReport(id int, ip string, protocol string, port int) Report {
	return Report{
		id:       id,
		ip:       ip,
		protocol: protocol,
		port:     port,
		result:   map[string]interface{}{},
		err:      nil,
	}
}

type Report struct {
	id       int
	ip       string
	protocol string
	port     int
	result   map[string]interface{}
	err      error
}

func (r Report) Dump() {
	pretty.Println(r)
}

func (r *Report) Add(key string, val interface{}) {
	r.result[key] = val
}

func (r *Report) catchError(err error) {
	r.err = err
}

func (r Report) GetString() string {
	return ""
}

func (r Report) GetJSON() string {
	return ""
}
