package petadb

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Database struct {
	DbType         string
	DriverName     string
	DataSourceName string
	IsDebug        bool
}

func NewDatabase(dbType string, driverName string, dataSourceName string, isDebug bool) Database {
	return Database{DbType: dbType, DriverName: driverName, DataSourceName: dataSourceName, IsDebug: isDebug}
}

func (database *Database) Insert(t interface{}) (int64, error) {
	mapper, err := getStructMapper(t)
	if err != nil {
		return -1, err
	}

	tableName := database.escapeTableName(mapper.TableName)
	cols := make([]string, 0)
	args := make([]interface{}, 0)
	argHolders := make([]string, 0)

	index := 0

	for k, v := range mapper.Columns {
		if strings.ToLower(k) == strings.ToLower(mapper.PrimaryKey) && mapper.AutoIncrement {
			continue
		}

		cols = append(cols, fmt.Sprintf("%v.%v", tableName, database.escapeSqlIdentifier(k)))
		args = append(args, v)
		argHolders = append(argHolders, database.escapeParamHolder(index))
		index++
	}

	query := fmt.Sprintf("INSERT INTO %v(%v) VALUES(%v)", tableName,
		strings.Join(cols, ","),
		strings.Join(argHolders, ","))

	res, err := database.Execute(query, args...)
	if mapper.AutoIncrement {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, err
		}

		reflect.Indirect(reflect.ValueOf(t)).FieldByName(mapper.PrimaryKey).SetInt(id)
		return id, nil
	}

	row, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return row, nil
}

func (database *Database) Update(t interface{}) (int64, error) {
	mapper, err := getStructMapper(t)
	if err != nil {
		return -1, err
	}

	tableName := database.escapeTableName(mapper.TableName)
	args := make([]interface{}, 0)
	updates := make([]string, 0)

	index := 0
	for k, v := range mapper.Columns {
		if strings.ToLower(k) == strings.ToLower(mapper.PrimaryKey) {
			continue
		}
		args = append(args, v)
		updates = append(updates, fmt.Sprintf("%v.%v = %v", tableName, database.escapeSqlIdentifier(k), database.escapeParamHolder(index)))
		index++
	}

	condition := fmt.Sprintf("%v.%v = %v", tableName, database.escapeSqlIdentifier(mapper.PrimaryKey), database.escapeParamHolder(index))
	args = append(args, mapper.Columns[mapper.PrimaryKey])

	query := fmt.Sprintf("UPDATE %v SET %v WHERE %v", tableName, strings.Join(updates, ","), condition)
	res, err := database.Execute(query, args...)
	if err != nil {
		return -1, err
	}

	rowAff, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return rowAff, nil
}

func (database *Database) UpdateSql(sqlBuilder *SqlBuilder) (int64, error) {
	res, err := database.Execute(sqlBuilder.SQL, sqlBuilder.Args...)
	if err != nil {
		return -1, err
	}

	rowAff, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return rowAff, nil
}

func (database *Database) Delete(t interface{}) (int64, error) {
	mapper, err := getStructMapper(t)
	if err != nil {
		return -1, err
	}

	args := make([]interface{}, 0)
	args = append(args, mapper.Columns[mapper.PrimaryKey])
	tableName := database.escapeTableName(mapper.TableName)
	condition := fmt.Sprintf("%v.%v = %v", tableName, database.escapeSqlIdentifier(mapper.PrimaryKey), database.escapeParamHolder(0))

	query := fmt.Sprintf("DELETE FROM %v WHERE %v", tableName, condition)
	res, err := database.Execute(query, args...)
	if err != nil {
		return -1, err
	}

	rowAff, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return rowAff, nil
}

func (database *Database) DeleteSql(t interface{}, sqlBuilder *SqlBuilder) (int64, error) {
	mapper, err := getStructMapper(t)
	if err != nil {
		return -1, err
	}

	tableName := database.escapeTableName(mapper.TableName)
	query := fmt.Sprintf("DELETE FROM %v %v", tableName, sqlBuilder.SQL)

	res, err := database.Execute(query, sqlBuilder.Args...)
	if err != nil {
		return -1, err
	}

	rowAff, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return rowAff, nil
}

