package main

// 注册任务
import (
	"yunfan/pkg/app-example/tasks"
	pkg_tasks "yunfan/pkg/tasks"
)

var (
	ptr_example_task = tasks.New_example_task()
)

func init() {
	pkg_tasks.New_task_second(0, 1, 3, false, ptr_example_task.Do_local_log)
}
