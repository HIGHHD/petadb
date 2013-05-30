package petadb

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 映射信息类
type Mapper struct {
	TableName     string
	PrimaryKey    string
	AutoIncrement bool
	Columns       map[string]interface{}
}

// 获取struct对象的映射信息
func getStructMapper(t interface{}) (Mapper, error) {
	var mapper Mapper

	data := reflect.Indirect(reflect.ValueOf(t))
	if data.Kind() != reflect.Struct || data.Type().String() == "time.Time" {
		return mapper, errors.New("请输入Struct类型的指针")
	}

	dataType := data.Type()
	tableName := dataType.Name()
	var primaryKey string
	var autoIncrement bool
	cols := make(map[string]interface{}, 0)
	hasSetPK := false

	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		fieldName := field.Name
		tag := field.Tag.Get("petadb")

		if tag != "notmap" {
			cols[fieldName] = data.Field(i).Interface()
		}

		if tag == "pk" || tag == "pk_notai" {
			if hasSetPK {
				return mapper, errors.New("重复设置主键")
			}

			primaryKey = fieldName
			hasSetPK = true
			autoIncrement = tag == "pk"
		}
	}

	if !hasSetPK {
		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)
			fieldName := field.Name
			if strings.ToLower(fieldName) == "id" {
				primaryKey = fieldName
				autoIncrement = true
				hasSetPK = true
				break
			}
		}
	}

	if !hasSetPK {
		return mapper, errors.New("请设置主键")
	}

	mapper.TableName = tableName
	mapper.AutoIncrement = autoIncrement
	mapper.Columns = cols
	mapper.PrimaryKey = primaryKey
	return mapper, nil
}

// 将数据库reader转成object
func readerToObject(t interface{}, reader map[string][]byte) error {
	data := reflect.Indirect(reflect.ValueOf(t))

	if data.Kind() == reflect.Struct && data.Type().String() != "time.Time" {
		if data.NumField() > len(reader) {
			return errors.New("将reader转成object时，列对应出现不匹配")
		}

		for k, v := range reader {
			field := data.FieldByName(k)
			if field.Interface() == nil || !field.CanSet() {
				continue
			}

			var item interface{}

			switch field.Type().Kind() {
			case reflect.String:
				item = string(v)
			case reflect.Bool:
				item = string(v) == "1"
			case reflect.Slice:
				item = v
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int:
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
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				i, err := strconv.ParseUint(string(v), 10, 64)
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
			case reflect.Struct:
				if field.Type().String() != "time.Time" {
					return errors.New("字段类型不支持:" + field.Type().Kind().String())
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
				return errors.New("字段类型不支持:" + field.Type().Kind().String())
			}
			field.Set(reflect.ValueOf(item))
		}
	} else {
		var item interface{}
		for _, v := range reader {
			switch data.Type().Kind() {
			case reflect.String:
				item = string(v)
			case reflect.Bool:
				item = string(v) == "1"
			case reflect.Slice:
				item = v
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int:
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
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				i, err := strconv.ParseUint(string(v), 10, 64)
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
			case reflect.Struct:
				if data.Type().String() != "time.Time" {
					return errors.New("字段类型不支持:" + data.Type().Kind().String())
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
				return errors.New("字段类型不支持:" + data.Type().Kind().String())
			}
			data.Set(reflect.ValueOf(item))
			break
		}
	}
	return nil
}
