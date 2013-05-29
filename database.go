package petadb

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Database struct {
	DbType         string
	DriverName     string
	DataSourceName string
	IsDebug        bool
}

type PagedList struct {
	CurrentPageIndex int
	PageSize         int
	List             []interface{}
	TotalItemCount   int64
	TotalPageCount   int64
}

func NewDatabase(dbType string, driverName string, dataSourceName string) Database {
	return Database{DbType: dbType, DataSourceName: dataSourceName, DriverName: driverName, IsDebug: false}
}

func (database *Database) FirstBySb(t interface{}, sqlBuilder *SqlBuilder) (bool, error) {
	return database.First(t, sqlBuilder.SQL, sqlBuilder.Args...)
}

func (database *Database) First(t interface{}, query string, args ...interface{}) (bool, error) {
	err := errors.New("")
	structData := reflect.Indirect(reflect.ValueOf(t))
	structType := structData.Type()

	isStruct := structType.Kind() == reflect.Struct && structType.Kind().String() != "time.Time"

	if isStruct {
		query, err = database.addSelectClause(t, query)
		if err != nil {
			return false, err
		}
	}

	if query, err = database.ProcessParam(query, args...); err != nil {
		return false, err
	}

	if database.IsDebug {
		fmt.Println(query)
		fmt.Println(args)
	}

	db, err := sql.Open(database.DriverName, database.DataSourceName)
	if err != nil {
		return false, err
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return false, err
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return false, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return false, err
	}

	if isStruct {
		tableInfo, err := GetTableInfo(t)
		if err != nil {
			return false, err
		}

		colToFieldOffset := make([]int, len(cols))
		for i := 0; i < len(cols); i++ {
			colToFieldOffset[i] = -1
			for k, _ := range tableInfo.Columns {
				if strings.ToLower(k) == strings.ToLower(cols[i]) {
					colToFieldOffset[i] = 1
					break
				}
			}
			if colToFieldOffset[i] == -1 {
				return false, errors.New("查询的列与Struct不一致")
			}
		}
	}

	isExists := false
	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return false, err
			}
			break
		}

		dest := make([]interface{}, 0)
		for i := 0; i < len(cols); i++ {
			if isStruct {
				dest = append(dest, structData.FieldByName(cols[i]).Addr().Interface())
			} else {
				dest = append(dest, structData.Interface())
			}
		}

		if err := rows.Scan(dest...); err != nil {
			return false, err
		}
		isExists = true
		break
	}

	return isExists, nil

}

func (database *Database) QueryBySb(sliceInput interface{}, sqlBuilder *SqlBuilder) error {
	return database.Query(sliceInput, sqlBuilder.SQL, sqlBuilder.Args...)
}

func (database *Database) Query(sliceInput interface{}, query string, args ...interface{}) error {
	err := errors.New("")
	sliceDataStruct := reflect.Indirect(reflect.ValueOf(sliceInput))
	if sliceDataStruct.Kind() != reflect.Slice {
		return errors.New("请输入数组的指针对象")
	}

	sliceElementType := sliceDataStruct.Type().Elem()
	element := reflect.New(sliceElementType)
	isStruct := sliceElementType.Kind() == reflect.Struct && sliceElementType.Kind().String() != "time.Time"

	if isStruct {
		query, err = database.addSelectClause(element.Interface(), query)
		if err != nil {
			return err
		}
	}

	if query, err = database.ProcessParam(query, args...); err != nil {
		return err
	}

	if database.IsDebug {
		fmt.Println(query)
		fmt.Println(args)
	}

	db, err := sql.Open(database.DriverName, database.DataSourceName)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return err
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	if isStruct {
		tableInfo, err := GetTableInfo(element.Interface())
		if err != nil {
			return err
		}

		colToFieldOffset := make([]int, len(cols))
		for i := 0; i < len(cols); i++ {
			colToFieldOffset[i] = -1
			for k, _ := range tableInfo.Columns {
				if strings.ToLower(k) == strings.ToLower(cols[i]) {
					colToFieldOffset[i] = 1
					break
				}
			}
			if colToFieldOffset[i] == -1 {
				return errors.New("查询的列与Struct不一致")
			}
		}
	}

	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return err
			}
			break
		}

		v := reflect.New(sliceElementType)
		dest := make([]interface{}, 0)
		for i := 0; i < len(cols); i++ {
			if isStruct {
				dest = append(dest, v.Elem().FieldByName(cols[i]).Addr().Interface())
			} else {
				dest = append(dest, v.Interface())
			}
		}

		if err := rows.Scan(dest...); err != nil {
			return err
		}
		sliceDataStruct.Set(reflect.Append(sliceDataStruct, reflect.Indirect(reflect.ValueOf(v.Interface()))))
	}

	return nil

}

