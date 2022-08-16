package tasks

import (
	"log"
	"sync"
	"time"
)

var (
	default_task_mgr = New_task_mgr()
)

type Task_mgr struct {
	task_id             uint64
	pending_tasks       map[uint64]*task
	m_task_mgr          sync.Mutex
	execute_num         uint64
	is_closing          bool
	current_execute_num uint64
}

type task struct {
	id                  uint64
	delay               int
	is_first_execute    bool
	interval            uint64
	callback            func(uint64) error
	is_idempotent       bool   // 幂等
	max_execute_num     uint64 // 0  一直运行
	current_execute_num uint64
	is_executing        bool
	m_task              sync.Mutex
}

func (this *task) execute(wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	if this.is_first_execute && this.delay != 0 {
		time.Sleep(time.Duration(this.delay) * time.Second)
		this.is_first_execute = false
	}

	err := this.callback(this.current_execute_num)
	if err != nil {
		log.Println(err)
	}

	this.m_task.Lock()
	this.is_executing = false
	this.m_task.Unlock()
}

func New_task_mgr() *Task_mgr {
	ptr_task_mgr := &Task_mgr{
		pending_tasks: make(map[uint64]*task),
	}
	return ptr_task_mgr
}

func (this *Task_mgr) New_task_second(delay int, interval, max_execute_num uint64, is_idempotent bool, callback func(uint64) error) {
	ptr_new_task := &task{
		id:               this.task_id,
		is_first_execute: true,
		callback:         callback,
		interval:         interval,
		is_idempotent:    is_idempotent,
		delay:            delay,
		max_execute_num:  max_execute_num,
	}
	this.m_task_mgr.Lock()
	this.pending_tasks[this.task_id] = ptr_new_task
	this.m_task_mgr.Unlock()
	this.task_id++
}

func (this *Task_mgr) loop() {
	log.Println("task mgr: ok")
	wg := new(sync.WaitGroup)
	done_tasks := make([]uint64, 0, 2000)

	for {
		this.m_task_mgr.Lock()
		for task_id, t := range this.pending_tasks {
			if this.is_closing {
				break
			}

			// is startup ?
			t.m_task.Lock()

			if t.max_execute_num != 0 && t.current_execute_num >= t.max_execute_num {
				if t.is_idempotent || !t.is_executing {
					done_tasks = append(done_tasks, task_id)
				}
				t.m_task.Unlock()
				continue
			}

			if t.is_executing && !t.is_idempotent {
				t.m_task.Unlock()
				continue
			}

			if this.current_execute_num%t.interval != 0 {
				t.m_task.Unlock()
				continue
			}

			t.is_executing = true
			t.m_task.Unlock()

			wg.Add(1)
			t.current_execute_num++

			go t.execute(wg)
		}

		if this.is_closing {
			this.m_task_mgr.Unlock()
			break
		} else {
			for _, task_id := range done_tasks {
				delete(this.pending_tasks, task_id)
			}
			done_tasks = done_tasks[:0]
		}

		this.m_task_mgr.Unlock()
		time.Sleep(1 * time.Second)
		this.current_execute_num++

	}

	wg.Wait()
}

func (this *Task_mgr) Handle_task() {
	go this.loop()
}

func (this *Task_mgr) Shutdown(deadline int) {
	this.m_task_mgr.Lock()
	defer this.m_task_mgr.Unlock()

	this.is_closing = true

	log.Println("task mgr: shutdown")
	time.Sleep(time.Duration(deadline) * time.Second)
}

func New_task_second(delay int, interval, max_execute_num uint64, is_idempotent bool, callback func(uint64) error) {
	default_task_mgr.New_task_second(delay, interval, max_execute_num, is_idempotent, callback)
}

func Handle_task() {
	default_task_mgr.Handle_task()
}

func Shutdown(deadline int) {
	default_task_mgr.Shutdown(deadline)
}
