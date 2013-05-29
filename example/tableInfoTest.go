package main

import (
	"fmt"
	"petadb"
	"time"
)

func main() {
	MapIntoObjectTest()
	MapIntoObjectTest2()
	MapIntoObjectTest3()
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

func MapIntoObjectTest() {
	objMap := make(map[string][]byte, 0)
	objMap["Id"] = []byte("1")
	objMap["Name"] = []byte("gejin")
	objMap["Age"] = []byte("29")

	human := Human{}
	if err := petadb.MapIntObject(&human, objMap); err != nil {
		panic(err)
	}

	fmt.Println(human)
}

func MapIntoObjectTest2() {
	objMap := make(map[string][]byte, 0)
	objMap["UserId"] = []byte("1")
	var i int
	if err := petadb.MapIntObject(&i, objMap); err != nil {
		panic(err)
	}
	fmt.Println(i)
}

func MapIntoObjectTest3() {
	objMap := make(map[string][]byte, 0)
	objMap["time"] = []byte("2013-05-29 15:04:05")

	objTime := time.Time{}

	if err := petadb.MapIntObject(&objTime, objMap); err != nil {
		panic(err)
	}

	fmt.Println(objTime)
}
