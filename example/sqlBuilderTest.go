package main

import (
	"fmt"
	"petadb"
	"time"
)

func main() {
	// AppendTest()
	// SelectTest()
	// FromTest()
	// WhereTest()
	// OrderByTest()
	// GroupByTest()
	// HavingTest()
	JoinTest()
}

func AppendTest() {
	sqlBuilder := petadb.NewSqlBuilder()
	sqlBuilder.Append("select * from userInfo where UserId = @0 and UserName = @1", 1, 2)
	fmt.Println(sqlBuilder.SQL)
	fmt.Println(sqlBuilder.Args)
}

func SelectTest() {
	sqlBuilder := petadb.NewSqlBuilder()
	sqlBuilder.Select("UserId,UserName")
	fmt.Println(sqlBuilder.SQL)
}

func FromTest() {
	sqlBuilder := petadb.NewSqlBuilder()
	sqlBuilder.Select("UserId,UserName").From("UserInfo")
	fmt.Println(sqlBuilder.SQL)
}

func WhereTest() {
	sqlBuilder := petadb.NewSqlBuilder()
	sqlBuilder.Select("UserId,UserName").From("UserInfo")
	sqlBuilder.Where("UserName = @0", "gejinlove")
	sqlBuilder.Where("CreateDate >= @0", time.Now())

	fmt.Println(sqlBuilder.SQL)
	fmt.Println(sqlBuilder.Args)
}

func OrderByTest() {
	sqlBuilder := petadb.NewSqlBuilder()
	sqlBuilder.Select("UserId,UserName").From("UserInfo").OrderBy("CreateDate DESC")
	sqlBuilder.OrderBy("UserId ASC,UserName DESC")
	fmt.Println(sqlBuilder.SQL)
}

func GroupByTest() {
	sqlBuilder := petadb.NewSqlBuilder()
	sqlBuilder.Select("UserId,UserName").From("UserInfo")
	sqlBuilder.GroupBy("UserId")
	sqlBuilder.GroupBy("UserName")
	fmt.Println(sqlBuilder.SQL)
}

func HavingTest() {
	sqlBuilder := petadb.NewSqlBuilder()
	sqlBuilder.Select("UserId,UserName").From("UserInfo")
	sqlBuilder.GroupBy("UserId")
	sqlBuilder.GroupBy("UserName")
	sqlBuilder.Having("Count(1) > @0", 1)
	fmt.Println(sqlBuilder.SQL)
	fmt.Println(sqlBuilder.Args)
}

func JoinTest() {
	sqlBuilder := petadb.NewSqlBuilder()
	sqlBuilder.Select("a.UserId,a.UserName").From("UserInfo a")
	sqlBuilder.InnerJoin("UserRole b").On("a.RoleId = b.RoleId AND a.UserId = @0", 1)
	sqlBuilder.LeftJoin("UserPassword c").On("a.UserId = c.UserId")
	fmt.Println(sqlBuilder.SQL)
	fmt.Println(sqlBuilder.Args)
}
