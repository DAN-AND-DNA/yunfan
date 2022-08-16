package log

import (
	"errors"
	"time"

	"sync"

	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	default_bind_log = New_bind_log()

	ERR_HAS_LOCAL  = errors.New("has register as local")
	ERR_HAS_REMOTE = errors.New("has register as remote")
	ERR_BAD_LEVEL  = errors.New("bad level")
)

type Bind_log struct {
	local_log       *zap.Logger
	local_log_sugar *zap.SugaredLogger
	m_log           sync.RWMutex
	has_local       bool
	zap_level       zapcore.Level
}

func New_bind_log() *Bind_log {
	return &Bind_log{
		has_local: false,
		zap_level: zap.DebugLevel,
	}
}

func (this *Bind_log) Connect(level string) error {
	this.m_log.Lock()
	defer this.m_log.Unlock()

	if this.has_local {
		return nil
	}

	this.zap_level = zap.DebugLevel
	switch level {
	case "debug":
		this.zap_level = zap.DebugLevel
	case "info":
		this.zap_level = zap.InfoLevel
	case "warn":
		this.zap_level = zap.WarnLevel
	case "error":
		this.zap_level = zap.ErrorLevel
	default:
		return ERR_BAD_LEVEL
	}

	atom := zap.NewAtomicLevel()
	atom.SetLevel(this.zap_level)

	cfg := zap.Config{
		Level:            atom,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "",
			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	this.local_log = logger
	this.local_log_sugar = this.local_log.Sugar()
	this.has_local = true

	go this.Report_runtime_metrics(37) // 37s
	return nil
}

func (this *Bind_log) Debug(vals []interface{}) {
	this.m_log.RLock()
	defer this.m_log.RUnlock()
	if this.has_local {
		this.local_log_sugar.Debugw("", vals...)
	}
}

func (this *Bind_log) Info(vals []interface{}) {
	this.m_log.RLock()
	defer this.m_log.RUnlock()

	if this.has_local {
		this.local_log_sugar.Infow("", vals...)
	}
}

func (this *Bind_log) Warn(vals []interface{}) {
	this.m_log.RLock()
	defer this.m_log.RUnlock()

	if this.has_local {
		this.local_log_sugar.Warnw("", vals...)
	}
}

func (this *Bind_log) Error(vals []interface{}) {
	this.m_log.RLock()
	defer this.m_log.RUnlock()

	if this.has_local {
		this.local_log_sugar.Errorw("", vals...)
	}
}

func (this *Bind_log) Shutdown(wait_s int) {
	this.m_log.Lock()
	defer this.m_log.Unlock()

	if !this.has_local {
		return
	}

	this.local_log.Sync()
	time.Sleep(time.Duration(wait_s) * time.Second)

	this.has_local = false
}

func (this *Bind_log) Report_runtime_metrics(interval int) {

	for {
		time.Sleep(time.Duration(interval) * time.Second)

		this.m_log.RLock()
		if !this.has_local {
			this.m_log.RUnlock()
			return
		}

		m := runtime.MemStats{}
		runtime.ReadMemStats(&m)

		this.local_log_sugar.Infow("", "===start_time_s===", time.Now().Unix())
		this.local_log_sugar.Infow("", "协程数", runtime.NumGoroutine())
		this.local_log_sugar.Infow("", "从操作系统申请作为go运行时的虚拟内存 KiB", m.Sys/(1024))
		this.local_log_sugar.Infow("", "从操作系统申请作为堆的虚拟内存 KiB", m.HeapSys/(1024))
		this.local_log_sugar.Infow("", "已经从堆里归还给操作系统的内存 KiB", m.HeapReleased/(1024))
		this.local_log_sugar.Infow("", "堆中分配的对象数", m.HeapObjects)
		this.local_log_sugar.Infow("", "堆的使用量 KiB", m.HeapAlloc/(1024))
		this.local_log_sugar.Infow("", "堆的占用量 KiB", m.HeapInuse/(1024))
		this.local_log_sugar.Infow("", "堆的空闲量 KiB", m.HeapIdle/(1024))
		this.local_log_sugar.Infow("", "触发下次GC的堆目标大小 KiB", m.NextGC/(1024))
		this.local_log_sugar.Infow("", "总GC次数", m.NumGC)
		this.local_log_sugar.Infow("", "应用强制GC", m.NumForcedGC)
		this.local_log_sugar.Infow("", "上次GC暂停时间 s", (float64)(m.PauseNs[(m.NumGC+255)%256])/(float64)(1000000000))
		this.local_log_sugar.Infow("", "总GC暂停时间 s", (float64)(m.PauseTotalNs)/(float64)(1000000000))
		this.m_log.RUnlock()
	}
}

func Connect(level string) error {
	return default_bind_log.Connect(level)
}

func Shutdown(wait_s int) {
	default_bind_log.Shutdown(wait_s)
}

func Debug(vals ...interface{}) {
	default_bind_log.Debug(vals)
}

func Info(vals ...interface{}) {
	default_bind_log.Info(vals)
}

func Warn(vals ...interface{}) {
	default_bind_log.Warn(vals)
}

func Error(vals ...interface{}) {
	default_bind_log.Error(vals)
}
