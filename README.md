# shorturl project



### 编写 api 文件
两个api
* convertHandler: 长链转短链
* showHandler: 查询短链对应的长链



### 根据api文件生成go-zero框架代码
```bash
# 生成代码
goctl api go -api short.api -dir .

#下载依赖
go mod tidy
```



### 编写数据库表，生成数据库部分代码
两张表
* 根据主键生成的自动发号器
* 长短链接映射表


```bash
# 生成代码
goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/short" -table="sequence" -dir="./model"
goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/short" -table="reflect_map" -dir="./model"

# 下载依赖
go mod tidy
```



### 修改配置文件
增加链接数据库配置



### Start
* 1.校验长链: validate
