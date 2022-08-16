package yoursql

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ERR_STMT_CLOSED      = errors.New("stmt has been closed")
	ERR_STMT_BAD_COMMAND = errors.New("stmt bad command")
)

type Faker_stmt struct {
	owner_conn      *Faker_conn
	next            *Faker_stmt
	expect_query    string
	is_closed       bool
	command         string
	base_id         int64
	placeholders    int
	str_placeholder string
}

// 1. Smt
func (this *Faker_stmt) Close() error {
	if this.owner_conn == nil {
		panic("no conn but stmt is bound to a conn")
	}

	if this.owner_conn.owner_db == nil {
		panic("no db but conn is bound to a database")
	}

	this.is_closed = true
	if this.next != nil {
		this.next.Close()
	}

	return nil
}

func (this *Faker_stmt) NumInput() int {
	return this.placeholders
}

func (this *Faker_stmt) Exec(args []driver.Value) (driver.Result, error) {
	panic("deprecated: use stmtExecContext instead")
}

func (this *Faker_stmt) Query(args []driver.Value) (driver.Rows, error) {
	panic("deprecated: use stmtQueryContext instead")
}

// 2. ColumnConverter
func (this *Faker_stmt) ColumnConverter(idx int) driver.ValueConverter {
	return driver.DefaultParameterConverter
}

// 3. StmtExecContext
func (this *Faker_stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if this.is_closed {
		return nil, ERR_STMT_CLOSED
	}

	for index, arg := range args {
		if this.str_placeholder == "?" {
			this.expect_query = strings.Replace(this.expect_query, "?", "%v", 1)
		} else {
			this.expect_query = strings.Replace(this.expect_query, "$"+strconv.Itoa(index+1), "%v", 1)
		}
		this.expect_query = fmt.Sprintf(this.expect_query, arg.Value)
	}

	resp := this.owner_conn.owner_db.owner_driver.expect.find_resp(this.expect_query)
	if resp.Expect_exec_bad_conn != nil {
		if resp.Expect_exec_bad_conn() == true {
			return nil, driver.ErrBadConn
		}
	}

	switch this.command {
	case "INSERT":
		ptr_result := &Faker_result{}
		ptr_result.Set(resp.Rows_affected, resp.Last_insert_id)
		return ptr_result, nil
	case "UPDATE":
		return driver.RowsAffected(resp.Rows_affected), nil
	case "DELETE":
		return driver.RowsAffected(resp.Rows_affected), nil
	default:
		return nil, ERR_STMT_BAD_COMMAND
	}

}

// 4. StmtQueryContext
func (this *Faker_stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {

	if this.is_closed {
		return nil, ERR_STMT_CLOSED
	}

	for index, arg := range args {
		if this.str_placeholder == "?" {
			this.expect_query = strings.Replace(this.expect_query, "?", "%v", 1)
		} else {
			this.expect_query = strings.Replace(this.expect_query, "$"+strconv.Itoa(index+1), "%v", 1)
		}
		this.expect_query = fmt.Sprintf(this.expect_query, arg.Value)
	}

	resp := this.owner_conn.owner_db.owner_driver.expect.find_resp(this.expect_query)
	if resp.Expect_query_bad_conn != nil {
		if resp.Expect_query_bad_conn() == true {
			return nil, driver.ErrBadConn
		}
	}

	if resp.Expect_error != nil {
		return nil, resp.Expect_error
	}

	//clone

	resp.result_rows.reuse()
	new_result_rows := &Faker_rows{
		is_closed:              resp.result_rows.is_closed,
		current_row_index:      resp.result_rows.current_row_index,
		result_set_index:       resp.result_rows.result_set_index,
		expect_error_row_index: resp.result_rows.expect_error_row_index,
		expect_error:           resp.result_rows.expect_error,
	}

	get_cols := false
	tmp := make([]*real_row, 0, len(resp.Resp))
	cols_index := map[string]int{}
	i := 0
	for _, vals := range resp.Resp {
		tmp_row := &real_row{
			cols: make([]interface{}, len(vals)),
		}
		for col_name, val := range vals {

			if get_cols == false {
				new_result_rows.cols = append(new_result_rows.cols, col_name)
				cols_index[col_name] = i
				i++
			}

			tmp_row.cols[cols_index[col_name]] = val
		}

		get_cols = true
		tmp = append(tmp, tmp_row)
	}

	new_result_rows.rows = append(new_result_rows.rows, tmp)
	return new_result_rows, nil
}