func (database *Database) DeleteBySb(t interface{}, sqlBuilder *SqlBuilder) (int64, error) {
	tableInfo, err := GetTableInfo(t)
	if err != nil {
		return -1, err
	}

	tableName := database.escapeTableName(tableInfo.TableName)

	query := fmt.Sprintf("DELETE FROM %v %v", tableName, sqlBuilder.SQL)
	return database.execAndReturnRowsAffected(query, sqlBuilder.Args...)
}

func (database *Database) Delete(t interface{}) (int64, error) {
	tableInfo, err := GetTableInfo(t)
	if err != nil {
		return -1, err
	}

	tableName := database.escapeTableName(tableInfo.TableName)
	condition := fmt.Sprintf("%v.%v = %v", tableName, database.escapeSqlIdentifier(tableInfo.PrimaryKey), database.escapeParamHolder(0))
	args := make([]interface{}, 0)
	args = append(args, tableInfo.Columns[tableInfo.PrimaryKey])

	query := fmt.Sprintf("DELETE FROM %v WHERE %v", tableName, condition)

	return database.execAndReturnRowsAffected(query, args...)
}

func (database *Database) UpdateBySb(sqlBuilder *SqlBuilder) (int64, error) {
	return database.execAndReturnRowsAffected(sqlBuilder.SQL, sqlBuilder.Args...)
}

func (database *Database) Update(t interface{}) (int64, error) {
	tableInfo, err := GetTableInfo(t)
	if err != nil {
		return -1, err
	}

	tableName := database.escapeTableName(tableInfo.TableName)
	args := make([]interface{}, 0)
	updates := make([]string, 0)

	index := 0
	for k, v := range tableInfo.Columns {
		if strings.ToLower(k) == strings.ToLower(tableInfo.PrimaryKey) {
			continue
		}

		updates = append(updates, fmt.Sprintf("%v.%v = %v", tableName, database.escapeSqlIdentifier(k), database.escapeParamHolder(index)))
		args = append(args, v)
		index++
	}

	condition := fmt.Sprintf("%v.%v = %v", tableName, database.escapeSqlIdentifier(tableInfo.PrimaryKey), database.escapeParamHolder(index))
	args = append(args, tableInfo.Columns[tableInfo.PrimaryKey])

	query := fmt.Sprintf("UPDATE %v SET %v WHERE %v", tableName, strings.Join(updates, ","), condition)
	return database.execAndReturnRowsAffected(query, args...)
}

