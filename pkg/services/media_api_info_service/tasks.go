package main

// 注册任务

import (
	"yunfan/pkg/services/media_api_info_service/tasks/auto_tasks"
	pkg_tasks "yunfan/pkg/tasks"
)

var (
	ptr_auto_tasks = auto_tasks.New()
)

func init() {
	pkg_tasks.New_task_second(20, 30, 0, false, ptr_auto_tasks.Refresh_toutiao_token)
}
