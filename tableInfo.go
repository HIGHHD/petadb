package petadb

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type TableInfo struct {
	TableName     string
	PrimaryKey    string
	AutoIncrement bool
	Columns       map[string]interface{}
}

func MapIntObject(t interface{}, objMap map[string][]byte) error {
	data := reflect.Indirect(reflect.ValueOf(t))
	fmt.Println(data.Type())
	// if t is struct (not time.Time)
	if data.Kind() == reflect.Struct && data.Type().String() != "time.Time" {
		for k, v := range objMap {
			field := data.FieldByName(k)

			if !field.CanSet() {
				continue
			}

			var item interface{}

			switch field.Type().Kind() {
			case reflect.Slice:
				item = v
			case reflect.String:
				item = string(v)
			case reflect.Bool:
				item = string(v) == "1"
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
				i, err := strconv.Atoi(string(v))
				if err != nil {
					return err
				}
				item = i
			case reflect.Int64:
				i, err := strconv.ParseInt(string(v), 10, 64)
				if err != nil {
					return err
				}
				item = i
			case reflect.Float32, reflect.Float64:
				i, err := strconv.ParseFloat(string(v), 64)
				if err != nil {
					return err
				}
				item = i
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				i, err := strconv.ParseUint(string(v), 10, 64)
				if err != nil {
					return err
				}
				item = i
			case reflect.Struct:
				if field.Type().String() != "time.Time" {
					return errors.New("字段类型不支持:" + field.Type().Kind().String())
				}
				i, err := time.Parse("2013-05-29 15:04:05", string(v))
				if err != nil {
					i, err = time.Parse("2013-05-29 15:04:05.000 -0700", string(v))
					if err != nil {
						return errors.New("不支持的时间格式:" + string(v))
					}
				}
				item = i
			default:
				return errors.New("字段类型不支持:" + field.Type().Kind().String())
			}
			field.Set(reflect.ValueOf(item))
		}
	} else {
		for _, v := range objMap {
			var item interface{}
			switch data.Kind() {
			case reflect.Slice:
				item = v
			case reflect.String:
				item = string(v)
			case reflect.Bool:
				item = string(v) == "1"
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
				i, err := strconv.Atoi(string(v))
				if err != nil {
					return err
				}
				item = i
			case reflect.Int64:
				i, err := strconv.ParseInt(string(v), 10, 64)
				if err != nil {
					return err
				}
				item = i
			case reflect.Float32, reflect.Float64:
				i, err := strconv.ParseFloat(string(v), 64)
				if err != nil {
					return err
				}
				item = i
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				i, err := strconv.ParseUint(string(v), 10, 64)
				if err != nil {
					return err
				}
				item = i
			case reflect.Struct:
				if data.Type().String() != "time.Time" {
					return errors.New("字段类型不支持:" + data.Kind().String())
				}
				i, err := time.Parse("2006-01-02 15:04:05", string(v))
				if err != nil {
					i, err = time.Parse("2006-01-02 15:04:05.000 -0700", string(v))
					if err != nil {
						return errors.New("不支持的时间格式:" + string(v))
					}
				}
				item = i
			default:
				return errors.New("字段类型不支持:" + data.Kind().String())
			}
			fmt.Println(item)
			data.Set(reflect.ValueOf(item))
			break
		}
	}

	return nil
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
