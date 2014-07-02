package core

import (
	"log"
	"math/rand"
)

const (
	HTTP = "http"
	FTP  = "ftp"
)

func NewWorker(id int, task *Task, reports *[]Report, timeout string) Worker {
	return Worker{
		id:      id,
		task:    task,
		reports: reports,
		timeout: timeout,
	}
}

type Worker struct {
	id       int
	task     *Task
	reports  *[]Report
	timeout  string
	protocol string
}

func (w Worker) ScanFtp() (int, error) {
	return 0, nil
}

func (w Worker) ScanHttp() (int, error) {
	ports := w.task.GetPorts()

	if len(ports) == 0 {
		ports = append(ports, 80)
	}
	s := make(chan (map[string]error))
	for _, port := range ports {
		go func(port int) {
			url_ := w.task.GetUrl(port)
			ret := map[string]error{
				url_: nil,
			}
			log.Println(url_)
			res, err := HttpHeadRequest(w.task.GetUrl(port), w.timeout)
			r := NewReport(rand.Int(), w.task.GetIp(), w.task.GetProtocol(), port)
			if err != nil {
				ret[url_] = err
				r.CatchError(err)
			} else {
				r.Add("Header", res.Header)
				r.Add("Server", res.Header.Get("Server"))
			}
			w.addReport(r)
			s <- ret
		}(port)
	}

	for i := 0; i < len(ports); i++ {
		r := <-s
		for k, v := range r {
			log.Println(k, v)
		}
	}

	return 0, nil
}

func (w *Worker) addReport(r Report) {
	*w.reports = append(*w.reports, r)
}
