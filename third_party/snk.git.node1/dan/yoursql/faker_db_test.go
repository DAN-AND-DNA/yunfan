package yoursql

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	//	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

/*
func Test_Query(t *testing.T) {
	r := require.New(t)
	reply := []map[string]interface{}{
		{
			"name": "dan",
			"age":  "28",
		},
		{
			"name": "lulu",
			"age":  27,
		},
	}

	sql_key := "SELECT * FROM `users`"
	db, err := sql.Open(YourSql, "dsn")
	r.Nil(err)
	defer db.Close()

	Set_expect(2, 0, sql_key, reply)
	defer Clean_expect(sql_key)

	rows, err := db.Query("SELECT * FROM `users`")
	r.Nil(err)
	defer rows.Close()

	for rows.Next() {
		col_name, err := rows.Columns()
		r.Nil(err)
		_ = col_name
		var name string
		var age int64
		rows.Scan(&name, &age)
	}
}
func Test_GORM_find(t *testing.T) {
	r := require.New(t)
	type User struct {
		Name string `gorm:"column:name"`
		Age  uint64 `gorm:"column:age"`
	}
	version := []map[string]interface{}{
		{
			"VERSION()": "5.7.26",
		},
	}

	reply := []map[string]interface{}{
		{
			"name": "dan",
			"age":  28,
		},
		{
			"name": "yang",
			"age":  30,
		},
	}

	db, err := sql.Open(YourSql, "dsn")
	r.Nil(err)
	defer db.Close()

	sql_key1 := "SELECT VERSION()"
	sql_key2 := "SELECT * FROM `users`"

	Set_expect(1, 0, sql_key1, version)
	defer Clean_expect(sql_key1)
	Set_expect(2, 0, sql_key2, reply)
	defer Clean_expect(sql_key2)

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	r.Nil(err)
	var users []User

	gdb.Find(&users)
	r.Equal(users, []User{
		{
			Name: "dan",
			Age:  28,
		},
		{
			Name: "yang",
			Age:  30,
		},
	})
}

func Test_GORM_first(t *testing.T) {
	r := require.New(t)
	type User struct {
		Name string `gorm:"column:name"`
		Age  uint64 `gorm:"column:age"`
	}

	//========================================
	version := []map[string]interface{}{
		{
			"VERSION()": "5.7.26",
		},
	}

	reply := []map[string]interface{}{
		{
			"name": "dan",
			"age":  28,
		},
		{
			"name": "yang",
			"age":  30,
		},
	}
	//========================================
	db, err := sql.Open(YourSql, "dsn")
	r.Nil(err)
	defer db.Close()

	sql_key1 := "SELECT VERSION()"
	Set_expect(1, 0, sql_key1, version)
	defer Clean_expect(sql_key1)

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	r.Nil(err)

	sql_key2 := "SELECT * FROM `users` ORDER BY `users`.`name` LIMIT 1"
	Set_expect(2, 0, sql_key2, reply)
	defer Clean_expect(sql_key2)

	var user User
	gdb.First(&user)
	r.Equal(user, User{
		Name: "dan",
		Age:  28,
	})

}

func Test_GORM_first_where(t *testing.T) {
	r := require.New(t)
	type User struct {
		Name string `gorm:"column:name"`
		Age  uint64 `gorm:"column:age"`
	}

	//========================================
	version := []map[string]interface{}{
		{
			"VERSION()": "5.7.26",
		},
	}

	reply := []map[string]interface{}{
		{
			"name": "dan",
			"age":  28,
		},
		{
			"name": "yang",
			"age":  30,
		},
	}
	//========================================

	db, err := sql.Open(YourSql, "dsn")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sql_key1 := "SELECT VERSION()"
	Set_expect(1, 0, sql_key1, version)
	defer Clean_expect(sql_key1)

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	sql_key2 := "SELECT * FROM `users` WHERE age = 28 ORDER BY `users`.`name` LIMIT 1"
	Set_expect(2, 0, sql_key2, reply)
	defer Clean_expect(sql_key2)

	var user User
	gdb.Where("age = ?", 28).First(&user)
	r.Equal(user, User{
		Name: "dan",
		Age:  28,
	})
}
*/
func Test_GORM_first_where_postgres(t *testing.T) {
	r := require.New(t)
	type User struct {
		Name string `gorm:"column:name"`
		Age  uint64 `gorm:"column:age"`
	}

	//========================================
	version := []map[string]interface{}{
		{
			"VERSION()": "5.7.26",
		},
	}

	reply := []map[string]interface{}{
		{
			"name": "dan",
			"age":  28,
		},
		{
			"name": "yang",
			"age":  30,
		},
	}
	//========================================

	db, err := sql.Open(YourSql, "dsn")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sql_key1 := "SELECT VERSION()"
	Set_expect(1, 0, sql_key1, version)
	defer Clean_expect(sql_key1)

	fdb, err := sql.Open(YourSql, "dsn")
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fdb,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	sql_key2 := `SELECT * FROM "users" WHERE age = 28 OR age = 35 ORDER BY "users"."name" LIMIT 1`
	Set_expect(2, 0, sql_key2, reply)
	defer Clean_expect(sql_key2)

	var user User
	gdb.Where("age = ? OR age = ?", 28, 35).First(&user)
	r.Equal(user, User{
		Name: "dan",
		Age:  28,
	})

}
