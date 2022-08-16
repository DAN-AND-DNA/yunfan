all: test media_api_info_service user_service company_agent_service task_id_service


APP_WHAT=example

.PHONY: clean
clean:
	@rm -rf _output


.PHONY: app-example
app-example:
	@echo "build app-example"
	@go build -o ./_output/app-example/ ./pkg/app-example

.PHONY: media_api_info_service
media_api_info_service:
	@echo "swag media_api_info_service"
	@swag init --dir ./pkg/services/media_api_info_service --output ./pkg/services/media_api_info_service/docs  --parseDependency=true --parseDepth 2
	@echo "build media_api_info_service"
	@CGO_ENABLED=1 go build -buildmode pie -tags "netgo osusergo static_build" -ldflags '-w -extldflags "-fno-PIC -static"' -o ./_output/media_api_info_service/ ./pkg/services/media_api_info_service
	@cp ./pkg/services/media_api_info_service/Dockerfile ./_output/media_api_info_service/

.PHONY: user_service
user_service:
	@echo "swag user_service"
	@swag init --dir ./pkg/services/user_service --output ./pkg/services/user_service/docs  --parseDependency=true --parseDepth 2
	@echo "build user_service"
	@CGO_ENABLED=1 go build -buildmode pie -tags "netgo osusergo static_build" -ldflags '-w -extldflags "-fno-PIC -static"' -o ./_output/user_service/ ./pkg/services/user_service
	@cp ./pkg/services/user_service/Dockerfile ./_output/user_service/

.PHONY: company_agent_service
company_agent_service:
	@echo "swag company_agent_service"
	@swag init --dir ./pkg/services/company_agent_service --output ./pkg/services/company_agent_service/docs  --parseDependency=true --parseDepth 2
	@echo "build company_agent_service"
	@CGO_ENABLED=1 go build -buildmode pie -tags "netgo osusergo static_build" -ldflags '-w -extldflags "-fno-PIC -static"' -o ./_output/company_agent_service/ ./pkg/services/company_agent_service
	@cp ./pkg/services/company_agent_service/Dockerfile ./_output/company_agent_service/


.PHONY: task_id_service
task_id_service:
	@echo "build task_id_service"
	@go build -o ./_output/task_id_service/ ./pkg/services/task_id_service
	@cp ./pkg/services/task_id_service/Dockerfile ./_output/task_id_service/


.PHONY: test
test:
	@go clean -testcache
	@CGO_ENABLED=1 go test -v  ./tests/*

.PHONY: docs
docs:
	@swag init

.PHONY: gc
gc:   
	@GODEBUG=gctrace=1 go run main.go
		          
.PHONY: release
release:
	@CGO_ENABLED=0 go build -buildmode pie -tags "netgo osusergo static_build" -ldflags '-w -extldflags "-fno-PIC -static"'

.PHONY: tree
tree:
	@tree -d -I vendor
