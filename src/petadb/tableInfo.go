package petadb

import (
	"errors"
	"reflect"
	"strings"
)

type TableInfo struct {
	TableName     string
	PrimaryKey    string
	AutoIncrement bool
	Columns       map[string]interface{}
}

func getTableInfo(t interface{}) (TableInfo, error) {
	var tableInfo TableInfo

	// 获取类型信息
	dataStruct := reflect.Indirect(reflect.ValueOf(t))
	if dataStruct.Kind() != reflect.Struct || dataStruct.Kind().String() == "time.Time" {
		return tableInfo, errors.New("请确认输入的是Struct类型（非时间类型）的指针对象")
	}

	dataStructType := dataStruct.Type()

	// 获取struct对应的表名，主键名，是否自增，字段列表
	tableName := dataStructType.Name()
	var primaryKey string
	var autoIncrement bool
	cols := make(map[string]interface{}, 0)
	hasSetPrimaryKey := false

	for i := 0; i < dataStructType.NumField(); i++ {
		field := dataStructType.Field(i)
		fieldName := field.Name
		tag := field.Tag.Get("petadb")

		if len(tag) == 0 || strings.ToLower(tag) != "notmap" {
			cols[fieldName] = dataStruct.FieldByName(fieldName).Interface()
		}

		if strings.ToLower(tag) == "pk" {
			if hasSetPrimaryKey {
				return tableInfo, errors.New("设置过多的主键")
			}

			primaryKey = fieldName
			autoIncrement = true
			hasSetPrimaryKey = true
		} else if strings.ToLower(tag) == "pk_na" {
			if hasSetPrimaryKey {
				return tableInfo, errors.New("设置过多的主键")
			}

			primaryKey = fieldName
			autoIncrement = false
			hasSetPrimaryKey = true
		}
	}

	if !hasSetPrimaryKey {
		for i := 0; i < dataStructType.NumField(); i++ {
			field := dataStructType.Field(i)
			fieldName := field.Name

			if strings.ToLower(fieldName) == "id" {
				primaryKey = fieldName
				autoIncrement = true
				hasSetPrimaryKey = true
				break
			}
		}
	}

	if !hasSetPrimaryKey {
		return tableInfo, errors.New("未能在对象中找到对应的主键")
	}

	tableInfo.TableName = tableName
	tableInfo.PrimaryKey = primaryKey
	tableInfo.AutoIncrement = autoIncrement
	tableInfo.Columns = cols
	return tableInfo, nil

}
