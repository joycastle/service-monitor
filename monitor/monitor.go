package monitor

import "time"

type Job interface {
	Frequency() int64
	Monitor()
}

var jobs []Job

func AddJob(job Job) {
	jobs = append(jobs, job)
}

func Start() {
	ch := make(chan Job, 100)
	m := make(map[Job]int64)
	for _, job := range jobs {
		m[job] = job.Frequency() + time.Now().Unix()
	}

	//crontab
	go func() {
		for {
			for job, t := range m {
				if time.Now().Unix() >= t {
					ch <- job
					m[job] = job.Frequency() + time.Now().Unix()
				}
			}
			time.Sleep(time.Second * 2)
		}
	}()

	//process
	for i := 0; i < 5; i++ {
		go func() {
			for {
				job := <-ch
				job.Monitor()
			}
		}()
	}
}
