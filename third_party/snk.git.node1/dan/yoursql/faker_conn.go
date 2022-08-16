package yoursql

import (
	"context"
	"database/sql/driver"
	"errors"
	"strings"
)

var (
	ERR_ALREADY_IN_A_TX = errors.New("error: already in a transaction")
	ERR_BAD_QUERY       = errors.New("error: bad query")
)

type Faker_conn struct {
	current_tx *Faker_tx
	owner_db   *Faker_db
}

func new_conn(ptr_db *Faker_db) *Faker_conn {
	return &Faker_conn{current_tx: nil, owner_db: ptr_db}
}

// 1. Conn
func (this *Faker_conn) Begin() (driver.Tx, error) {
	if this.current_tx != nil {
		return nil, ERR_ALREADY_IN_A_TX
	}

	this.current_tx = new_tx(this)
	return this.current_tx, nil
}

func (this *Faker_conn) Close() error {
	if this.current_tx != nil {
		this.current_tx.Rollback()
	}

	this.owner_db = nil
	return nil
}

func (this *Faker_conn) Prepare(query string) (driver.Stmt, error) {
	panic("use PrepareContext instead")
}

// 2. Execer
func (this *Faker_conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	panic("deprecated! use ExecContext instead")
}

// 3. ExecerContext
func (this *Faker_conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return nil, driver.ErrSkip
}

// 4. Queryer
func (this *Faker_conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	panic("deprecated! use QueryContext instead")
}

// 5. QueryerContext
func (this *Faker_conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return nil, driver.ErrSkip
}

// 6. ConnPrepareContext
func (this *Faker_conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if len(query) == 0 {
		return nil, ERR_BAD_QUERY
	}

	var ptr_first = &Faker_stmt{owner_conn: this, expect_query: query}

	ptr_first.placeholders = strings.Count(query, "?")
	if ptr_first.placeholders != 0 {
		ptr_first.str_placeholder = "?"
	} else {
		ptr_first.placeholders = strings.Count(query, "$")
		ptr_first.str_placeholder = "$"
	}

	query_lists := strings.Split(query, " ")
	if len(query_lists) < 2 {
		return nil, ERR_BAD_QUERY
	}

	ptr_first.command = strings.ToUpper(query_lists[0])
	return ptr_first, nil
}
