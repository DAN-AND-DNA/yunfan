# yunfan
学习使用，玩具微服务框架："云帆"，了解一下k8s的基础服务开发，简单的rpc，swag文档输出，今日头条接口对接等功能

## 使用方法

1. 安装kube参考(http://snk.git.node1/dan/ops/src/branch/master/gen_kube_install.sh)
2. 登陆华为云镜像SWR (https://console.huaweicloud.com/swr/xxxxxxxxxxxxxxxxxxx)
3. 在kube的工作节点设置docker或containerd登陆密码
4. 拷贝开发版本配置如: cp ./k8s/development ./k8s/my-work
5. 启动: kubectl apply -f ./k8s/my-work

## 在页面建立分支后，如何使用git跟踪

```sh
# 更新服务端的分支到本地仓库
git fetch --all

# 列出所有的分支，如remotes/origin/develop
git branch --all

# 切换到一个已有的分支，如develop
git checkout <你建立的分支名>

# 删除本地的一个分支，如develop
git branch -D <本地分支名>

# 删除服务端的一个分支，如develop
git push origin :<远程分支>
```

## 如何开发一个服务

1. 在pkg/services里建立服务
2. 在sdk/arpc里建立服务的rpc协议
3. 在sdk/dbs里建立服务的数据库结构
4. 在sdk/cmd里建立服务命令
5. 在Makefile里添加服务的构建脚本
6. 例子参考user_service和media_api_info_service


## 基准

```
=== RUN   Test_gob_rpc_call
--- PASS: Test_gob_rpc_call (0.00s)
=== RUN   Test_json_rpc_call
--- PASS: Test_json_rpc_call (0.00s)
goos: linux
goarch: amd64
pkg: template-project/benchmark
cpu: Intel(R) Core(TM) i5-9600K CPU @ 3.70GHz
Benchmark_server_handle_json_rpc_call
Benchmark_server_handle_json_rpc_call-6   	 1291899	       922.6 ns/op	    1960 B/op	      25 allocs/op
Benchmark_server_handle_gob_rpc_call
Benchmark_server_handle_gob_rpc_call-6    	  195685	      5492 ns/op	   16828 B/op	     220 allocs/op
PASS
```

## 测试

```
=== RUN   Test_client_json_rpc_call_easy
--- PASS: Test_client_json_rpc_call_easy (0.00s)
=== RUN   Test_client_json_rpc_call_map
--- PASS: Test_client_json_rpc_call_map (0.00s)
=== RUN   Test_client_json_rpc_call_bytes
--- PASS: Test_client_json_rpc_call_bytes (0.00s)
PASS
```
