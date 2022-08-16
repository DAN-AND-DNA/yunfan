package app_example

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type Env_runtime struct {
	Run_as                        string `json:"run_as"`
	Http_listen_port              string `json:"http_listen_port"`
	Http_allow_origins            string `json:"http_allow_origins"`
	Http_header_name              string `json:"http_header_name"`
	Http_body_limit_bytes         int    `json:"http_body_limit_bytes"`
	Http_read_timeout_secs        int    `json:"http_read_timeout_secs"`
	Http_write_timeout_secs       int    `json:"http_write_timeout_secs"`
	Http_max_conns                int    `json:"http_max_conns"`
	DB_clickhouse_dns             string `json:"db_clickhouse_dns"`
	Rpc_listen_port               string `json:"rpc_listen_port"`
	Rpc_protocol_type             string `json:"rpc_protocol_type"`
	Rpc_key                       string `json:"rpc_key"`
	Rpc_salt                      string `json:"rpc_salt"`
	Rpc_idle                      int    `json:"rpc_idle"`
	Rpc_app_protocol_type         string `json:"rpc_app_protocol_type"`
	Log_level                     string `json:"log_level"`
	Log_remote_heartbeat_interval int    `json:"log_remote_heartbeat_interval"`
	Log_remote_flush_interval     int    `json:"log_remote_flush_interval"`
	Log_remote_host               string `json:"log_remote_host"`
	Log_remote_port               string `json:"log_remote_port"`
	Tracer_service_name           string `json:"tracer_service_name"`
	Tracer_remote_hostname        string `json:"tracer_remote_hostname"`
	Tracer_remote_port            string `json:"tracer_remote_port"`
}

var (
	ptr_env_runtime = &Env_runtime{}

	base_cmd = &cobra.Command{
		Use:   "template-project [flags]",
		Short: "复制该模板项目，然后参考例子修改",
		RunE: func(_ *cobra.Command, _ []string) error {
			if ptr_env_runtime.Run_as == "" {
				return errors.New("need run as")
			}

			switch ptr_env_runtime.Run_as {
			case "deployment":
			case "production":
			case "test":
			default:
				return errors.New("bad run as")
			}

			if ptr_env_runtime.Http_listen_port != "" {
				if _, err := strconv.Atoi(ptr_env_runtime.Http_listen_port); err != nil {
					return errors.New("bad http port")
				}
			}

			if len(ptr_env_runtime.Http_allow_origins) == 0 {
				return errors.New("bad all origins")
			}

			if ptr_env_runtime.Http_body_limit_bytes < 0 {
				return errors.New("bad body limit bytes")
			}

			if ptr_env_runtime.Http_read_timeout_secs < 0 {
				return errors.New("bad read timeout")
			}

			if ptr_env_runtime.Http_write_timeout_secs < 0 {
				return errors.New("bad write timeout")
			}

			if ptr_env_runtime.Http_max_conns < 0 {
				return errors.New("bad max conns")
			}

			if ptr_env_runtime.Rpc_listen_port != "" {
				if _, err := strconv.Atoi(ptr_env_runtime.Rpc_listen_port); err != nil {
					return errors.New("bad rpc port")
				}

				if ptr_env_runtime.Rpc_protocol_type != "" {
					switch ptr_env_runtime.Rpc_protocol_type {
					case "tcp":
					case "udp":
						if ptr_env_runtime.Rpc_key == "" || ptr_env_runtime.Rpc_salt == "" {
							return errors.New("need rpc key and rpc salt for safty")
						}

					default:
						return errors.New("bad protocol type")
					}
				}

				if ptr_env_runtime.Rpc_idle < 0 {
					return errors.New("bad rpc idle time")
				}
			}

			if ptr_env_runtime.Rpc_app_protocol_type != "" {
				switch ptr_env_runtime.Rpc_app_protocol_type {
				case "json":
				case "gob":
				default:
					return errors.New("bad rcp app protocol type")
				}
			}

			if ptr_env_runtime.Log_level != "" {
				switch ptr_env_runtime.Log_level {
				case "debug":
				case "info":
				case "warn":
				case "error":
				default:
					return errors.New("bad log level")
				}
			}

			if ptr_env_runtime.Log_remote_heartbeat_interval <= 0 {
				return errors.New("bad log service heartbeat interval")
			}

			if ptr_env_runtime.Log_remote_flush_interval <= 0 {
				return errors.New("bad log service flush interval")
			}

			return nil
		},
	}
)

