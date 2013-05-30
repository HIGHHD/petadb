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
 