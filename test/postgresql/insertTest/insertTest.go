// hello.go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	"time"
)

func main() {

	db, err := sql.Open("postgres",
		`host=192.168.209.128
			user=postgres
			password=Changeme_123
			dbname=testdb
			sslmode=disable`)

	if err != nil {
		fmt.Println(err)
	}

	age := 2
	r, err := db.Query("SELECT name,age FROM USERS WHERE age > $1", age)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Close()

	for r.Next() {
		var name string
		var age int
		_ = r.Scan(&name, &age)

		fmt.Printf("%s  %d\n", name, age)
	}

	fmt.Println(time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"))
	var a uint64
	a = 10000000000
	fmt.Println(strconv.FormatUint(a, 10))

	if 1 > 0 && 10 > 0 {
		fmt.Println("ok")
	}
}