func (database *Database) FindOne(t interface{}, query string, args ...interface{}) (bool, error) {

	var err error
	structData := reflect.Indirect(reflect.ValueOf(t))
	structType := structData.Type()

	isStruct := structType.Kind() == reflect.Struct && structType.Kind().String() != "time.Time"

	if isStruct {
		query, err = database.addSelectClause(t, query)
		if err != nil {
			return false, err
		}
	}

	if query, err = database.processParam(query, args...); err != nil {
		return false, err
	}

	readerList, err := database.executeReader(query, args...)
	if err != nil {
		return false, err
	}

	isExists := len(readerList) > 0
	if isExists {
		err := readerToObject(t, readerList[0])
		if err != nil {
			return isExists, err
		}
	}
	return isExists, nil
}

func (database *Database) FindOneSql(t interface{}, sqlBuilder *SqlBuilder) (bool, error) {
	return database.FindOne(t, sqlBuilder.SQL, sqlBuilder.Args...)
}

func (database *Database) FindPagedList(pagedInfo *PagedInfo, slice interface{}, pageIndex int, pageSize int, query string, args ...interface{}) (err error) {
	sliceData := reflect.Indirect(reflect.ValueOf(slice))
	if sliceData.Kind() != reflect.Slice {
		return errors.New("请输入数组的指针对象")
	}

	sliceElementType := sliceData.Type().Elem()
	element := reflect.New(sliceElementType)
	isStruct := sliceElementType.Kind() == reflect.Struct && sliceElementType.Kind().String() != "time.Time"

	if isStruct {
		query, err = database.addSelectClause(element.Interface(), query)
		if err != nil {
			return err
		}
	}

	if query, err = database.processParam(query, args...); err != nil {
		return err
	}

	sqlCount, sqlPage, err := database.buildPagingQueries((pageIndex-1)*pageSize, pageSize, query)
	if err != nil {
		return err
	}

	_, err = database.FindOne(&pagedInfo.TotalItemCount, sqlCount, args...)
	if err != nil {
		return err
	}

	if err := database.FindList(slice, sqlPage, args...); err != nil {
		return err
	}

	pagedInfo.TotalPageCount = int(pagedInfo.TotalItemCount / pageSize)
	pagedInfo.CurrentPageIndex = pageIndex
	pagedInfo.PageSize = pageSize

	return nil
}

func (database *Database) FindList(slice interface{}, query string, args ...interface{}) error {
	var err error

	sliceDataStruct := reflect.Indirect(reflect.ValueOf(slice))
	if sliceDataStruct.Kind() != reflect.Slice {
		return errors.New("请输入数组的指针对象")
	}

	sliceElementType := sliceDataStruct.Type().Elem()
	element := reflect.New(sliceElementType)
	isStruct := sliceElementType.Kind() == reflect.Struct && sliceElementType.String() != "time.Time"

	if isStruct {
		query, err = database.addSelectClause(element.Interface(), query)
		if err != nil {
			return err
		}
	}

	if query, err = database.processParam(query, args...); err != nil {
		return err
	}

	reader, err := database.executeReader(query, args...)
	if err != nil {
		return err
	}

	for i := 0; i < len(reader); i++ {
		newStructValue := reflect.New(sliceElementType)
		err := readerToObject(newStructValue.Interface(), reader[i])
		if err != nil {
			return err
		}
		sliceDataStruct.Set(reflect.Append(sliceDataStruct, reflect.Indirect(reflect.ValueOf(newStructValue.Interface()))))
	}
	return nil
}

func (database *Database) FindListSql(slice interface{}, sqlBuilder *SqlBuilder) error {
	return database.FindList(slice, sqlBuilder.SQL, sqlBuilder.Args...)
}

