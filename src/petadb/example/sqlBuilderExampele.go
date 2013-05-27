package main

import (
	"fmt"
	"petadb"
)

func main() {
	// AppendTest()
	// SelectTest()
	// WhereTest()
	GroupByHavingTest()
}

func AppendTest() {
	sb := petadb.NewSqlBuilder()
	sb.Append("SELECT * FROM UserInfo WHERE UserId = @0", 1)
	fmt.Println(sb)
}

func SelectTest() {
	sb := petadb.NewSqlBuilder()
	sb.Select("*").From("UserInfo").Where("UserId = @0 AND UserState = @1", 1, 2)
	fmt.Println(sb)
}

func WhereTest() {
	sb := petadb.NewSqlBuilder()
	sb.Select("UserID").From("UserInfo")
	sb.Where("UserId = @0", 1)
	sb.Where("UserState = @0", 2)
	fmt.Println(sb)
}

func GroupByHavingTest() {
	sb := petadb.NewSqlBuilder()
	sb.Select("UserID,UserName").From("UserInfo")
	sb.Where("UserId = @0", 1)
	sb.Where("UserState = @0", 2)
	sb.GroupBy("UserId")
	sb.GroupBy("UserName").Having("COUNT(1) > @0", 1)
	fmt.Println(sb)
}
