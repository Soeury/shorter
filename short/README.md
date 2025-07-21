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
# 生成代码(不带缓存版)
goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/sequence" -table="sequence" -dir="./model"
goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/short" -table="reflect_map" -dir="./model"
goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/short2" -table="reflect_map2" -dir="./model"


# 带缓存版
goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/short" -table="reflect_map" -dir="./model" -c 
goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/short2" -table="reflect_map2" -dir="./model" -c 


# 下载依赖
go mod tidy
```



### 修改配置文件
增加链接数据库配置



### Start
* 1. 校验长链: validate, 修改handler部分和api，增加validate校验部分，长链不为空
* 2. 长链有效：实现自定义部分connect, 调用检查长链是否能get并返回200
* 3. 获取md5值, 查询长链是否被转存过
* 4. 获取URL path, 查询传入的长链是否是短链
* 5. 针对urltool, md5, base62编写单元测试
* 6. 实现发号器 (mysql:replace, redis:incr) 复用接口
* 7. 号码转短链：10 ---> 62进制
* 8. 短链特殊词检查(准备短链黑名单, map检查特殊词)
* 9. 存储长短链映射, 返回响应
* 10. 编写重定向模块(show), handler进行validate校验，logic编写查询逻辑并返回响应, 拿到长链后handler进行重定向
* 11. 编写重定向模块缓存, 生成带缓存的reflect_map_model层代码, 内嵌singleflight进行请求合并
* 12. 增加fileter防止缓存击穿，生成的所有短链会保存到filter, 查询长链进行重定向时会先查询过滤器。(做缓存击穿)
* 13. 拆分奇偶表和取号表, 采用3个不同的数据库, (存储分离, 查询分离) 