func (database *Database) Execute(query string, args ...interface{}) (sql.Result, error) {
	if database.IsDebug {
		fmt.Println(query)
		fmt.Println(args)
	}

	finalQuery, err := database.processParam(query, args)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(database.DriverName, database.DataSourceName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(finalQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (database *Database) processParam(query string, args ...interface{}) (string, error) {
	regx, err := regexp.Compile("\\@\\d*")
	if err != nil {
		return query, err
	}

	index := 0
	var finalSql string
	switch strings.ToLower(database.DbType) {
	case "mysql":
		finalSql = regx.ReplaceAllString(query, "?")
	case "mssql":
		finalSql = regx.ReplaceAllStringFunc(query, func(string) string {
			str := fmt.Sprintf("@%d", index)
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
		mapper, err := getStructMapper(t)
		if err != nil {
			return query, err
		}

		tableName := database.escapeTableName(mapper.TableName)
		cols := make([]string, 0)
		for k, _ := range mapper.Columns {
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

func (database *Database) executeReader(query string, args ...interface{}) ([]map[string][]byte, error) {
	if database.IsDebug {
		fmt.Println(query)
		fmt.Println(args)
	}

	db, err := sql.Open(database.DriverName, database.DataSourceName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var resultSlice []map[string][]byte

	for rows.Next() {
		result := make(map[string][]byte, 0)
		var dest []interface{}
		for i := 0; i < len(cols); i++ {
			var destItem interface{}
			dest = append(dest, &destItem)
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}

		for index, keyName := range cols {
			colData := reflect.Indirect(reflect.ValueOf(dest[index]))
			if colData.Interface() == nil {
				continue
			}

			colType := reflect.TypeOf(colData.Interface())
			colValue := reflect.ValueOf(colData.Interface())

			var tempStr string
			switch colType.Kind() {
			case reflect.String:
				tempStr = colValue.String()
				result[keyName] = []byte(tempStr)
			case reflect.Slice:
				if colType.Elem().Kind() == reflect.Uint8 {
					result[keyName] = colValue.Interface().([]byte)
				}
				break
			case reflect.Float32, reflect.Float64:
				tempStr = strconv.FormatFloat(colValue.Float(), 'f', -1, 64)
				result[keyName] = []byte(tempStr)
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
				tempStr = strconv.FormatInt(colValue.Int(), 10)
				result[keyName] = []byte(tempStr)
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				tempStr = strconv.FormatUint(colValue.Uint(), 10)
				result[keyName] = []byte(tempStr)
			case reflect.Struct:
				tempStr = colValue.Interface().(time.Time).Format("2006-01-02 15:04:05.000 -0700")
				result[keyName] = []byte(tempStr)
				fmt.Println(keyName)
			}
		}
		resultSlice = append(resultSlice, result)
	}
	return resultSlice, nil
}

func (database *Database) buildPagingQueries(skip int, take int, querySql string) (sqlCount string, sqlPage string, err error) {
	regColumns, _ := regexp.Compile("select(.*)from(.*)")
	regSelect, _ := regexp.Compile("select ")
	regOrderBy, _ := regexp.Compile("order by (.*)")
	querySqlBytes := []byte(strings.ToLower(querySql))
	matchColumnsArr := regColumns.FindSubmatch(querySqlBytes)
	matchOrderByArr := regOrderBy.FindSubmatch(querySqlBytes)

	var removeSelect string
	var columns string
	var afterFrom string
	var orderBy string

	removeSelect = regSelect.ReplaceAllString(strings.ToLower(querySql), "")
	removeSelect = regOrderBy.ReplaceAllString(removeSelect, "")

	if len(matchColumnsArr) >= 2 {
		columns = string(matchColumnsArr[1])
	}

	if len(matchColumnsArr) >= 3 {
		afterFrom = string(matchColumnsArr[2])
	}

	if len(matchColumnsArr) == 0 {
		return sqlCount, sqlPage, errors.New("创建分页SQL时出现异常")
	}

	if len(matchOrderByArr) == 2 {
		orderBy = string(matchOrderByArr[1])
	}

	if strings.Contains(columns, "distinct ") {
		sqlCount = fmt.Sprintf("select count(%v) AS Num from %v", columns, afterFrom)
	} else {
		sqlCount = fmt.Sprintf("select count(1) AS Num from %v", afterFrom)
	}

	switch strings.ToLower(database.DbType) {
	case "mysql":
		sqlPage = fmt.Sprintf("%v LIMIT %d OFFSET %d", querySql, take, skip)
	case "mssql":
		if strings.Contains(columns, "distinct ") {
			columns = "peta_inner.* FROM (SELECT " + removeSelect + ") peta_inner"
		}

		if len(orderBy) == 0 {
			orderBy = "ORDER BY (SELECT NULL)"
		}

		sqlPage = fmt.Sprintf("SELECT * FROM (SELECT ROW_NUMBER() OVER (%v) peta_rn, %v) peta_paged WHERE peta_rn > %d AND peta_rn <= %d", orderBy, removeSelect, skip, take+skip)
	}

	return sqlCount, sqlPage, nil
}
