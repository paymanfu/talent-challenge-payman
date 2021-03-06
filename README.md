# 一、模块
## 1.1 Server 模块
    Server加载目录`./ware/house/`下的 'data/' 或者 'data_num_priority/'的文件，
    并解析成 *JobInfo* 结构体存储。
    此外，提供 GetJobInfo API接口获取Job工作数据。
## 1.2 Worker 模块
    Worker对Job数据进行调度并执行，核心模块由以下两个组件组成：
    a. Actuator: Task执行管理器，执行最小粒度为Task。
    b. Scheduler: Job调度管理器，执行最小粒度为Job。
    Worker分别支持以下几种创建方式：
    1. NewBaseWorker: 创建一个无容量、忽略优先级的Worker；
    2. NewWorkerWithCapacity: 创建一个有容量、忽略优先级的Worker；
    3. NewWorkerWithSimplePriority: 创建一个有容量、优先级有效的Worker，Worker采用'期待执行时间相同的情况下，才比较优先级'；
    4. NewWorkerWithSmartPriority: 创建一个有容量、智能优先级的Worker，Worker采用'期待执行时间加权优先级比较排序'。

# 二、组件
## 2.1 Task执行管理器 Actuator
Task执行管理器 [actuator.go](./worker/actuator.go) 核心工作是执行满足条件的Task，核心条件为capacity容量。

    1. 方法 Ticking() 执行Task
    2. 方法 String() 输出图标功能

## 2.2 Job调度管理器 Scheduler
Scheduler调度管理器 [scheduler.go](./worker/scheduler.go) 基于**优先队列Priority Queue**进行Job调度。

    1. Enqueue入队: 方法用于将一个Job压入优先队列中；
    2. Enqueue出队: 方法用于从优先队列中获取一个优先级最高的Job。

此外，优先队列的实现在 [priority_queue.go](./worker/util/priority_queue.go) 文件中，基于go容器 [heap](https://pkg.go.dev/container/heap@go1.18.2) 实现优先队列。队列的每个元素是一个go容器 [list](https://pkg.go.dev/container/list@go1.18.2)，同一优先级的Job进入同一个 [list](https://pkg.go.dev/container/list@go1.18.2)，从而保证同一优先级的Job有先后顺序。


Task状态转移
    
    TaskWait: 等待被执行
	TaskRunning: 正在被Actuator执行
	TaskFinished: 执行完成

Job状态转移
    
    JobCreated  = Job已创建
	JobWait     = Job等待执行
	JobRunning  = Job正在执行
	JobSleep    = Job休眠中
	JobFinished = Job执行完成

![](job状态转移.png)

# 三、测试用例
    
笔试的Task在测试用例 [smart_worker_test.go](./worker/smart_worker_test.go)文件中实现。
    
    分别实现了 Task 1.2、Task 2.1、Task 2.2、Task 2.3和Task 3.1、Task 3.3。
    
    0. DataHubServer的方法GetJobInfo: 实现 Task 1.1;
    1. TestNewBaseWorker: 实现Task 1.2;
    2. TestNewWorkerWithCapacity: 实现 Task 2.1;
    3. TestNewWorkerWithSimplePriority: 实现 Task 2.2;
    4. TestNewWorkerWithSmartPriority: 实现 Task 2.3;
    5. TestNewWorkerWithNumPriority: 实现 Task 3.1;
    6. TestNewWorkerWithTaskSpeed: 实现 Task 3.3。