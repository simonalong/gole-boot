# td-orm
对tdengine-orm进行封装，用于简化配置使用
## 快速使用

### 单库
```yaml
# application.yaml
gole:
   tdengine:
      enable: true
      # 连接类型，支持三种：original, restful, websocket；默认为websocket（v1.5.1 支持original和websocket）
      connect-type: websocket
      host: localhost
      username: root
      password: taosdata
      db-name: td_orm
      # 6030：tcp端口：给original模式使用
      # 6041：http端口：给websocket和restful使用
      port: 6041
```

##### 示例：
```go
// 代码示例
package test

import (
    "github.com/stretchr/testify/assert"
    "github.com/simonalong/gole-boot/extend/tdengine"
    "github.com/simonalong/gole/maps"
    "testing"
    "time"
)

func TestInsert(t *testing.T) {
    tdClient, err := tdengine.GetClient()
    if err != nil {
        panic(err)
        return
    }

    ormMap := maps.Of("ts", time.Now(), "name", "大牛市-boot", "age", 28, "address", "浙江杭州市")
    num, err := tdClient.Insert("td_china", ormMap)
    assert.Equal(t, 1, num)
}
```

# 多库
### 配置文件
```yaml
# application.yaml

gole:
  tdengine:
    enable: true
    # 数据源配置名1
    xxx-name1:
       # 连接类型，支持三种：original, restful, websocket；默认为websocket（v1.5.1 支持original和websocket）
       connect-type: websocket
       host: localhost
       username: root
       password: taosdata
       db-name: td_orm1
       # 6030：tcp端口：给original模式使用
       # 6041：http端口：给websocket和restful使用
       port: 6041
     # 数据源配置名2
    xxx-name2:
       # 连接类型，支持三种：original, restful, websocket；默认为websocket（v1.5.1 支持original和websocket）
       connect-type: websocket
       host: localhost
       username: root
       password: taosdata
       db-name: td_orm2
       # 6030：tcp端口：给original模式使用
       # 6041：http端口：给websocket和restful使用
       port: 6041
```

### 提供api
```go
import(
    "github.com/simonalong/gole-boot/extend/tdengine"
    orm "github.com/simonalong/tdorm"
)

// 获取tdengine单库实例（单例模式获取对象）
tdengine.GetClient() (*orm.TdClient, error){}
// 新建tdengine单库实例（非单例模式获取对象）
tdengine.NewClient() (*orm.TdClient, error){}

// 获取tdengine多库实例（单例模式获取对象）
tdengine.GetClientWithName(dbName string) (*orm.TdClient, error){}
// 新建tdengine多库实例（非单例模式获取对象）
tdengine.NewClientWithName(dbName string) (*orm.TdClient, error){}
```

### 开启日志打印
这边也支持tdorm的sql的打印
```yaml
gole:
   tdengine:
      enable: true
      logger:
         level: debug
```
这个也可以在logger这边开启，等价的
```yaml
gole:
   logger:
      group:
         # 日志分组
         tdorm:
            # 分组级别
            level: debug
```