func (database *Database) Insert(t interface{}) (int64, error) {
	tableInfo, err := GetTableInfo(t)
	if err != nil {
		return -1, err
	}

	tableName := database.escapeTableName(tableInfo.TableName)
	cols := make([]string, 0)
	args := make([]interface{}, 0)
	argsHolders := make([]string, 0)

	index := 0
	for k, v := range tableInfo.Columns {
		if strings.ToLower(k) == strings.ToLower(tableInfo.PrimaryKey) && tableInfo.AutoIncrement {
			continue
		}

		cols = append(cols, fmt.Sprintf("%v.%v", tableName, database.escapeSqlIdentifier(k)))
		args = append(args, v)
		argsHolders = append(argsHolders, database.escapeParamHolder(index))
		index++
	}

	query := fmt.Sprintf("INSERT INTO %v(%v) VALUES(%v)", tableName, strings.Join(cols, ","), strings.Join(argsHolders, ","))

	if tableInfo.AutoIncrement {
		id, err := database.execAndReturnLastInsertId(query, args...)
		if err != nil {
			return -1, err
		}
		reflect.Indirect(reflect.ValueOf(t)).FieldByName(tableInfo.PrimaryKey).SetInt(id)
		return id, nil
	}

	return database.execAndReturnRowsAffected(query, args...)
}

func (database *Database) execAndReturnRowsAffected(query string, args ...interface{}) (int64, error) {
	res, err := database.Exec(query, args...)
	if err != nil {
		return -1, err
	}

	row, err := res.RowsAffected()
	if err != nil {
		return -1, nil
	}
	return row, nil
}

func (database *Database) execAndReturnLastInsertId(query string, args ...interface{}) (int64, error) {
	res, err := database.Exec(query, args...)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (database *Database) Exec(query string, args ...interface{}) (sql.Result, error) {

	if database.IsDebug {
		fmt.Println(query)
		fmt.Println(args)
	}

	finalQuery, err := database.ProcessParam(query, args)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(database.DriverName, database.DataSourceName)
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(finalQuery)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

var regxParam = regexp.MustCompile("=\\s*\\@\\d*")

func (database *Database) ProcessParam(query string, args ...interface{}) (string, error) {
	regx, err := regexp.Compile("=\\s\\@\\d*")
	if err != nil {
		return query, err
	}

	index := 0
	var finalSql string
	switch strings.ToLower(database.DbType) {
	case "mysql":
		finalSql = regx.ReplaceAllString(query, "= ?")
	case "mssql":
		finalSql = regx.ReplaceAllStringFunc(query, func(string) string {
			str := fmt.Sprintf("= @%d", index)
			index++
			return str
		})
	}

	return finalSql, nil
}

var regxSelect = regexp.MustCompile(`\A\s*(SELECT|EXECUTE|CALL)\s`)
var regxFrom = regexp.MustCompile(`\A\s*FROM\s`)

func (database *Database) addSelectClause(t interface{}, query string) (string, error) {
	if query == "" {
		return query, nil
	}

	if query[0] == ';' {
		return query, nil
	}

	if !regxSelect.MatchString(strings.ToUpper(query)) {
		tableInfo, err := GetTableInfo(t)
		if err != nil {
			return query, err
		}

		tableName := database.escapeTableName(tableInfo.TableName)
		cols := make([]string, 0)
		for k, _ := range tableInfo.Columns {
			cols = append(cols, fmt.Sprintf("%v.%v", tableName, database.escapeSqlIdentifier(k)))
		}

		if regxFrom.MatchString(strings.ToUpper(query)) {
			return fmt.Sprintf("SELECT %v %v", strings.Join(cols, ","), query), nil
		} else {
			return fmt.Sprintf("SELECT %v FROM %v %v", strings.Join(cols, ","), tableName, query), nil
		}
	}
	return query, nil

}

func (database *Database) escapeTableName(table string) string {
	if strings.Index(table, ".") >= 0 {
		return table
	}
	return database.escapeSqlIdentifier(table)
}

func (database *Database) escapeSqlIdentifier(str string) string {
	switch strings.ToLower(database.DbType) {
	case "mysql":
		return fmt.Sprintf("`%v`", str)
	case "mssql":
		return fmt.Sprintf("[%v]", str)
	}
	return str
}

func (database *Database) escapeParamHolder(paramIndex int) string {
	switch strings.ToLower(database.DbType) {
	case "mysql":
		return "?"
	case "mssql":
		return fmt.Sprintf("@%d", paramIndex)
	}
	return ""
}
