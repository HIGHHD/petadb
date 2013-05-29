package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"petadb"
	"time"
)

var database = petadb.NewDatabase("mysql", "mysql", "root:123456@/test?charset=utf8")

func main() {
	database.IsDebug = true
	// InsertTest()
	// QueryStructTest()
	// QuerySbTest()
	// QueryPrtTest()
	// FirstTest()
	DeleteTest()
}

func InsertTest() {
	userInfo := UserInfo{UserName: "Deng", CreateDate: time.Now().String()}
	id, err := database.Insert(&userInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println(id)
	fmt.Println(userInfo.UserId)
}

func QueryStructTest() {
	userList := make([]UserInfo, 0)
	if err := database.Query(&userList, "SELECT * FROM UserInfo"); err != nil {
		panic(err)
	}
	fmt.Println(userList)
}

func QuerySbTest() {
	userList := make([]UserInfo, 0)

	sb := petadb.NewSqlBuilder()
	sb.Where("UserId = @0", 1)
	sb.Where("UserName = @0", "Deng")

	if err := database.QueryBySb(&userList, &sb); err != nil {
		panic(err)
	}
	fmt.Println(userList)
}

func QueryPrtTest() {
	userNameList := make([]string, 0)
	if err := database.Query(&userNameList, "select UserName from UserInfo"); err != nil {
		panic(err)
	}

	fmt.Println(userNameList)
}

func FirstTest() {
	userInfo := UserInfo{}
	isExists, err := database.First(&userInfo, "select * from userInfo where userId = @0", 1)
	if err != nil {
		panic(err)
	}

	fmt.Println(isExists)
	fmt.Println(userInfo)
}

func DeleteTest() {
	userInfo := UserInfo{}

	sb := petadb.NewSqlBuilder()
	sb.Where("UserId = @0", 2)

	isExists, err := database.FirstBySb(&userInfo, &sb)
	if err != nil {
		panic(err)
	}

	if isExists {
		id, err := database.Delete(&userInfo)
		if err != nil {
			panic(err)
		}

		fmt.Println(id)
	}
}

type UserInfo struct {
	UserId     int `pk`
	UserName   string
	CreateDate string
}
