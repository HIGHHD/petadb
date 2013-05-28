package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	regColumns, _ := regexp.Compile("select(.*)from(.*)")
	querySql := "SELECT UserId FROM UserInfo WHERE UserName = @0 ORDER BY CreateDate DESC"
	matchColumnsArr := regColumns.FindSubmatch([]byte(strings.ToLower(querySql)))

	var columns string
	var afterFrom string
	var sqlCount string

	if len(matchColumnsArr) >= 2 {
		columns = string(matchColumnsArr[1])
	}

	if len(matchColumnsArr) >= 3 {
		afterFrom = string(matchColumnsArr[2])
	}

	if len(columns) == 0 {
		fmt.Println("error")
		return
	}

	if strings.Contains(columns, "distinct ") {
		sqlCount = fmt.Sprintf("select count(%v) from %v", columns, afterFrom)
	} else {
		sqlCount = fmt.Sprintf("select count(1) from %v", afterFrom)
	}

	sqlPage := fmt.Sprintf("%v LIMIT %d OFFSET %d", querySql, 1, 10)
	fmt.Printf(sqlCount)
	fmt.Printf(sqlPage)
}
