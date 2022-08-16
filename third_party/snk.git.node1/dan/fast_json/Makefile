all: test benchmark

.PHONY: test
test:
	go test -v fast_json.go fast_json_test.go


.PHONY: benchmark
benchmark:
	go test -v -bench=. -benchtime=3s fast_json_timing_test.go fast_json.go
