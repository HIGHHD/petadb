# petadb
一个简洁的Golang ORM

##  快速开始

### 1.初始化
在main package 中，初始化一个全局的database,如下：
 ```go
var database = petadb.NewDatabase("mysql", "mysql", "root:123456@/test?charset=utf8", false)
```

### 2.增删改查
* 前置条件：
	test数据库中存在UserInfo数据表：UserId(int),UserName(varchar(32)),CreateDate(DateTime)
	