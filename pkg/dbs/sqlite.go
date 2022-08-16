package dbs

import (
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	default_sqlite_client = New_sqlite_client()
)

type Sqlite_client struct {
	sync.RWMutex
	is_open   bool
	is_mock   bool
	db        *gorm.DB
	db_dsn    string
	tables    map[string]interface{}
	table_ops map[string]string
}

func New_sqlite_client() *Sqlite_client {
	return &Sqlite_client{
		tables:    map[string]interface{}{},
		table_ops: map[string]string{},
	}
}

func (this *Sqlite_client) Open(str_dsn string, is_mock bool) error {
	this.Lock()
	defer this.Unlock()

	if this.is_open {
		return nil
	}

	if str_dsn == "" {
		return nil
	}

	this.db_dsn = str_dsn
	var err error

	if is_mock {
		// do mock
		// just in memory
		if this.db, err = gorm.Open(sqlite.Open("file:mockdb?mode=memory&cache=shared&_auto_vacuum=none"), &gorm.Config{}); err != nil {
			return err
		}

		this.is_mock = true
		log.Println("fake sqlite: ok")
		for table_name, table_struct := range this.tables {

			log.Printf("auto migrate: %s", table_name)
			if table_op, ok := this.table_ops[table_name]; ok {
				if err = this.db.Set("gorm:table_options", table_op).AutoMigrate(table_struct); err != nil {
					return err
				}
			} else {
				if err = this.db.AutoMigrate(table_struct); err != nil {
					return err
				}
			}
		}

	} else {
		if this.db, err = gorm.Open(sqlite.Open(this.db_dsn), &gorm.Config{}); err != nil {
			return err
		}

		log.Println("sqlite: ok")
		for table_name, table_struct := range this.tables {

			log.Printf("auto migrate: %s", table_name)
			if table_op, ok := this.table_ops[table_name]; ok {
				if err = this.db.Set("gorm:table_options", table_op).AutoMigrate(table_struct); err != nil {
					return err
				}
			} else {
				if err = this.db.AutoMigrate(table_struct); err != nil {
					return err
				}
			}
		}
	}

	this.is_open = true
	return nil
}

func (this *Sqlite_client) Shutdown() {
	this.Lock()
	defer this.Unlock()

	if this.is_open {
		raw_db, err := this.db.DB()
		if err == nil {
			raw_db.Close()
		}
	}

	this.is_open = false
}

func (this *Sqlite_client) Reopen() error {
	this.Lock()
	defer this.Unlock()

	if !this.is_open {
		return nil
	}

	raw_db, err := this.db.DB()
	if err == nil {
		raw_db.Close()
	}

	this.is_open = false

	if this.is_mock {
		if this.db, err = gorm.Open(sqlite.Open("file:mockdb?mode=memory&cache=shared&_auto_vacuum=none"), &gorm.Config{}); err != nil {
			return err
		}
		log.Println("fake sqlite: reopen ok")
		for table_name, table_struct := range this.tables {

			log.Printf("auto migrate: %s", table_name)
			if table_op, ok := this.table_ops[table_name]; ok {
				if err = this.db.Set("gorm:table_options", table_op).AutoMigrate(table_struct); err != nil {
					return err
				}
			} else {
				if err = this.db.AutoMigrate(table_struct); err != nil {
					return err
				}
			}
		}

	} else {

		if this.db, err = gorm.Open(sqlite.Open(this.db_dsn), &gorm.Config{}); err != nil {
			return err
		}

		log.Println("sqlite: reopen ok")
		for table_name, table_struct := range this.tables {

			log.Printf("auto migrate: %s", table_name)
			if table_op, ok := this.table_ops[table_name]; ok {
				if err = this.db.Set("gorm:table_options", table_op).AutoMigrate(table_struct); err != nil {
					return err
				}
			} else {
				if err = this.db.AutoMigrate(table_struct); err != nil {
					return err
				}
			}
		}
	}

	this.is_open = true
	return nil

}

func (this *Sqlite_client) Register(table_name, table_op string, table_struct interface{}) {
	this.RLock()
	defer this.RUnlock()

	if this.is_open {
		return
	}

	this.tables[table_name] = table_struct
	this.table_ops[table_name] = table_op
}

func (this *Sqlite_client) Begin() (*gorm.DB, bool) {
	this.RLock()

	if !this.is_open {
		this.RUnlock()
		return nil, false
	}

	return this.db, true
}

func (this *Sqlite_client) End() {
	this.RUnlock()
}

func Open_sqlite(str_dsn string, is_test bool) error {
	return default_sqlite_client.Open(str_dsn, is_test)
}

func Register_sqlite(table_name, table_op string, table_struct interface{}) {
	default_sqlite_client.Register(table_name, table_op, table_struct)
}

func Reopen_sqlite() error {
	return default_sqlite_client.Reopen()
}

func Begin_sqlite() (*gorm.DB, bool) {
	return default_sqlite_client.Begin()
}

func End_sqlite() {
	default_sqlite_client.End()
}
