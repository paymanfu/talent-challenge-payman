package worker

import (
	"context"
	"github.com/fujiahui/talnet-challenge-payman/common"
	"github.com/fujiahui/talnet-challenge-payman/logger"
	"github.com/fujiahui/talnet-challenge-payman/worker/util"
	"sync"
)

type handler func(created common.TimestampType) *common.JobInfoArray

type SmartWorker struct {
	// Job调度管理器
	scheduler *Scheduler
	// Job执行管理器
	actuator *Actuator

	speedFlag bool
}

// NewBaseWorker Task 1.2
func NewBaseWorker(startTimestamp common.TimestampType) *SmartWorker {
	capacity := common.MaxCapacity // 16位最大整数 == 不限制容量
	return &SmartWorker{
		scheduler: NewScheduler(false, util.BaseCmp),
		actuator:  NewActuator(startTimestamp, capacity),
	}
}

// NewWorkerWithCapacity Task 2.1
func NewWorkerWithCapacity(startTimestamp common.TimestampType, capacity common.PointType) *SmartWorker {
	return &SmartWorker{
		scheduler: NewScheduler(false, util.BaseCmp),
		actuator:  NewActuator(startTimestamp, capacity),
	}
}

// NewWorkerWithSimplePriority Task 2.2
func NewWorkerWithSimplePriority(startTimestamp common.TimestampType, capacity common.PointType) *SmartWorker {
	return &SmartWorker{
		scheduler: NewScheduler(true, util.SimpleCmp),
		actuator:  NewActuator(startTimestamp, capacity),
	}
}

// NewWorkerWithSmartPriority Task 2.3
func NewWorkerWithSmartPriority(startTimestamp common.TimestampType, capacity common.PointType) *SmartWorker {
	return &SmartWorker{
		scheduler: NewScheduler(true, util.SmartCmp),
		actuator:  NewActuator(startTimestamp, capacity),
	}
}

func (w *SmartWorker) EnableTaskSpeed() {
	w.speedFlag = true
}

func (w *SmartWorker) DisableTaskSpeed() {
	w.speedFlag = false
}

func (w *SmartWorker) Start(ctx context.Context, h handler) {
	wg := &sync.WaitGroup{}

	tick := 100
	for {
		// 0. 按照时间进行迭代图标
		jobs := w.actuator.Ticking(tick)

		// 1. 每隔1ms / 1s 获取一批次的job列表
		if jobArray := h(w.actuator.CurrTimestamp()); jobArray != nil {
			for _, info := range jobArray.JobInfos {

				tooManyPoint := false
				for _, task := range info.Tasks {
					point := task
					if w.speedFlag && point%2 == 0 {
						point /= 2
					}
					if task > w.actuator.Capacity() {
						logger.Warnf("%d-%d's task point more than capacity %d", info.ID, task, w.actuator.Capacity())
						tooManyPoint = true
						break
					}
				}
				if tooManyPoint {
					continue
				}
				job := util.NewJobFromCommon(info)
				if job == nil {
					continue
				}

				if w.speedFlag {
					job.EnableTaskSpeed()
				} else {
					job.DisableTaskSpeed()
				}
				jobs = append(jobs, job)
			}
		}

		// 2. 把jobs放入优先队列中
		for _, job := range jobs {
			w.scheduler.Enqueue(job)
		}

		// 3. 分发Task
		for {
			freePoint := w.actuator.FreePoint()
			job := w.scheduler.Dequeue(freePoint)
			if job == nil {
				break
			}
			w.actuator.Execute(job)

			wg.Add(1)
			go func(t *util.Task) {
				defer wg.Done()
				t.Running(ctx, tick)
			}(job.CurrTask())

		}

		logger.ChartLogger.Println(w.actuator.String())
		if ctx.Err() != nil {
			logger.Errorf("Worker Exit, error %v", ctx.Err())
			break
		}

	}

	wg.Wait()
}
