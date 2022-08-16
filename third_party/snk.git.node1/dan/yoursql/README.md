# yoursql

## 基准

```
=== RUN   Test_Query
--- PASS: Test_Query (0.00s)
=== RUN   Test_GORM_find
--- PASS: Test_GORM_find (0.00s)
=== RUN   Test_GORM_first
--- PASS: Test_GORM_first (0.00s)
=== RUN   Test_GORM_first_where
--- PASS: Test_GORM_first_where (0.00s)
goos: linux
goarch: amd64
pkg: snk.git.node1/dan/yoursql
cpu: Intel(R) Core(TM) i5-9600K CPU @ 3.70GHz
Benchmark_Query
Benchmark_Query-6        	 1877784	       617.5 ns/op	     702 B/op	      15 allocs/op
Benchmark_GORM_first
Benchmark_GORM_first-6   	  716335	      1607 ns/op	    3442 B/op	      43 allocs/op
PASS
ok  	snk.git.node1/dan/yoursql	2.980s
```
