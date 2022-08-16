package yoursql

import (
	"database/sql/driver"
	"errors"
	"sync"
)

var (
	Err_faker_driver_bad_dsn = errors.New("faker_driver: bad dsn")
)

type Faker_driver struct {
	dbs    sync.Map // dsn : Faker_db
	expect *Faker_resps
}

func new_driver() *Faker_driver {
	return &Faker_driver{expect: new_resps()}
}

func (this *Faker_driver) Open(dsn string) (driver.Conn, error) {
	if len(dsn) == 0 {
		return nil, Err_faker_driver_bad_dsn
	}

	var ptr_db *Faker_db
	if ptr_raw_db, ok := this.dbs.Load(dsn); !ok {
		ptr_db = new_db(dsn, this)
		this.dbs.Store(dsn, ptr_db)
	} else {
		ptr_db, ok = ptr_raw_db.(*Faker_db)
		if !ok {
			ptr_db = new_db(dsn, this)
			this.dbs.Store(dsn, ptr_db)
		}
	}

	return new_conn(ptr_db), nil
}

// extensions
func (this *Faker_driver) Set_expect(rows_affected, last_insert_id int64, query string, rows []map[string]interface{}) error {
	return this.expect.attach([]*Faker_resp{
		{
			Pattern:        query,
			Resp:           rows,
			Rows_affected:  rows_affected,
			Last_insert_id: last_insert_id,
		},
	})
}

func (this *Faker_driver) Clean_expect(query string) {
	this.expect.detach(query)
}
