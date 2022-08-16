package errcode

const (
	Code_s2s_ok = 0

	// for service to service
	// 1000 -- 1999
	Code_s2s_unkown_error = 1000
	Code_s2s_need_arg     = 1001
	Code_s2s_bad_arg      = 1002
	Code_s2s_network      = 1003
	Code_s2s_retry        = 1004

	// for db
	// 2000 -- 2999
	Code_db_internal_error      = 2000
	Code_db_no_found_error      = 2001
	Code_db_node_internal_error = 2002
)

const (
	From_local = 0

	// for app_example
	// 1000 -- 1999
	From_app_example = 1000

	// for media_api_info_service
	// 2000 -- 2999
	From_media_api_info_service = 2000
	From_create_token           = 2001
	From_manual_refresh_token   = 2002
	From_auto_refresh_task      = 2003

	// for user_service
	// 3000 -- 3999
	From_user_service   = 3000
	From_create_system  = 3001
	From_create_company = 3002
	From_create_user    = 3003
	From_gen_id         = 3004

	// for company_agent_service
	// 4000 --- 4999
	From_company_agent_service = 4000
	From_reopen_db             = 4001

	// for task_id_service
	// 5000 -- 5999
	From_task_id_service = 5000
)