# 更多用法
更多详细用法请见[tdengine-orm](https://github.com/simonalong/tdorm)，这里截取部分功能

---

我们这里对Orm框架的封装主要提供如下的常见功能
## 用法api
```go
func NewConnectOriginal(host string, port int, user, password, dbName string) *TdClient {}
func NewConnectWebsocket(host string, port int, user, password, dbName string) *TdClient {}
// 暂时不支持restful连接的各种功能
func NewConnectRest(host string, port int, user, password, dbName string) *TdClient {}
```
### 基础api
```go
Exec(sql string, args ...driver.Value) (driver.Result, error)
Query(sql string, args ...driver.Value) ([]*maps.BaseMap, error)
```
### 高级查询
```go
// Select 高级查询：支持所有的查询，包括不限于各种特色查询
Select(hintsDistinctTagsAndColumns *column.Columns, fromClause string, whereConditionAndClause *condition.Condition) ([]*maps.BaseMap, error)
```

### 基础功能封装
对Orm框架的封装主要提供如下的常见功能
1. 新增
   - insert
   - insertWithTag
   - insertEntity
   - insertEntityWithTag
   - insertBatch
   - insertBatchWithTag
   - insertBatchEntity
   - insertBatchEntityWithTag
2. 删除
   - delete
3. 查询：
   - one：查询一行
   - list：查询多行
   - listOfDistinct：查询多行
   - values：查询一列值
   - valuesOfDistinct：查询一列值
   - value：查询一值
   - count：查询个数

### 数据库信息
```go
 SelectDatabase() (string, error)
 SelectClientVersion() (string, error)
 SelectServerVersion() (string, error)
 SelectServerStatus() (string, error) // 服务器状态检测语句。如果服务器正常，返回一个数字（例如 1）。如果服务器异常，返回 error code。该 SQL 语法能兼容连接池对于 TDengine 状态的检查及第三方工具对于数据库服务器状态的检查。并可以避免出现使用了错误的心跳检测 SQL 语句导致的连接池连接丢失的问题
 SelectNow() (string, error)
 SelectToday() (string, error)
 SelectTimeZone() (string, error)
 SelectCurrentUser() (string, error)
 SelectUser() (string, error)
```

### show功能
```go
ShowApps() ([]*maps.BaseMap, error)                                   // 显示接入集群的应用（客户端）信息
ShowCluster() (*maps.BaseMap, error)                                  // 显示当前集群的信息
ShowClusterAlive() (int, error)                                       // 查询当前集群的状态是否可用，返回值： 0：不可用 1：完全可用 2：部分可用（集群中部分节点下线，但其它节点仍可以正常使用）
ShowConnections() ([]*maps.BaseMap, error)                            // 显示当前系统中存在的连接的信息
ShowConsumers() ([]*maps.BaseMap, error)                              // 显示当前数据库下所有消费者的信息
ShowCreateDatabase(dbName string) (string, error)                     // 显示 db_name 指定的数据库的创建语句
ShowCreateStable(dbName, stableName string) (string, error)           // 显示 tb_name 指定的超级表的创建语句
ShowCreateTable(dbName, tableName string) (string, error)             // 显示 tb_name 指定的表的创建语句。支持普通表、超级表和子表
ShowDatabases(dbType string) ([]string, error)                        // dbType：system：指定只显示系统数据库; user：指定只显示用户创建的数据库
ShowDNodes() ([]*maps.BaseMap, error)                                 // 显示当前系统中 DNODE 的信息
ShowFunctions() ([]*maps.BaseMap, error)                              // 显示用户定义的自定义函数。
ShowLicences() ([]*maps.BaseMap, error)                               // 显示企业版许可授权的信息；注：企业版独有
ShowGrants() ([]*maps.BaseMap, error)                                 // 显示企业版许可授权的信息；注：企业版独有
ShowIndexes(dbName, tableName string) ([]*maps.BaseMap, error)        // 显示已创建的索引
ShowLocalVariables() ([]*maps.BaseMap, error)                         // 显示当前客户端配置参数的运行值。
ShowMNodes() ([]*maps.BaseMap, error)                                 // 显示当前系统中 MNODE 的信息
ShowQNodes() ([]*maps.BaseMap, error)                                 // 显示当前系统中 QNODE （查询节点）的信息
ShowQueries() ([]*maps.BaseMap, error)                                // 显示当前系统中正在进行的查询
ShowScores() ([]*maps.BaseMap, error)                                 // 显示系统被许可授权的容量的信息 注：企业版独有
ShowStables(dbName, stableNamePart string) ([]string, error)          // 显示当前数据库下的所有超级表的信息。会自动使用 LIKE 对表名进行模糊匹配
ShowStreams() ([]*maps.BaseMap, error)                                // 显示当前系统内所有流计算的信息
ShowSubscriptions() ([]*maps.BaseMap, error)                          // 显示当前系统内所有的订阅关系
ShowTables(tableRole, dbName, tableNamePart string) ([]string, error) // 显示当前数据库下的所有普通表和子表的信息。会自动使用 LIKE 对表名进行模糊匹配。tableRole：NORMAL 指定只显示普通表信息， CHILD 指定只显示子表信息
ShowTags(dbName, tableName string) ([]*maps.BaseMap, error)           // 显示子表的标签信息
ShowTopics() ([]string, error)                                        // 显示当前数据库下的所有主题的信息
ShowTransactions() ([]*maps.BaseMap, error)                           // 显示当前系统中正在执行的事务的信息(该事务仅针对除普通表以外的元数据级别)
ShowUsers() ([]*maps.BaseMap, error)                                  // 显示当前系统中所有用户的信息。包括用户自定义的用户和系统默认用户
ShowClusterVariables() ([]*maps.BaseMap, error)                       // 显示当前系统中各节点需要相同的配置参数的运行值，也可以指定 DNODE 来查看其的配置参数
ShowVGroups(dbName string) ([]*maps.BaseMap, error)                   // 显示当前数据库中所有 VGROUP 的信息
ShowVNodes(dNodeId string) ([]*maps.BaseMap, error)                   // 显示当前系统中所有 VNODE 或某个 DNODE 的 VNODE 的信息
```

### 用法示例
```go
func TestTdOrm(t *testing.T) {
    tdOrm, err := tdengine.NewTdengine()
    if err != nil {
        panic(err)
        return
    }
    
    // 新增：使用map，对应SQL：【insert into td_orm.td_china (`ts`,`name`,`age`,`address`) values (?,?,?,?)】
    insertMap := maps.NewSort().Put("ts", time.Now()).Put("name", "大牛市").Put("age", "18").Put("address", "浙江杭州市")
    tdOrm.Insert("td_china", insertMap)

    // 新增：使用entity，对应SQL：【insert into td_orm.td_china (`ts`,`name`,`age`,`address`) values (?,?,?,?)】
    type OrmChinaDomain struct {
    	Timestamp time.Time `column:"ts"`
    	Na        string    `column:"name"`
    	Ag        int       `column:"age"`
    	Add       string    `column:"address"`
    }
    tdChinaDomain := OrmChinaDomain{Timestamp: time.Now(), Na: "大牛市2", Ag: 19, Add: "浙江温州市"}
    tdOrm.InsertEntity("td_china", tdChinaDomain)

    // 删除，对应SQL：【delete from td_orm.td_china where `ts` > ?】
    tdOrm.Delete("td_china", condition.New().Gt("ts", "2024-07-12 12:00:00.000"))

    // 查询：一行，对应SQL：【select `name`, `age` from td_orm.td_china where `ts` < ? limit 1】
    tdOrm.One("td_china", column.Of("name", "age"), condition.New().Lt("ts", time.Now()))

    // 查询：多行，对应SQL：【select `name`, `age` from td_orm.td_china where `age` = ?】
    tdOrm.List("td_china", column.Of("name", "age"), condition.New().Eq("age", 18))

    // 查询：一个，对应SQL：【select `name` from td_orm.td_china where `age` > ? and `ts` = ? limit 1】
    tdOrm.Value("td_china", "name", condition.New().Gt("age", 12).Eq("ts", "2024-07-16 11:19:23.291"))

    // 查询：一列，对应SQL：【select `name` from td_orm.td_china where `age` = ?】
    tdOrm.Values("td_china", "name", condition.New().Eq("age", 18))

    // 查询：个数，对应SQL：【select count(*) as cnt from td_orm.td_china where `age` = ?】
    tdOrm.Count("td_china", condition.New().Eq("age", 18))

    // 新增：批量新增，对应SQL：【insert into td_orm.td_china (`ts`,`name`,`age`,`address`) values (?,?,?,?)】
    timeData1, _ := baseTime.ParseTimeYmdHmsS("2024-07-15 14:01:23.391")
    timeData2, _ := baseTime.ParseTimeYmdHmsS("2024-07-15 15:01:23.391")
    var insertMaps []*maps.OrmMap
    insertMaps = append(insertMaps, maps.NewSort().Put("ts", timeData1).Put("name", "大牛市1").Put("age", "18").Put("address", "浙江杭州市1"))
    insertMaps = append(insertMaps, maps.NewSort().Put("ts", timeData2).Put("name", "大牛市2").Put("age", "28").Put("address", "浙江杭州市2"))
    tdOrm.BatchInsert("td_china", insertMaps)
    // 也支持批量插入实体
    //tdOrm.BatchInsertEntity("td_china", insertEntities)
}	
```

## 参数column用法
### column用法
支持功能：
- hints    *Hint     // 提示，就是tdengine的提示功能
- fields   []string  // 属性
- distinct bool      // 是否展示distinct
- tags     bool      // 是否展示tags


```go
// 普通用法
assert.Equal(t, "`id`, `name`, `age`", column.Of("id", "name", "age").ToSql())

// 全部列：通配符
assert.Equal(t, "*", column.Of("*").ToSql())
assert.Equal(t, "tableName.*", column.Of("tableName.*").ToSql())

// 别名
assert.Equal(t, "id as dataId, `name`, `age`", column.Of("id as dataId", "name", "age").ToSql())

// 标签列
assert.Equal(t, "tags id, `name`, `age`", column.Of("tags id", "name", "age").ToSql())
assert.Equal(t, "tags `id`, `name`, `age`", column.Tags().Of("id", "name", "age").ToSql())

// 去重
assert.Equal(t, "distinct `id`, `name`, `age`", column.Distinct().Of("id", "name", "age").ToSql())

// Hints 提示
assert.Equal(t, "/*+ BATCH_SCAN() */ `ts`", column.Hints(&column.Hint{BatchScan: true}).Of("ts").ToSql())
assert.Equal(t, "/*+ NO_BATCH_SCAN() */ `ts`", column.Hints(&column.Hint{NoBatchScan: true}).Of("ts").ToSql())
assert.Equal(t, "/*+ SORT_FOR_GROUP() */ `ts`", column.Hints(&column.Hint{SortForGroup: true}).Of("ts").ToSql())
assert.Equal(t, "/*+ PARTITION_FIRST() */ `ts`", column.Hints(&column.Hint{PartitionFirst: true}).Of("ts").ToSql())
assert.Equal(t, "/*+ PARA_TABLES_SORT() */ `ts`", column.Hints(&column.Hint{ParaTablesSort: true}).Of("ts").ToSql())
assert.Equal(t, "/*+ SMALLDATA_TS_SORT() */ `ts`", column.Hints(&column.Hint{SmalldataTsSort: true}).Of("ts").ToSql())
assert.Equal(t, "/*+ SKIP_TSMA() */ `ts`", column.Hints(&column.Hint{SkipTsma: true}).Of("ts").ToSql())

// 多个提示
assert.Equal(t, "/*+ BATCH_SCAN() SKIP_TSMA() */ `ts`", column.Hints(&column.Hint{BatchScan: true, SkipTsma: true}).Of("ts").ToSql())
```


## 查询条件
查询条件主要是condition支持功能，也支持BaseMap作为条件，这里主要讲condition.Condition的功能；支持如下函数
```go
And(datas ...any) *Condition
AndEm(datas ...any) *Condition
Or(datas ...any) *Condition
OrEm(datas ...any) *Condition
Em(datas ...any) *Condition
Gt(key string, value interface{}) *Condition
Ge(key string, value interface{}) *Condition
Lt(key string, value interface{}) *Condition
Le(key string, value interface{}) *Condition
Eq(key string, value interface{}) *Condition
UnEq(key string, value interface{}) *Condition
IsNull(key string) *Condition
IsNotNull(key string) *Condition
BetweenAnd(key string, leftValue, rightValue interface{}) *Condition
NotBetweenAnd(key string, leftValue, rightValue interface{}) *Condition
In(key string, values interface{}) *Condition
NotIn(key string, values interface{}) *Condition
Like(key string, value interface{}) *Condition
NotLike(key string, value interface{}) *Condition
Match(key string, value interface{}) *Condition
NMatch(key string, value interface{}) *Condition
NotMatch(key string, value interface{}) *Condition
Contains(key string, value interface{}) *Condition

// 这个是tdengine的特色查询，这里提供直接拼接功能
Append(otherClause string) *Condition 
```
query是查询条件
```go
assert.Equal(t, "where (`age` = ? and `name` = ? and `address` = ?)", condition.New().And("age", 12, "name", "zhou", "address", "杭州").ToSql())
assert.Equal(t, "where `age` = ? and `name` = ? and `address` = ?", condition.New().AndEm("age", 12, "name", "zhou", "address", "杭州").ToSql())

assert.Equal(t, "where (`age` = ?) and (`name` = ?) and (`address` = ?)", condition.New().And("age", 12).And("name", "zhou").And("address", "杭州").ToSql())
assert.Equal(t, "where `age` = ? and `name` = ? and `address` = ?", condition.New().AndEm("age", 12).AndEm("name", "zhou").AndEm("address", "杭州").ToSql())

assert.Equal(t, "where (`age` = ?)", condition.New().And("age", 12).ToSql())
assert.Equal(t, "where `age` = ?", condition.New().AndEm("age", 12).ToSql())

// value数据为nil，则sql不拼接
assert.Equal(t, "where `age` = ? and `name` = ?", condition.New().AndEm("age", 12, "name", "zhou", "address", nil).ToSql())
assert.Equal(t, "where `age` = ? and `name` = ? and `address` = ?", condition.New().AndEm("age", 12, "name", "zhou", "address", "").ToSql())
assert.Equal(t, "", condition.New().AndEm("age", nil).ToSql())
```
一个条件简单查询，可以使用简化如下
```go
// 比较操作符：两种写法，最后生成的sql是一样的
assert.Equal(t, "where `age` > ?", condition.New().Gt("age", 12).ToSql())
assert.Equal(t, "where `age` > ?", condition.New().AndEm(op.Gt("age", 12)).ToSql())

assert.Equal(t, "where `age` > ? and `name` = ?", condition.New().Gt("age", 12).AndEm("name", "zhou").ToSql())
assert.Equal(t, "where `age` > ? and `name` = ?", condition.New().AndEm(op.Gt("age", 12), "name", "zhou").ToSql())

// 其他比较操作符
assert.Equal(t, "where `age` >= ?", condition.New().Ge("age", 12).ToSql())
assert.Equal(t, "where `age` < ?", condition.New().Lt("age", 12).ToSql())
assert.Equal(t, "where `age` <= ?", condition.New().Le("age", 12).ToSql())
assert.Equal(t, "where `age` = ?", condition.New().Eq("age", 12).ToSql())
assert.Equal(t, "where `age` != ?", condition.New().UnEq("age", 12).ToSql())

// is null；is not null
assert.Equal(t, "where `name` is null", condition.New().IsNull("name").ToSql())
assert.Equal(t, "where `name` is not null", condition.New().IsNotNull("name").ToSql())

// between and
assert.Equal(t, "where `age` between ? and ?", condition.New().BetweenAnd("age", 12, 20).ToSql())
assert.Equal(t, "where `age` not between ? and ?", condition.New().NotBetweenAnd("age", 12, 20).ToSql())

// in；not in
assert.Equal(t, "where `age` in (?, ?)", condition.New().In("age", []int{12, 18}).ToSql())
assert.Equal(t, "where `age` in (?, ?)", condition.New().In("age", []int{12, 18}).ToSql())
assert.Equal(t, "where `age` in (?, ?, ?, ?, ?, ?, ?, ?, ?)", condition.New().In("age", []int{12, 18, 43, 32, 43, 54, 65, 12, 64}).ToSql())
assert.Equal(t, "where `age` not in (?, ?, ?)", condition.New().NotIn("age", []int{12, 18, 119}).ToSql())

// like; not like
assert.Equal(t, "where `name` like ?", condition.New().Like("name", "%牛_").ToSql())
assert.Equal(t, "where `name` like ?", condition.New().Like("name", "%牛__").ToSql())
assert.Equal(t, "where `name` like ?", condition.New().Like("name", "%牛%").ToSql())
assert.Equal(t, "where `name` not like ?", condition.New().NotLike("name", "%牛%").ToSql())

// match; nmatch
assert.Equal(t, "where `name` match ?", condition.New().Match("name", "%牛_").ToSql())
assert.Equal(t, "where `name` match ?", condition.New().Match("name", "%牛__").ToSql())
assert.Equal(t, "where `name` match ?", condition.New().Match("name", "%牛%").ToSql())
assert.Equal(t, "where `name` nmatch ?", condition.New().NotMatch("name", "%牛%").ToSql())
assert.Equal(t, "where `name` nmatch ?", condition.New().NMatch("name", "%牛%").ToSql())

// contains
assert.Equal(t, "where `info` contains ?", condition.New().Contains("info", "k1").ToSql())
```
也支持嵌套`op.xxx`；op操作可以作为condition.And()或者Or()的参数，如下示例
```go
// 使用And 和op.Or混用
assert.Equal(t, "where ((`age` = ?) and `name` = ? and `address` = ?)", condition.New().And(op.Or("age", 12), "name", "zhou", "address", "杭州").ToSql())
assert.Equal(t, "where ((`age` = ? or `name` = ?) and `address` = ?)", condition.New().And(op.Or("age", 12, "name", "zhou"), "address", "杭州").ToSql())
assert.Equal(t, "where ((`age` = ? or `name` = ? or `address` = ?))", condition.New().And(op.Or("age", 12, "name", "zhou", "address", "杭州")).ToSql())
```
op的函数有如下，与condition.And或者Or函数可以组成更加复杂的条件
```go
// 比较操作符
assert.Equal(t, "where (`age` > ? and `name` = ? and `address` = ?)", condition.New().And(op.Gt("age", 12), "name", "zhou", "address", "杭州").ToSql())
assert.Equal(t, "where (`age` >= ? and `name` = ? and `address` = ?)", condition.New().And(op.Ge("age", 12), "name", "zhou", "address", "杭州").ToSql())
assert.Equal(t, "where (`age` < ? and `name` = ? and `address` = ?)", condition.New().And(op.Lt("age", 12), "name", "zhou", "address", "杭州").ToSql())
assert.Equal(t, "where (`age` <= ? and `name` = ? and `address` = ?)", condition.New().And(op.Le("age", 12), "name", "zhou", "address", "杭州").ToSql())
//assert.Equal(t, "where (`age` <> ? and `name` = ? and `address` = ?)", condition.New().And(op.UnEq("age", 12), "name", "zhou", "address", "杭州").ToSql())
assert.Equal(t, "where (`age` != ? and `name` = ? and `address` = ?)", condition.New().And(op.UnEq("age", 12), "name", "zhou", "address", "杭州").ToSql())

// is null；is not null
assert.Equal(t, "where `name` is null and `address` = ?", condition.New().AndEm(op.IsNull("name"), "address", "杭州").ToSql())
assert.Equal(t, "where `name` is not null and `address` = ?", condition.New().AndEm(op.IsNotNull("name"), "address", "杭州").ToSql())

// between and
assert.Equal(t, "where `age` between ? and ? and `address` = ?", condition.New().AndEm(op.BetweenAnd("age", 12, 20), "address", "杭州").ToSql())
assert.Equal(t, "where `age` not between ? and ? and `address` = ?", condition.New().AndEm(op.NotBetweenAnd("age", 12, 20), "address", "杭州").ToSql())

// in；not in
assert.Equal(t, "where `age` in (?, ?) and `address` = ?", condition.New().AndEm(op.In("age", []int{12, 18}), "address", "杭州").ToSql())
assert.Equal(t, "where `age` in (?, ?) and `address` = ?", condition.New().AndEm(op.In("age", []int{12, 18}), "address", "杭州").ToSql())
assert.Equal(t, "where `age` in (?, ?, ?, ?, ?, ?, ?, ?, ?) and `address` = ?", condition.New().AndEm(op.In("age", []int{12, 18, 43, 32, 43, 54, 65, 12, 64}), "address", "杭州").ToSql())
assert.Equal(t, "where `age` not in (?, ?, ?) and `address` = ?", condition.New().AndEm(op.NotIn("age", []int{12, 18, 119}), "address", "杭州").ToSql())

// like; not like
assert.Equal(t, "where `name` like ?", condition.New().Em(op.Like("name", "%牛_")).ToSql())
assert.Equal(t, "where `name` like ?", condition.New().Em(op.Like("name", "%牛__")).ToSql())
assert.Equal(t, "where `name` like ?", condition.New().Em(op.Like("name", "%牛%")).ToSql())
assert.Equal(t, "where `name` not like ?", condition.New().Em(op.NotLike("name", "%牛%")).ToSql())

// match; nmatch
assert.Equal(t, "where `name` match ?", condition.New().Em(op.Match("name", "%牛_")).ToSql())
assert.Equal(t, "where `name` match ?", condition.New().Em(op.Match("name", "%牛__")).ToSql())
assert.Equal(t, "where `name` match ?", condition.New().Em(op.Match("name", "%牛%")).ToSql())
assert.Equal(t, "where `name` nmatch ?", condition.New().Em(op.NotMatch("name", "%牛%")).ToSql())
assert.Equal(t, "where `name` nmatch ?", condition.New().Em(op.NMatch("name", "%牛%")).ToSql())

// contains
assert.Equal(t, "where `info` contains ?", condition.New().Em(op.Contains("info", "k1")).ToSql())
```

高级查询
```go
// tdengine的高级功能
assert.Equal(t, "where `status` = ?", condition.New().Eq("status", 2).ToSql())
assert.Equal(t, "where `status` = ? limit 10", condition.New().Eq("status", 2).Append("limit 10").ToSql())
assert.Equal(t, "where `status` = ? partition by location interval(10m)", condition.New().Eq("status", 2).Append("partition by location interval(10m)").ToSql())
assert.Equal(t, "where `status` = ? partition by tbname state_window(case when voltage >= 205 and voltage <= 235 then 1 else 0 end)", condition.New().Eq("status", 2).Append("partition by tbname state_window(case when voltage >= 205 and voltage <= 235 then 1 else 0 end)").ToSql())
assert.Equal(t, "session(ts, tol_val)", condition.New().Append("session(ts, tol_val)").ToSql())
```
