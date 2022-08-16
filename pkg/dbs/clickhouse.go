package dbs

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

var (
	default_ck_client = New_clickhouse_client()

	ERR_OPEN_FIRST = errors.New("need open db first")
)

type Clickhouse_client struct {
	db        *gorm.DB
	is_open   bool
	db_dsn    string
	tables    map[string]interface{}
	table_ops map[string]string
	m_cc      sync.Mutex
}

func New_clickhouse_client() *Clickhouse_client {
	return &Clickhouse_client{
		tables:    map[string]interface{}{},
		table_ops: map[string]string{},
	}
}

func (this *Clickhouse_client) Open(str_dsn string) error {
	if this.is_open {
		return nil
	}

	if str_dsn == "" {
		return nil
	}

	this.db_dsn = str_dsn

	var err error
	if this.db, err = gorm.Open(clickhouse.Open(this.db_dsn), &gorm.Config{}); err != nil {
		return err
	} else {
		if err := this.Health_check(0, 0, true); err != nil {
			return err
		}

		log.Println("clickhouse is online")
		go this.Health_check(10, 7, false)
	}

	for table_name, table_struct := range this.tables {

		log.Printf("auto migrate: %s", table_name)
		if table_op, ok := this.table_ops[table_name]; ok {
			if err := this.db.Set("gorm:table_options", table_op).AutoMigrate(table_struct); err != nil {
				return err
			}
		} else {
			if err := this.db.AutoMigrate(table_struct); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *Clickhouse_client) Health_check(delay, interval int, once bool) error {
	if interval < 0 {
		interval = 0
	}

	if delay < 0 {
		delay = 0
	}

	if delay != 0 {
		time.Sleep(time.Duration(delay) * time.Second)
	}

	raw_db, err := this.db.DB()
	if err != nil {
		return err
	}

	if once {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := raw_db.PingContext(ctx); err != nil {
			this.m_cc.Lock()
			this.is_open = false
			this.m_cc.Unlock()
			return err
		} else {
			this.m_cc.Lock()
			this.is_open = true
			this.m_cc.Unlock()
		}
	} else {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

			if err := raw_db.PingContext(ctx); err != nil {
				this.m_cc.Lock()
				this.is_open = false
				this.m_cc.Unlock()

				log.Println("db heartbeat: ", err)

				// just retry
				cancel()
				time.Sleep(3 * time.Second)
				continue
			} else {
				this.m_cc.Lock()
				this.is_open = true
				this.m_cc.Unlock()
			}

			cancel()
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}

	return nil
}

func (this *Clickhouse_client) Register(table_name, table_op string, table_struct interface{}) {
	this.m_cc.Lock()
	defer this.m_cc.Unlock()

	if this.is_open {
		return
	}

	this.table_ops[table_name] = table_op
	this.tables[table_name] = table_struct
}

func (this *Clickhouse_client) DB() (*gorm.DB, bool) {
	this.m_cc.Lock()
	defer this.m_cc.Unlock()

	if !this.is_open {
		return nil, false
	}
	return this.db, true
}

func Open_ck(str_dsn string) error {
	return default_ck_client.Open(str_dsn)
}

func Register_ck(table_name, table_op string, table_struct interface{}) {
	default_ck_client.Register(table_name, table_op, table_struct)
}

func Ck() (*gorm.DB, bool) {
	return default_ck_client.DB()
}
