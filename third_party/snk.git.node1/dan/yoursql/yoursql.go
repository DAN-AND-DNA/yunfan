package yoursql

import (
	"database/sql"
)

var (
	YourSql        = "yoursql_driver"
	default_driver = new_driver()
)

func init() {
	drivers := sql.Drivers()
	for _, name := range drivers {
		if name == YourSql {
			return
		}
	}

	sql.Register(YourSql, default_driver)
}

// extensions
func Set_expect(rows_affected, last_insert_id int64, query string, rows []map[string]interface{}) error {
	return default_driver.Set_expect(rows_affected, last_insert_id, query, rows)
}

func Clean_expect(query string) {
	default_driver.Clean_expect(query)
}