func Execute() {
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Run_as), "run-as", "", "", "运行模式：[deployment | production | test]")
	base_cmd.MarkFlagRequired("run-as")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Http_listen_port), "http-listen-port", "", "", "http监听端口")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Http_allow_origins), "http-allow-origns", "", "*", "跨域允许域名，以逗号分割")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Http_header_name), "http-header-name", "", "template-server", "请求")
	base_cmd.Flags().IntVarP(&(ptr_env_runtime.Http_body_limit_bytes), "http-body-limit-bytes", "", 1*1024*1024, "http body 大小")
	base_cmd.Flags().IntVarP(&(ptr_env_runtime.Http_read_timeout_secs), "http-read-timeout-secs", "", 10, "连接读超时")
	base_cmd.Flags().IntVarP(&(ptr_env_runtime.Http_write_timeout_secs), "http-write-timeout-secs", "", 10, "连接写超时")
	base_cmd.Flags().IntVarP(&(ptr_env_runtime.Http_max_conns), "http-max-conns", "", 10000, "最大连接数")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.DB_clickhouse_dns), "db-clickhouse-dns", "", "", "clickhouse的dns")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Rpc_listen_port), "rpc-listen-port", "", "", "rpc监听端口")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Rpc_protocol_type), "rpc-protocol-type", "", "udp", "rpc通讯协议：[tcp | udp]")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Rpc_app_protocol_type), "rpc-app-protocol-type", "", "gob", "rpc通讯协议：[gob | json]")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Rpc_key), "rpc-key", "", "", "rpc key")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Rpc_salt), "rpc-salt", "", "", "rpc salt")
	base_cmd.Flags().IntVarP(&(ptr_env_runtime.Rpc_idle), "rpc-idle", "", 30, "rpc连接保持时间")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Log_level), "log-level", "", "debug", "日志等级：[debug | info | warn | error]")
	base_cmd.Flags().IntVarP(&(ptr_env_runtime.Log_remote_heartbeat_interval), "log-remote-heartbeat-interval", "", 5, "远程日志服务的心跳间隔")
	base_cmd.Flags().IntVarP(&(ptr_env_runtime.Log_remote_flush_interval), "log-remote-flush-interval", "", 10, "远程日志服务的推送间隔")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Log_remote_host), "log-remote-host", "", "", "远程日志服务的hostname")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Tracer_service_name), "tracer-service-name", "", "", "追踪服务名")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Tracer_remote_hostname), "tracer-remote-hostname", "", "", "追踪服务器hostname")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Tracer_remote_port), "tracer-remote-port", "", "", "追踪服务器端口")

	base_cmd.SetHelpFunc(func(*cobra.Command, []string) {
		fmt.Println(`
Usage:
  template-project [flags]

Flags:
      --db-clickhouse-dns string            clickhouse的dns
  -h, --help                                help for template-project
      --http-allow-origns string            跨域允许域名，以逗号分割 (default "*")
      --http-body-limit-bytes int           http body 大小 (default 1048576)
      --http-header-name string             请求 (default "template-server")
      --http-listen-port string             http监听端口
      --http-max-conns int                  最大连接数 (default 10000)
      --http-read-timeout-secs int          连接读超时 (default 10)
      --http-write-timeout-secs int         连接写超时 (default 10)
      --log-level string                    日志等级：[debug | info | warn | error] (default "debug")
      --log-remote-flush-interval int       远程日志服务的推送间隔 (default 10)
      --log-remote-heartbeat-interval int   远程日志服务的心跳间隔 (default 5)
      --log-remote-host string              远程日志服务的hostname
      --rpc-app-protocol-type string        rpc通讯协议：[gob | json] (default "gob")
      --rpc-idle int                        rpc连接保持时间 (default 30)
      --rpc-key string                      rpc key
      --rpc-listen-port string              rpc监听端口
      --rpc-protocol-type string            rpc通讯协议：[tcp | udp] (default "udp")
      --rpc-salt string                     rpc salt
      --run-as string                       运行模式：[deployment | production | test]
      --tracer-remote-hostname string       追踪服务器hostname
      --tracer-remote-port string           追踪服务器端口
      --tracer-service-name string          追踪服务名
`)
		os.Exit(0)
	})

	if err := base_cmd.Execute(); err != nil {
		os.Exit(1)
	}

	str_env_runtime, _ := json.MarshalIndent(ptr_env_runtime, "", "	")

	fmt.Println(string(str_env_runtime))
}

func Get_env() *Env_runtime {
	return ptr_env_runtime
}
