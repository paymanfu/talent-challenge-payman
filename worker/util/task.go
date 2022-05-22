package util

import (
	"fmt"
	"github.com/fujiahui/talnet-challenge-payman/common"
	"time"
)

type TaskStatusType uint8

const (
	TaskWait     = TaskStatusType(1)
	TaskRunning  = TaskStatusType(2)
	TaskFinished = TaskStatusType(4)
)

type Task struct {
	jobID       int64
	taskID      int
	taskPoint   common.PointType // Task总的要执行的Points数
	remainPoint common.PointType // Task剩余要执行的Points数
	/*
		预计开始的时间
			1. 刚创建的Job的第一个Task的 ExpectedTime = Job.Created
			2. Job的其他Task 等于前一个Task结束时间 +1s
	*/
	expectedTimestamp int64
	status            TaskStatusType

	speedFlag bool
}

func NewTask(point common.PointType) *Task {
	return &Task{
		taskPoint:   point,
		remainPoint: point,
		status:      TaskWait,
	}
}

func (t *Task) Ticking() {
	PointOfTicks := common.PointType(0)
	if t.speedFlag && t.remainPoint > 0 && t.taskPoint%2 == 0 {
		PointOfTicks = common.PointType(2)
	} else {
		PointOfTicks = common.PointType(1)
	}

	t.remainPoint -= PointOfTicks

	if t.remainPoint == 0 {
		t.status = TaskFinished
	}

	return
}

// TaskPoint 用于调度器中进行比较排序
func (t *Task) TaskPoint() common.PointType {
	if t.speedFlag && t.remainPoint > 0 && t.taskPoint%2 == 0 {
		return t.taskPoint / 2
	}
	return t.taskPoint
}

// RemainPoint 用于图标的输出
func (t *Task) RemainPoint() common.PointType {
	if t.speedFlag && t.remainPoint > 0 && t.taskPoint%2 == 0 {
		return t.remainPoint / 2
	}
	return t.remainPoint
}

func (t *Task) Status() TaskStatusType {
	return t.status
}

func (t *Task) EnableSpeed() {
	t.speedFlag = true
}

func (t *Task) DisableSpeed() {
	t.speedFlag = false
}

func (t *Task) SetJobID(jobID int64) {
	t.jobID = jobID
}

func (t *Task) SetTaskID(taskID int) {
	t.taskID = taskID
}

func (t *Task) SetRunning() {
	t.status = TaskRunning
}

func (t *Task) SetExpectedTime(currTimestamp int64) {
	t.expectedTimestamp = currTimestamp
}

func (t *Task) ExpectedTimestamp() int64 {
	return t.expectedTimestamp
}

func (t *Task) Finished() bool {
	return t.status == TaskFinished
}

func (t *Task) String() string {
	return fmt.Sprintf("%d(%d)", t.taskID, t.remainPoint)
}

func (t *Task) StringWithExpectedTime() string {
	return fmt.Sprintf("%d(%d-%d)", t.taskID, t.expectedTimestamp, t.remainPoint)
}

func (t *Task) Running(tick int) {
	remainPoint := t.remainPoint
	for point := remainPoint; point > 0; point-- {
		time.Sleep(time.Duration(tick) * time.Millisecond)
		fmt.Printf("JobRunning %d-%d(%d)", t.jobID, t.taskID, point)
	}
}
