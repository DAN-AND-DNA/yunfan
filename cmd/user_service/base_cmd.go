package user_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ptr_env_runtime = &Env_runtime{}
	base_cmd        = &cobra.Command{
		Use:   "user_service [flags]",
		Short: "用户服务",
		RunE: func(_ *cobra.Command, _ []string) error {
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

			if ptr_env_runtime.Arpc_conn_idle < 0 {
				return errors.New("the idle time of connection should bigger than 0")
			}

			if ptr_env_runtime.Arpc_network_protocol == "udp" {
				if ptr_env_runtime.Arpc_udp_key == "" {
					return errors.New("need udp key")
				}

				if ptr_env_runtime.Arpc_udp_salt == "" {
					return errors.New("need udp salt")
				}
			}

			if ptr_env_runtime.Arpc_swag != "" {
				switch ptr_env_runtime.Arpc_swag {
				case "on":
				case "off":
				default:
					return errors.New("bad swag option")
				}
			}

			return nil
		},
	}
)

type Env_runtime struct {
	DB_postgresql_dsn     string `json:"db_postgresql_dsn"`
	Log_level             string `json:"log_level"`
	Arpc_network_protocol string `json:"arpc_network_protocol"`
	Arpc_app_protocol     string `json:"arpc_app_protocol"`
	Arpc_listen_port      string `json:"arpc_listen_port"`
	Arpc_udp_key          string `json:"arpc_udp_key"`
	Arpc_udp_salt         string `json:"arpc_udp_salt"`
	Arpc_conn_idle        int    `json:"arpc_conn_idle"`
	Arpc_swag             string `json:"arpc_swag"`
}

func Execute() {
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.DB_postgresql_dsn), "db-postgresql-dsn", "", "host=localhost user=dan password=12345678 dbname=dan port=5432 sslmode=disable TimeZone=Asia/Shanghai", "postgresql的dsn")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Log_level), "log-level", "", "", "日志等级：[debug | info | warn | error]")
	base_cmd.MarkFlagRequired("log-level")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Arpc_network_protocol), "arpc-network-protocol", "", "tcp", "arpc网络层协议类型：[tcp | udp]")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Arpc_app_protocol), "arpc-app-protocol", "", "json", "arpc应用层协议类型：[gob | json]")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Arpc_listen_port), "arpc-listen-port", "", "", "arpc监听端口")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Arpc_udp_key), "arpc-udp-key", "", "", "udp协议密钥")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Arpc_udp_salt), "arpc-udp-salt", "", "", "udp协议盐")
	base_cmd.Flags().IntVarP(&(ptr_env_runtime.Arpc_conn_idle), "arpc-conn-idle", "", 30, "arpc连接保持时间")
	base_cmd.Flags().StringVarP(&(ptr_env_runtime.Arpc_swag), "arpc-swag", "", "off", "是否开启arpc的swag：[off | on]")

	base_cmd.SetHelpFunc(func(*cobra.Command, []string) {
		fmt.Println(`
用户服务

Usage:
  user_service [flags]

Flags:
      --arpc-app-protocol string            arpc应用层协议类型：[gob | json] (default "json")
      --arpc-conn-idle int                  arpc连接保持时间 (default 30)
      --arpc-listen-port string             arpc监听端口
      --arpc-network-protocol string        arpc网络层协议类型：[tcp | udp] (default "tcp")
      --arpc-swag string                    是否开启arpc的swag：[off | on] (default "off")
      --arpc-udp-key string                 udp协议密钥
      --arpc-udp-salt string                udp协议盐
      --db-postgresql-dsn string            postgresql的dsn (default "host=localhost user=dan password=12345678 dbname=dan port=5432 sslmode=disable TimeZone=Asia/Shanghai")
  -h, --help                                help for user_service
      --log-level string                    日志等级：[debug | info | warn | error]
`)
		os.Exit(0)
	})

	if err := base_cmd.Execute(); err != nil {
		os.Exit(1)
	}

	str_env_runtime, err := json.MarshalIndent(ptr_env_runtime, "", "	")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(str_env_runtime))
}

func Get_env() *Env_runtime {
	return ptr_env_runtime
}
