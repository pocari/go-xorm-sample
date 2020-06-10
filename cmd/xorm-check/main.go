package main

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

// Railsと違ってテーブル名そのままのstruct名になる
type Hoges struct {
	Id   int64  `xorm:"'id'"`
	Name string `xorm:"'name'"`
}

func mysqlConnectionString(database, user, pass, host string, port int) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		user,
		pass,
		host,
		port,
		database,
	)
}

func createHoge(engine *xorm.Engine, name string) error {
	has, err := engine.Exist(&Hoges{
		Name: name,
	})
	if err != nil {
		return err
	}
	if !has {
		engine.Table("hoges").Insert(Hoges{
			Name: name,
		})
	}
	return nil
}

func InsertHoges(engine *xorm.Engine) {
	session := engine.NewSession()
	defer session.Clone()

	if err := createHoge(engine, "hoge_name01"); err != nil {
		panic(err)
	}
	if err := createHoge(engine, "hoge_name02"); err != nil {
		panic(err)
	}
	if err := createHoge(engine, "hoge_name03"); err != nil {
		panic(err)
	}
}

func main() {
	dbPort, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		panic(err)
	}
	engine, err := xorm.NewEngine("mysql", mysqlConnectionString(
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		dbPort,
	))
	if err != nil {
		panic(err)
	}

	defer engine.Close()

	InsertHoges(engine)

	var hoges []Hoges
	if err := engine.Find(&hoges); err != nil {
		panic(err)
	}

	for i, h := range hoges {
		fmt.Printf("%02d: hoge: %+v\n", i, h)
	}
}
