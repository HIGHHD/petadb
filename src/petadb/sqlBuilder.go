package petadb

import (
	"strings"
)

type SqlBuilder struct {
	SQL  string
	Args []interface{}
}

func NewSqlBuilder() SqlBuilder {
	return SqlBuilder{SQL: "", Args: make([]interface{}, 0)}
}

func (sqlBuilder *SqlBuilder) Append(query string, args ...interface{}) *SqlBuilder {
	sqlBuilder.SQL += "\n" + query
	sqlBuilder.Args = append(sqlBuilder.Args, args...)
	return sqlBuilder
}

func (sqlBuilder *SqlBuilder) Select(columns string) *SqlBuilder {
	sqlBuilder.SQL += "SELECT " + columns
	return sqlBuilder
}

func (sqlBuilder *SqlBuilder) From(tables string) *SqlBuilder {
	sqlBuilder.SQL += "\nFROM " + tables
	return sqlBuilder
}

func (sqlBuilder *SqlBuilder) Where(condition string, args ...interface{}) *SqlBuilder {
	if strings.Contains(strings.ToLower(sqlBuilder.SQL), "where ") {
		sqlBuilder.SQL += "\nAND (" + condition + ") "
	} else {
		sqlBuilder.SQL += "\nWHERE (" + condition + ") "
	}
	sqlBuilder.Args = append(sqlBuilder.Args, args...)
	return sqlBuilder
}

func (sqlBuilder *SqlBuilder) GroupBy(groupBy string) *SqlBuilder {
	if strings.Contains(strings.ToLower(sqlBuilder.SQL), "group by ") {
		sqlBuilder.SQL += "," + groupBy + " "
	} else {
		sqlBuilder.SQL += "\nGROUP BY " + groupBy
	}
	return sqlBuilder
}

func (sqlBuilder *SqlBuilder) Having(condition string, args ...interface{}) *SqlBuilder {
	if strings.Contains(strings.ToLower(sqlBuilder.SQL), "having ") {
		sqlBuilder.SQL += "\nAND (" + condition + ") "
	} else {
		sqlBuilder.SQL += "\nHAVING (" + condition + ") "
	}
	sqlBuilder.Args = append(sqlBuilder.Args, args...)
	return sqlBuilder
}

func (sqlBuilder *SqlBuilder) OrderBy(orderBy string) *SqlBuilder {
	if strings.Contains(strings.ToLower(sqlBuilder.SQL), "order by ") {
		sqlBuilder.SQL += "," + orderBy + " "
	} else {
		sqlBuilder.SQL += "\nORDER BY " + orderBy
	}
	return sqlBuilder
}

func (sqlBuilder *SqlBuilder) Join(table string) *SqlBuilder {
	sqlBuilder.SQL += " INNER JOIN " + table
	return sqlBuilder
}

func (sqlBuilder *SqlBuilder) LeftJoin(table string) *SqlBuilder {
	sqlBuilder.SQL += " LEFT JOIN " + table
	return sqlBuilder
}

func (sqlBuilder *SqlBuilder) On(condition string, args ...interface{}) *SqlBuilder {
	sqlBuilder.SQL += "\n    ON (" + condition + ") "
	sqlBuilder.Args = append(sqlBuilder.Args, args...)
	return sqlBuilder
}
