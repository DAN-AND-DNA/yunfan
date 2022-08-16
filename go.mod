module yunfan

go 1.16

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-kit/kit v0.10.0
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/gofiber/fiber/v2 v2.13.0
	github.com/json-iterator/go v1.1.11
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/http-swagger v1.0.0
	github.com/swaggo/swag v1.7.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/xtaci/kcp-go/v5 v5.6.1
	go.uber.org/zap v1.18.1
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/tools v0.1.5 // indirect
	gorm.io/driver/clickhouse v0.2.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
	snk.git.node1/dan/fast_json v0.0.0 // indirect
	snk.git.node1/dan/go_request v0.0.0
	snk.git.node1/dan/yoursql v0.0.0-00010101000000-000000000000
	snk.git.node1/yunfan/arpc v0.0.0
	snk.git.node1/yunfan/yunfan_dev v0.0.0 // indirect
)

// replace 项目名 => 项目位置
replace (
	snk.git.node1/dan/fast_json => ./third_party/snk.git.node1/dan/fast_json
	snk.git.node1/dan/go_request => ./third_party/snk.git.node1/dan/go_request
	snk.git.node1/dan/yoursql => ./third_party/snk.git.node1/dan/yoursql
	snk.git.node1/yunfan/arpc => ./third_party/snk.git.node1/yunfan/arpc
	snk.git.node1/yunfan/yunfan_dev => ./third_party/snk.git.node1/yunfan/yunfan_dev
)
