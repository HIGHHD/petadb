petadb
=====

一个微型Golang ORM，思路来源于.NET社区的PetaPoco微型ORM

### 组件介绍
sqlBuilder.go ：主要用于SQL语句的组装，通过此组件的组装，SQL语句会自动匹配不同种类的数据库，现暂支持：mysql与MS Sql Server 

mapper.go ：主要用于poco类型与数据库表之间的映射，database/sql组件读取数据库数据后自动反射为poco对象

database.go：主要用于数据库的查询，已经封装好各种查询API

### 快速开始 
#### 初始化 

在main package中，初始化一个全局的database对象，如下：
```go
	var database = petadb.NewDatabase("mysql", "mysql", "root:123456@/test?charset=utf8")
```

#### 增删改查
##### 环境
数据库类型：mysql

数据库：test

数据表： UserInfo => UserId (int)   UserName (varchar(32))   CreateDate(DateTime)

##### 映射
 
poco实体类：
```go

type UserInfo struct {
	UserId     int `petadb:"pk"` // 主键自增，若属性名为Id时，则默认为自增主键，非自增主键：pk_notai
	UserName   string
	CreateDate time.Time
	Other string `petadb:"notmap"`  // 不映射至数据表字段
}
``` 

##### 新增 database.Insert
```go
func insertTest() {
	userInfo := UserInfo{UserName: "gejin", CreateDate: time.Now()}
	id, err := database.Insert(&userInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(userInfo)

	fmt.Println(id)
}
```
##### 查找第一条 database.FindOne
```go 
func FindOneTest() {
	userInfo := UserInfo{}
	isEixsts, err := database.FindOne(&userInfo, "SELECT * FROM UserInfo WHERE UserName = 'gejin'")
	if err != nil {
		panic(err)
	}

	fmt.Println(isEixsts)
	// 是否存在
	if isEixsts {
		fmt.Println(userInfo)
	}
}
```

##### 修改 database.Update
```go 
func UpdateTest() {
	var userInfo UserInfo
	// 取出要修改的对象
	isEixsts, err := database.FindOne(&userInfo, "SELECT * FROM UserInfo WHERE UserName = 'gejin'")
	if err != nil {
		panic(err)
	}

	fmt.Println(isEixsts)
	// 如果存在，则Update
	if isEixsts {
		fmt.Println(userInfo)

		userInfo.CreateDate = time.Now()
		row, err := database.Update(&userInfo)
		if err != nil {
			panic(err)
		}

		fmt.Println(row)
		fmt.Println(userInfo)
	}
}
```

##### 删除 database.Delete 

```go 
func DeleteTest() {
	var userInfo UserInfo
	// 取出要修改的对象
	isEixsts, err := database.FindOne(&userInfo, "SELECT * FROM UserInfo WHERE UserName = 'gejin'")
	if err != nil {
		panic(err)
	}

	fmt.Println(isEixsts)
	// 如果存在，则Delete
	if isEixsts {
		fmt.Println(userInfo)

		row, err := database.Delete(&userInfo)
		if err != nil {
			panic(err)
		}

		fmt.Println(row)
	}
}
```

##### 查找列表 database.FindList
```go  
func FindListTest() {
	var userList []UserInfo
	if err := database.FindList(&userList, "SELECT * FROM UserInfo"); err != nil {
		panic(err)
	}

	fmt.Println(userList)
}
```

##### 分页查询  database.FindPagedList
```go 
func PagedListTest() {
	var pagedInfo petadb.PagedInfo
	userList := make([]UserInfo, 0)

	// SQL语句会自动转换为分页语句(1.查询总数语句,2.查询列表语句)
	if err := database.FindPagedList(&pagedInfo, &userList, 1, 10, "SELECT * FROM UserInfo"); err != nil {
		panic(err)
	}

	fmt.Println(pagedInfo)
	fmt.Println(userList)
}
```