package util

import (
	"container/heap"
	"github.com/fujiahui/talnet-challenge-payman/common"
)

type CmpHandler func(x any, y any) bool

type JobPriorityQueue struct {
	queueArray []*JobQueue
	queueMap   map[common.PriorityType]*JobQueue

	cmp CmpHandler
}

func NewJobPriorityQueue(cmp CmpHandler) *JobPriorityQueue {

	pq := &JobPriorityQueue{
		queueArray: make([]*JobQueue, 0, 16),
		queueMap:   make(map[common.PriorityType]*JobQueue),
		cmp:        cmp,
	}
	heap.Init(pq)
	return pq
}

func (pq JobPriorityQueue) Len() int { return len(pq.queueArray) }

func (pq JobPriorityQueue) Less(i, j int) bool {
	x := pq.queueArray[i].Front()
	y := pq.queueArray[j].Front()

	return pq.cmp(x, y)
}

func (pq JobPriorityQueue) Swap(i, j int) {
	pq.queueArray[i], pq.queueArray[j] = pq.queueArray[j], pq.queueArray[i]
	pq.queueArray[i].index = i
	pq.queueArray[j].index = j
}

func (pq *JobPriorityQueue) Push(x any) {
	n := len(pq.queueArray)
	item := x.(*JobQueue)
	item.index = n
	pq.queueArray = append(pq.queueArray, item)
	//
	priority := item.priority
	pq.queueMap[priority] = item
}

func (pq *JobPriorityQueue) Pop() any {
	old := pq.queueArray
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	pq.queueArray = old[0 : n-1]
	//
	priority := item.priority
	delete(pq.queueMap, priority)
	return item
}

func (pq *JobPriorityQueue) PopFront() any {

	if pq.Len() == 0 {
		return nil
	}

	q := heap.Pop(pq).(*JobQueue)
	x := q.PopFront()
	if q.Len() > 0 {
		heap.Push(pq, q)
	}

	return x
}

func (pq *JobPriorityQueue) PushBack(priority common.PriorityType, x any) {
	if q, ok := pq.at(priority); !ok {
		q = NewJobQueue(priority)
		q.PushBack(x)
		heap.Push(pq, q)
	} else {
		q.PushBack(x)
	}
}

func (pq *JobPriorityQueue) PushFront(priority common.PriorityType, x any) {
	if q, ok := pq.at(priority); !ok {
		q = NewJobQueue(priority)
		q.PushFront(x)
		heap.Push(pq, q)
	} else {
		q.PushFront(x)
		heap.Fix(pq, q.index)
	}
}

func (pq *JobPriorityQueue) at(priority common.PriorityType) (*JobQueue, bool) {
	q, ok := pq.queueMap[priority]
	return q, ok
}
