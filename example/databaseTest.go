package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"petadb"
	"time"
)

var database = petadb.NewDatabase("mysql", "mysql", "root:123456@/test?charset=utf8", false)

func main() {
	//insertTest()

	// FindOneTest()
	// FindOneTest2()
	// FindOneTest3()
	// FindOneTest4()

	// QueryTest1()
	// QuerySqlTest()

	PagedListTest()
	// UpdateTest()
}

func insertTest() {
	userInfo := UserInfo{UserName: "gejin", CreateDate: time.Now()}
	id, err := database.Insert(&userInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(userInfo)

	fmt.Println(id)
}

func FindOneTest() {
	userInfo := UserInfo{}
	isEixsts, err := database.FindOne(&userInfo, "SELECT * FROM UserInfo WHERE UserName = 'gejin'")
	if err != nil {
		panic(err)
	}

	fmt.Println(isEixsts)
	if isEixsts {
		fmt.Println(userInfo)
	}
}

func FindOneTest2() {
	var userId int
	isExists, err := database.FindOne(&userId, "SELECT UserId FROM UserInfo WHERE UserName = 'gejin'")
	if err != nil {
		panic(err)
	}

	fmt.Println(isExists)
	if isExists {
		fmt.Println(userId)
	}
}

func FindOneTest3() {
	var createDate time.Time
	isExists, err := database.FindOne(&createDate, "SELECT CreateDate FROM UserInfo WHERE UserName = 'gejin'")
	if err != nil {
		panic(err)
	}

	fmt.Println(isExists)
	if isExists {
		fmt.Println(createDate)
	}
}

func FindOneTest4() {
	sb := petadb.NewSqlBuilder()
	sb.Where("UserName = @0", "gejin")
	sb.Where("UserId = @0", 12)

	var userInfo UserInfo
	isExists, err := database.FindOneSql(&userInfo, &sb)
	if err != nil {
		panic(err)
	}

	fmt.Println(isExists)
	if isExists {
		fmt.Println(userInfo)
	}
}

func QueryTest1() {
	var userList []UserInfo
	if err := database.FindList(&userList, "SELECT * FROM UserInfo WHERE UserId > 1"); err != nil {
		panic(err)
	}

	fmt.Println(userList)
}

func FindListSqlTest() {

	sb := petadb.NewSqlBuilder()
	sb.Where("UserId > @0", 1)

	userList := make([]UserInfo, 0)
	if err := database.FindListSql(&userList, &sb); err != nil {
		panic(err)
	}

	fmt.Println(userList)
}

func PagedListTest() {
	var pagedInfo petadb.PagedInfo
	userList := make([]UserInfo, 0)

	if err := database.FindPagedList(&pagedInfo, &userList, 1, 10, "SELECT * FROM UserInfo"); err != nil {
		panic(err)
	}

	fmt.Println(pagedInfo)
	fmt.Println(userList)
}

func UpdateTest() {
	var userInfo UserInfo
	isEixsts, err := database.FindOne(&userInfo, "SELECT * FROM UserInfo WHERE UserName = 'gejin'")
	if err != nil {
		panic(err)
	}

	fmt.Println(isEixsts)
	if isEixsts {
		fmt.Println(userInfo)

		userInfo.CreateDate = time.Now()
		row, err := database.Update(&userInfo)
		if err != nil {
			panic(err)
		}

		fmt.Println(row)
		fmt.Println(userInfo)
	}
}

type UserInfo struct {
	UserId     int `petadb:"pk"`
	UserName   string
	CreateDate time.Time
}
