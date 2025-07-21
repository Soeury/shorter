# 短URL链接系统
基于go-zero开发的短URL链接系统，针对读多写少场景实现的长链接生成短链重定向

## 项目结构
```bash
short/
├─etc   # 配置文件
├─internal   # 内部包
│  ├─config   # 配置文件
│  ├─handler  # 控制层
│  ├─logic    # 业务逻辑
│  ├─pkg      # 公共包
│  ├─sequence  # 取号器
│  ├─svc      # 服务层  
│  └─types   # api接口
└─model   # 数据模型
```

## Start
1. 拉取代码
```bash
go clone git@github.com:Soeury/shorter.git
cd short 
```

2. 安装依赖
```bash
go mod tidy 
```

3. 配置变量
```bash 
cp etc/short-example.yaml etc/short-api.yaml
# 编辑 etc/short-api.yaml 文件，填入必要的配置信息
```

4. 启动服务
```bash 
go run short.go 
```

## 主要功能
* 生成可靠短链
* 短链重定向