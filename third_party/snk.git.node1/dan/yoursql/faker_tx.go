package yoursql

import (
	"database/sql/driver"
)

var (
	Trigger_bad_conn_commit   func() bool
	Trigger_bad_conn_rollback func() bool
)

type Faker_tx struct {
	owner_conn *Faker_conn
}

func new_tx(ptr_conn *Faker_conn) *Faker_tx {
	return &Faker_tx{ptr_conn}
}

func (this *Faker_tx) Commit() error {
	this.owner_conn.current_tx = nil

	if Trigger_bad_conn_commit != nil {
		if Trigger_bad_conn_commit() == true {
			return driver.ErrBadConn
		}
	}

	return nil
}

func (this *Faker_tx) Rollback() error {
	this.owner_conn.current_tx = nil

	if Trigger_bad_conn_rollback != nil {
		if Trigger_bad_conn_rollback() == true {
			return driver.ErrBadConn
		}
	}

	return nil
}
