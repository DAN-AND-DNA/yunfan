# httpclient-go

fasthttp 封装的 http client, 适合http rpc或者api调用

## 版本

0.0.1

## 用法

参考test, 尽可能使用Send方法

## 测试
```
Benchmark_Get_json-4             7706832               741 ns/op             165 B/op          7 allocs/op
Benchmark_Get_json_unsafe-4      9295737               650 ns/op             104 B/op          4 allocs/op
Benchmark_Get_json-8            10144230               554 ns/op             166 B/op          7 allocs/op
Benchmark_Get_json_unsafe-8     12676928               512 ns/op             104 B/op          4 allocs/op
```


