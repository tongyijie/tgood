// hello.go
package main

import (
	"bytes"
	"database/sql"
	"errors"
	log "github.com/alecthomas/log4go"
	_ "github.com/lib/pq"
	"math/rand"
	"strconv"
	"time"
)

var (
	db *sql.DB
)

func initDB() error {
	db, err := sql.Open("postgres",
		`host=192.168.209.128
			user=postgres
			password=Changeme_123
			dbname=testdb
			sslmode=disable`)

	if err != nil {
		log.Error("init db error. error=", err)
		db = nil

		return errors.New("init db error.")
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(4)

	return nil
}

func main() {
	log.LoadConfiguration("log4go.xml")
	defer log.Close()

	log.Info("test begin")

	//抓取异常
	defer func() {
		if r := recover(); r != nil {
			log.Error("Recovered in f", r)
		}
	}()

	//初始化数据库
	err := initDB()
	if err != nil {
		log.Error("init db error. program will exit.")
		return
	}
	defer db.Close()

	total := 10000
	index := 1
	sqlBuf := bytes.NewBufferString("")
	manageId := uint64(10000000000001706277)
	dataId := uint64(100)
	rand.Seed(time.Now().Unix())
	start := time.Now()

	for i := 0; i < total; i++ {
		tx, _ := db.Begin()

		dataId = uint64(100)
		for j := 0; i < 100 && index <= total; j++ {
			sqlBuf.Truncate(0)

			sqlBuf.WriteString("insert into table (manage_id,data_id,val,record_time) values ")
			sqlBuf.WriteString("(")
			sqlBuf.WriteString(strconv.FormatUint(manageId, 10))
			sqlBuf.WriteString(strconv.FormatUint(dataId, 10))
			sqlBuf.WriteString(strconv.Itoa(rand.Intn(200)))
			sqlBuf.WriteString("'" + time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05") + "'")
			sqlBuf.WriteString(")")

			stmt, err := tx.Prepare(sqlBuf.String())
			if err != nil {
				log.Error("prepare sql statment error. error=%s \nsql=%s", err, sqlBuf.String())
				tx.Rollback()
				panic(err)
			}

			_, err = stmt.Exec()
			if err != nil {
				log.Error("exec sql statment error. error=%s \nsql=%s", err, sqlBuf.String())
				tx.Rollback()
				panic(err)
			}

			tx.Commit()

		}
		manageId += 1
	}

	end := time.Now()

	log.Info("test end. insert %d datas, Elapsed %d second.", total, end.Sub(start).Seconds())
}
