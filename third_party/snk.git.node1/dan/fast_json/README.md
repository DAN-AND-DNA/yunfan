# fast_json

稍微快的json库

## 基准

1. 运行 make 命令 

```
$ make
go test -v fast_json.go fast_json_test.go
=== RUN   Test_Marshal_simple
--- PASS: Test_Marshal_simple (0.00s)
=== RUN   Test_Marshal
--- PASS: Test_Marshal (0.00s)
=== RUN   Test_Unmarshal_simple
--- PASS: Test_Unmarshal_simple (0.00s)
=== RUN   Test_Unmarshal
--- PASS: Test_Unmarshal (0.00s)
PASS
ok      command-line-arguments  0.004s
go test -v -bench=. -benchtime=3s fast_json_timing_test.go fast_json.go
goos: linux
goarch: amd64
Benchmark_Marshal_simple
Benchmark_Marshal_simple-8      10227925               326 ns/op              64 B/op          1 allocs/op
Benchmark_Marshal_en
Benchmark_Marshal_en-8          15957931               188 ns/op               0 B/op          0 allocs/op
Benchmark_Marshal_cn
Benchmark_Marshal_cn-8          17812866               210 ns/op               0 B/op          0 allocs/op
Benchmark_Unmarshal_simple
Benchmark_Unmarshal_simple-8     4861639               696 ns/op              32 B/op          3 allocs/op
Benchmark_Unmarshal_en
Benchmark_Unmarshal_en-8         7669627               475 ns/op               0 B/op          0 allocs/op
Benchmark_Unmarshal_cn
Benchmark_Unmarshal_cn-8         3631740               868 ns/op              24 B/op          3 allocs/op
PASS
ok      command-line-arguments  28.908s

```
