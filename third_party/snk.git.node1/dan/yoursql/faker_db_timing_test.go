package yoursql

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func Benchmark_Query(b *testing.B) {
	b.ReportAllocs()

	version := []map[string]interface{}{
		{
			"VERSION()": "5.7.26",
		},
	}

	reply := []map[string]interface{}{
		{
			"age": 28,
		},
		{
			"age": 28,
		},
	}

	var err error
	db, err := sql.Open(YourSql, "dsn")
	if err != nil {
		panic(err)
	}

	Set_expect(1, 0, "SELECT VERSION()", version)
	Set_expect(2, 0, "SELECT * FROM `users` ORDER BY `users`.`age` LIMIT 1", reply)

	type User struct {
		Age uint64 `gorm:"column:age"`
	}

	b.RunParallel(func(pb *testing.PB) {
		var user User
		for pb.Next() {

			rows, err := db.Query("SELECT * FROM `users` ORDER BY `users`.`age` LIMIT 1")
			if err != nil {
				panic(err)
			}

			for rows.Next() {
				if err := rows.Scan(&(user.Age)); err != nil {
					panic(err)
				}

				if user.Age == 0 {
					panic("age not 0")
				}
				break
			}
			rows.Close()
		}
	})

}
func Benchmark_GORM_first(b *testing.B) {
	b.ReportAllocs()

	version := []map[string]interface{}{
		{
			"VERSION()": "5.7.26",
		},
	}

	reply := []map[string]interface{}{
		{
			"age": 28,
		},
		{
			"age": 28,
		},
	}

	var err error
	db, err := sql.Open(YourSql, "dsn")
	if err != nil {
		panic(err)
	}

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	Set_expect(1, 0, "SELECT VERSION()", version)
	Set_expect(2, 0, "SELECT * FROM `users` ORDER BY `users`.`age` LIMIT 1", reply)

	type User struct {
		Age uint64 `gorm:"column:age"`
	}

	b.RunParallel(func(pb *testing.PB) {

		var user User
		for pb.Next() {
			err = gdb.First(&user).Error
			if err != nil {
				panic(err)
			}

			if user.Age == 0 {
				panic("age not 0")
			}
		}
	})

}
