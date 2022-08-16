package yoursql

import (
	//"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	//	"reflect"
)

var (
	ERR_ROWS_CLOSED = errors.New("error: rows has been closed")
)

type Faker_rows struct {
	cols                   []string
	is_closed              bool
	current_row_index      int
	result_set_index       int
	expect_error_row_index int
	expect_error           error
	rows                   [][]*real_row
	//bytes_clone            map[*byte][]byte
}

type real_row struct {
	cols []interface{}
}

func (this *Faker_rows) reuse() {
	this.is_closed = false
	this.current_row_index = -1
}

// 1. Rows
func (this *Faker_rows) Columns() []string {
	return this.cols
}

func (this *Faker_rows) Close() error {
	if this.is_closed == true {
		return nil
	}
	this.is_closed = true
	return nil
}

func (this *Faker_rows) Next(dest []driver.Value) error {
	if this.is_closed == true {
		return ERR_ROWS_CLOSED
	}

	has_cols := false
	for _, sub_rows := range this.rows {
		for _, row := range sub_rows {
			if len(row.cols) != 0 {
				has_cols = true
			}
		}
	}

	if !has_cols {
		return io.EOF
	}

	this.current_row_index++
	if this.current_row_index == this.expect_error_row_index {
		return this.expect_error
	}

	if this.current_row_index >= len(this.rows[this.result_set_index]) {
		return io.EOF
	}

	for i, val := range this.rows[this.result_set_index][this.current_row_index].cols {
		dest[i] = val
	}

	return nil
}

// 2. RowsNextResultSet
func (this *Faker_rows) HasNextResultSet() bool {
	return this.result_set_index < len(this.rows)-1
}

func (this *Faker_rows) NextResultSet() error {
	if this.HasNextResultSet() == true {
		this.result_set_index++
		this.current_row_index = -1
		return nil
	}
	return io.EOF
}
