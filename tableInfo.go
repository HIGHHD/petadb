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

func GetTableInfo(t interface{}) (TableInfo, error) {
	var tableInfo TableInfo
	dataStruct := reflect.Indirect(reflect.ValueOf(t))
	if dataStruct.Kind() != reflect.Struct || dataStruct.Kind().String() == "time.Time" {
		return tableInfo, errors.New("请输入Sturct类型的对象指针")
	}

	dataStructType := dataStruct.Type()

	tableName := dataStructType.Name()
	var primaryKey string
	var autoIncrement = false
	cols := make(map[string]interface{}, 0)

	hasSetPK := false
	for i := 0; i < dataStructType.NumField(); i++ {
		field := dataStructType.Field(i)
		fieldName := field.Name
		tag := strings.ToLower(reflect.ValueOf(field.Tag).String())

		if tag != "notmap" {
			cols[fieldName] = dataStruct.Field(i).Interface()
		}

		if tag == "pk" || tag == "pk_notauto" {
			if hasSetPK {
				return tableInfo, errors.New("多次设置主键")
			}

			primaryKey = fieldName
			hasSetPK = true
			autoIncrement = tag == "pk"
		}
	}

	if !hasSetPK {
		for i := 0; i < dataStructType.NumField(); i++ {
			field := dataStructType.Field(i)
			fieldName := field.Name

			if strings.ToLower(fieldName) == "id" {
				primaryKey = fieldName
				autoIncrement = true
				hasSetPK = true
			}
		}
	}

	if !hasSetPK {
		return tableInfo, errors.New("没有设置主键")
	}

	tableInfo.AutoIncrement = autoIncrement
	tableInfo.PrimaryKey = primaryKey
	tableInfo.Columns = cols
	tableInfo.TableName = tableName
	return tableInfo, nil
}
