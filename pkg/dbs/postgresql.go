package dbs

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"snk.git.node1/dan/yoursql"
)

var (
	default_postgresql_client = New_postgresql_client()
)

type Postgresql_client struct {
	is_open   bool
	db        *gorm.DB
	db_dsn    string
	m_cc      sync.RWMutex
	tables    map[string]interface{}
	table_ops map[string]string
}

func New_postgresql_client() *Postgresql_client {
	return &Postgresql_client{
		tables:    map[string]interface{}{},
		table_ops: map[string]string{},
	}
}

func (this *Postgresql_client) Open(str_dsn string, is_test bool) error {
	if this.is_open {
		return nil
	}

	if str_dsn == "" {
		return nil
	}

	this.db_dsn = str_dsn
	var err error

	if is_test {
		// do mock
		fdb, err := sql.Open(yoursql.YourSql, "dsn")
		if this.db, err = gorm.Open(postgres.New(postgres.Config{Conn: fdb}), &gorm.Config{}); err != nil {
			return err
		} else {
			this.is_open = true
			log.Println("faker postgresql: ok")
		}

	} else {

		if this.db, err = gorm.Open(postgres.Open(this.db_dsn), &gorm.Config{}); err != nil {
			return err
		} else {
			if err := this.Health_check(0, 0, true); err != nil {
				return err
			}

			log.Println("postgresql: ok")
			go this.Health_check(15, 10, false)
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
	}

	return nil
}

func (this *Postgresql_client) Health_check(delay, interval int, once bool) error {
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

func (this *Postgresql_client) Register(table_name, table_op string, table_struct interface{}) {
	this.m_cc.RLock()
	defer this.m_cc.RUnlock()

	if this.is_open {
		return
	}

	this.tables[table_name] = table_struct
	this.table_ops[table_name] = table_op
}

func (this *Postgresql_client) DB() (*gorm.DB, bool) {
	this.m_cc.Lock()
	defer this.m_cc.Unlock()

	if !this.is_open {
		return nil, false
	}
	return this.db, true
}

func (this *Postgresql_client) Begin() (*gorm.DB, bool) {
	this.m_cc.RLock()

	if !this.is_open {
		this.m_cc.RUnlock()
		return nil, false
	}

	return this.db, true
}

func (this *Postgresql_client) End() {
	this.m_cc.RUnlock()
}

func Open_postgres(str_dsn string, is_test bool) error {
	return default_postgresql_client.Open(str_dsn, is_test)
}

func Register_postgres(table_name, table_op string, table_struct interface{}) {
	default_postgresql_client.Register(table_name, table_op, table_struct)
}

func Postgres() (*gorm.DB, bool) {
	return default_postgresql_client.DB()
}

func Begin_postgres() (*gorm.DB, bool) {
	return default_postgresql_client.Begin()
}

func End_postgres() {
	default_postgresql_client.End()
}
