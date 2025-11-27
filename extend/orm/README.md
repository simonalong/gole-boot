# orm
对业内的常见Orm进行封装，进行方便使用，目前支持的有
- gorm
- xorm
  
注意：<br/>
xorm: 这个是xorm.io不是go-xorm，go-xorm暂时不支持

### 单数据源
#### 代码
```go
import "github.com/simonalong/gole-boot/extend/orm"

/// ------------------------ gorm ------------------------
// 获取gorm客户端（单例模式获取实例）
func GetGormClient() (*gorm.DB, error) {}
// 新建gorm客户端（非单例模式获取实例）
func NewGormClient() (*gorm.DB, error) {}

// gorm：获取默认配置库实例，自定义配置
func NewGormClientWitConfig(gormConfig *gorm.Config) (*gorm.DB, error) {}

/// ------------------------ xorm ------------------------
// xorm：获取默认配置库实例（单例模式获取实例）
func GetXormClient()(*xorm.Engine, error){}
// xorm：获取默认配置库实例（非单例模式获取实例）
func NewXormClient()(*xorm.Engine, error){}

// xorm：获取默认配置库实例，自定义参数
func NewXormClientWithParams(params map[string]string)(*xorm.Engine, error){}
```
#### 全部配置
```yaml
gole:
  datasource:
    # 是否启用，默认关闭
    enable: true
    username: user
    password: passwd
    host: 10.33.33.33
    port: 8080
    # 目前支持: mysql、postgresql、sqlserver
    driver-name: mysql
    # 数据库名
    db-name: xx_db
    # 示例：charset=utf8&parseTime=True&loc=Local 等url后面的配置，直接配置即可
    url-config:
      xxx: xxxxx
      yyy: yyyyy
    # 连接池配置
    connect-pool:
      # 最大空闲连接数
      max-idle-conns: 10
      # 最大连接数
      max-open-conns: 1000
      # 连接可重用最大时间；带字符（s：秒，m：分钟，h：小时）
      max-life-time: 10s
      # 连接空闲的最大时间；带字符（s：秒，m：分钟，h：小时）
      max-idle-time: 10s
```

### 多数据源
#### 代码
```go
import "github.com/simonalong/gole-boot/extend/orm"
/// ------------------------ gorm ------------------------
// 获取gorm客户端，根据数据源配置名（单例模式获取实例）
func GetGormClientWithName(datasourceName string) (*gorm.DB, error) {}
// 新建gorm客户端，根据数据源配置名（非单例模式获取实例）
func NewGormClientWithName(datasourceName string) (*gorm.DB, error) {}

// gorm：根据数据源配置名获取库实例，自定义配置
func NewGormDbWithNameAndConfig(gormConfig *gorm.Config) (*gorm.DB, error) {}

/// ------------------------ xorm ------------------------
// xorm：获取默认配置库实例，根据数据源配置名（单例模式获取实例）
func GetXormClientWithName(datasourceName string)(*xorm.Engine, error){}
// xorm：获取默认配置库实例，根据数据源配置名（非单例模式获取实例）
func NewXormClientWithName(datasourceName string)(*xorm.Engine, error){}

// xorm：根据数据源配置名获取库实例，自定义参数
func NewXormDbWithNameParams(datasourceName string, params map[string]string)

// xorm：主从接口
func NewXormDbMasterSlave(masterDatasourceName string, slaveDatasourceNames []string, policies ...xorm.GroupPolicy)
```
#### 配置
```yml
gole:
  datasource:
    # 是否启用，默认关闭
    enable: true
    # 数据源配置名1
    xxx-name1:
      username: xxx
      password: xxx
      host: xxx
      port: xxx
      # 目前支持: mysql、postgresql、sqlserver
      driver-name: xxx
      # 数据库名
      db-name: xx_db
      # 示例：charset=utf8&parseTime=True&loc=Local 等url后面的配置，直接配置即可
      url-config:
        xxx: xxx
        yyy: yyy
      # 连接池配置
      connect-pool:
        # 最大空闲连接数
        max-idle-conns: 10
        # 最大连接数
        max-open-conns: 10
        # 连接可重用最大时间；带字符（s：秒，m：分钟，h：小时）
        max-life-time: 10s
        # 连接空闲的最大时间；带字符（s：秒，m：分钟，h：小时）
        max-idle-time: 10s
    # 数据源配置名2
    xxx-name2:
      username: xxx
      password: xxx
      host: xxx
      port: xxx
      # 目前支持: mysql、postgresql、sqlserver
      driver-name: xxx
      # 数据库名
      db-name: xx_db
      # 示例：charset=utf8&parseTime=True&loc=Local 等url后面的配置，直接配置即可
      url-config:
        xxx: xxxxx
        yyy: yyyyy
      # 连接池配置
      connect-pool:
        # 最大空闲连接数
        max-idle-conns: 10
        # 最大连接数
        max-open-conns: 10
        # 连接可重用最大时间；带字符（s：秒，m：分钟，h：小时）
        max-life-time: 10s
        # 连接空闲的最大时间；带字符（s：秒，m：分钟，h：小时）
        max-idle-time: 10s
```
### 示例：gorm
```go
func TestGorm1(t *testing.T) {
    client, _ := orm.GetGormClient()

    // 删除表
	client.Exec("drop table test.base_demo")

    //新增表
	client.Exec("CREATE TABLE base_demo(\n" +
        "  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',\n" +
        "  `name` char(20) NOT NULL COMMENT '名字',\n" +
        "  `age` INT NOT NULL COMMENT '年龄',\n" +
        "  `address` char(20) NOT NULL COMMENT '名字',\n" +
        "  \n" +
        "  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n" +
        "  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n" +
        "\n" +
        "  PRIMARY KEY (`id`)\n" +
    ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='测试表'")

    // 新增
    client.Create(&BaseDemo{Name: "zhou", Age: 18, Address: "杭州"})
    client.Create(&BaseDemo{Name: "zhou", Age: 11, Address: "杭州2"})
    
    // 查询：一行
    var demo BaseDemo
    client.First(&demo).Where("name=?", "zhou")
    
    fmt.Println(demo)
}
```

