package worker

import (
	"fmt"
	"github.com/fujiahui/talnet-challenge-payman/common"
	"github.com/fujiahui/talnet-challenge-payman/logger"
	"github.com/fujiahui/talnet-challenge-payman/worker/util"
	"strings"
	"time"
)

// Actuator Task执行管理器
type Actuator struct {
	currTimestamp  common.TimestampType
	jobs           map[common.JobIDType]*util.Job
	capacity       common.PointType
	executingPoint common.PointType
}

func NewActuator(startTimestamp common.TimestampType, capacity common.PointType) *Actuator {
	return &Actuator{
		currTimestamp:  startTimestamp,
		jobs:           make(map[common.JobIDType]*util.Job),
		capacity:       capacity,
		executingPoint: 0,
	}
}

func (c *Actuator) Capacity() common.PointType {
	return c.capacity
}

func (c *Actuator) ExecutingPoint() common.PointType {
	return c.executingPoint
}

func (c *Actuator) CurrTimestamp() common.TimestampType {
	return c.currTimestamp
}

func (c *Actuator) FreePoint() common.PointType {
	return c.capacity - c.executingPoint
}

// Ticking 滴答滴答向前一步步大胆的滴答
func (c *Actuator) Ticking(tick int) []*util.Job {
	time.Sleep(time.Duration(tick) * time.Millisecond)
	c.currTimestamp++

	ids := make([]common.JobIDType, 0, 16)
	jobs := make([]*util.Job, 0, 16)
	for id, job := range c.jobs {
		c.executingPoint -= 1
		t := job.CurrTask()
		t.Ticking()
		if t.Finished() {
			ids = append(ids, id)
			job.NextTask(c.currTimestamp)
			if !job.Finished() {
				jobs = append(jobs, job)
			}
		}
	}
	for _, id := range ids {
		delete(c.jobs, id)
	}

	return jobs
}

// Execute 执行一个Job
func (c *Actuator) Execute(job *util.Job) {
	if _, ok := c.jobs[job.ID()]; ok {
		// job的tasks存在非顺序执行
		logger.Panicf("job not sequentially executes job tasks")
	}
	c.jobs[job.ID()] = job
	job.SetRunning()

	t := job.CurrTask()
	t.SetRunning()

	c.executingPoint += t.TaskPoint()
	if c.executingPoint > c.capacity {
		// 容量超限
		logger.Panicf("Actuator's executingPoint %d more than capacity %d ", c.executingPoint, c.capacity)
	}
	return
}

func (c *Actuator) String() string {
	jj := make([]string, 0, 16)
	for _, job := range c.jobs {
		jj = append(jj, job.String())
	}

	tt := make([]string, 3, 4)
	currTimestamp := int(c.currTimestamp)
	for i := 2; i >= 0; i-- {
		tt[i] = fmt.Sprintf("%.2d", currTimestamp%60)
		currTimestamp /= 60
	}

	return fmt.Sprintf("%s | %s | %d",
		strings.Join(tt, ":"),
		strings.Join(jj, ","),
		c.executingPoint)
}

func (c *Actuator) StringWithPriority() string {
	ss := make([]string, 0, 16)
	for _, job := range c.jobs {
		ss = append(ss, job.StringWithPriority())
	}

	tt := make([]string, 3, 4)
	currTimestamp := int(c.currTimestamp)
	for i := 2; i >= 0; i-- {
		tt[i] = fmt.Sprintf("%.2d", currTimestamp%60)
		currTimestamp /= 60
	}

	return fmt.Sprintf("%s | %s | %d",
		strings.Join(tt, ":"),
		strings.Join(ss, ","),
		c.executingPoint)
}
