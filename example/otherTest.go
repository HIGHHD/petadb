package main

import (
	"fmt"
	"reflect"
)

func main() {
	userInfo := UserInfo{}
	reflectSet(&userInfo)
	fmt.Println(userInfo)

	var tempClass TempClass
	tempClass.Item = &userInfo
	reflectSet(tempClass.Item)
	fmt.Println(tempClass.Item)
}

func reflectSet(t interface{}) {
	data := reflect.Indirect(reflect.ValueOf(t))
	for i := 0; i < data.NumField(); i++ {
		filed := data.Field(i)
		var i interface{}
		switch filed.Type().Kind() {
		case reflect.Int:
			i = 1
		case reflect.String:
			i = "gejin"
		}
		filed.Set(reflect.ValueOf(i))
	}
}

type UserInfo struct {
	UserId   int
	UserName string
}

type TempClass struct {
	Item *interface{}
}