### 示例：xorm
单数据源
```go
func TestXorm1(t *testing.T) {
    db, _ := orm.GetXormClient()
    
    // 删除表
    db.Exec("drop table test.base_demo")
    
    //新增表
    db.Exec("CREATE TABLE base_demo(\n" +
        "  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',\n" +
        "  `name` char(20) NOT NULL COMMENT '名字',\n" +
        "  `age` INT NOT NULL COMMENT '年龄',\n" +
        "  `address` char(20) NOT NULL COMMENT '名字',\n" +
        "  \n" +
        "  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n" +
        "  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n" +
        "\n" +
        "  PRIMARY KEY (`id`)\n" +
    ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='测试表'")
    
    db.Table("base_demo").Insert(&BaseDemo{Name: "zhou", Age: 18, Address: "杭州"})
    // 新增
    db.Table("base_demo").Insert(&BaseDemo{Name: "zhou", Age: 18, Address: "杭州"})
    
    var demo BaseDemo
    db.Table("base_demo").Where("name=?", "zhou").Get(&demo)
    
    fmt.Println(demo)
}
```

### 注意
请不要在业务中使用init方法获取db，因为这个时候config的配置还没有加载出来

## 框架配置
上面全都是数据库的配置，对于一些orm框架本身也会有一些配置，这里支持下(version >= 1.5.2)
支持配置：
- 打印sql
```yaml
gole:
  orm:
    show-sql: true
```
#### 线上动态开启和关闭sql的话
```yaml
curl -X PUT http://localhost:xxx/{api-prefix}/config/update -d '{"key":"baes.orm.show-sql", "value":"true"}'
```
或者如下
```yaml
curl -X PUT http://localhost:xxx/{api-prefix}/config/update -d '{"key":"baes.logger.group.orm.level", "value":"debug"}'
```
这两个配置功能是等同的，一个是直接基于logger来修改，一个是基于orm的配置来修改



### 更多配置
在一些场景下，也需要mysql本身提供一些配置，就是最近遇到gorm默认在mariadb下面是报失败，因此增加了这样的配置（version >= 1.5.2）
```yaml
gole:
  datasource:
    mysql:
      server-version: ""
      skip-initialize-with-version: false
      default-string-size: 0
      disable-with-returning: false
      disable-datetime-precision: false
      dont-support-rename-index: false
      dont-support-rename-column: false
      dont-support-for-share-clause: false
      dont-support-null-as-default-value: false
```
以上这些配置其实对应的是如下的代码，示例
```go
// 其中的`DisableWithReturning` 对应的就是上面的 gole.datasource.mysql.disable-with-returning，其他更多的配置都在里面
gorm.Open(mysql.New(mysql.Config{Conn: conn, DisableWithReturning: true}))
```
