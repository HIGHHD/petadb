package main

import (
	"fmt"
	"petadb"
	"time"
)

func main() {
	Test1()
	Test2()
}

type Human struct {
	Id   int
	Name string
	Age  int
}

type UserInfo struct {
	Guid       string `pk_notauto` // 主键，非自增
	UserName   string
	CreateDate time.Time `notmap` // 不映射到数据表
}

func Test1() {
	tableInfo, err := petadb.GetTableInfo(Human{})
	if err != nil {
		panic(err)
	}

	fmt.Println(tableInfo)
}

func Test2() {
	tableInfo, err := petadb.GetTableInfo(UserInfo{})
	if err != nil {
		panic(err)
	}

	fmt.Println(tableInfo)
}
