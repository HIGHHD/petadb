package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"petadb"
	"time"
)

var database = petadb.NewDatabase("mysql", "mysql", "root:123456@/test?charset=utf8")

func main() {
	//insertTest()

	// FindOneTest()
	// FindOneTest2()
	// FindOneTest3()
	// FindOneTest4()

	// FindListTest()
	FindListSqlTest()

	// PagedListTest()
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
	// 是否存在
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

func FindListSqlTest() {
	sb := petadb.NewSqlBuilder()

	sb.Where("UserId > @0", 1)
	sb.Where("CreateDate < @0", time.Now())
	// 当有多个查询条件时，Sql组装是一件非常痛苦的事情
	// sqlBuilder应运而生，多个查询下，简单的sql.Where即可完全应对

	var userList []UserInfo
	// 什么？ 不需要写select ... from ??
	// 对，不需要，组件会自动帮你匹配
	err := database.FindListSql(&userList, &sb)
	if err != nil {
		panic(err)
	}

	fmt.Println(userList)
}

func QueryTest1() {
	var userList []UserInfo
	if err := database.FindList(&userList, "SELECT * FROM UserInfo WHERE UserId > 1"); err != nil {
		panic(err)
	}

	fmt.Println(userList)
}

func FindListTest() {
	var userList []UserInfo
	if err := database.FindList(&userList, "SELECT * FROM UserInfo"); err != nil {
		panic(err)
	}

	fmt.Println(userList)
}

func PagedListTest() {
	var pagedInfo petadb.PagedInfo
	userList := make([]UserInfo, 0)

	// SQL语句会自动转换为分页语句(1.查询总数语句,2.查询列表语句)
	if err := database.FindPagedList(&pagedInfo, &userList, 1, 10, "SELECT * FROM UserInfo"); err != nil {
		panic(err)
	}

	fmt.Println(pagedInfo)
	fmt.Println(userList)
}

func UpdateTest() {
	var userInfo UserInfo
	// 取出要修改的对象
	isEixsts, err := database.FindOne(&userInfo, "SELECT * FROM UserInfo WHERE UserName = 'gejin'")
	if err != nil {
		panic(err)
	}

	fmt.Println(isEixsts)
	// 如果存在，则Update
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

func DeleteTest() {
	var userInfo UserInfo
	// 取出要修改的对象
	isEixsts, err := database.FindOne(&userInfo, "SELECT * FROM UserInfo WHERE UserName = 'gejin'")
	if err != nil {
		panic(err)
	}

	fmt.Println(isEixsts)
	// 如果存在，则Delete
	if isEixsts {
		fmt.Println(userInfo)

		row, err := database.Delete(&userInfo)
		if err != nil {
			panic(err)
		}

		fmt.Println(row)
	}
}

type UserInfo struct {
	UserId     int `petadb:"pk"`
	UserName   string
	CreateDate time.Time
}
