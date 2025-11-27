
# 所有接入的指标列表

### orm（gorm）
```docker
# HELP base_boot_gorm_db_idle_connections 空闲连接数
# TYPE base_boot_gorm_db_idle_connections gauge
base_boot_gorm_db_idle_connections{db="test"} 1

# HELP base_boot_gorm_db_in_use_connections 当前正在使用的连接数
# TYPE base_boot_gorm_db_in_use_connections gauge
base_boot_gorm_db_in_use_connections{db="test"} 0

# HELP base_boot_gorm_db_max_open_connections 数据库的最大打开连接数
# TYPE base_boot_gorm_db_max_open_connections gauge
base_boot_gorm_db_max_open_connections{db="test"} 10

# HELP base_boot_gorm_db_open_connections 正在使用和空闲的已建立连接数
# TYPE base_boot_gorm_db_open_connections gauge
base_boot_gorm_db_open_connections{db="test"} 1

# HELP base_boot_gorm_max_idle_closed_connections 由于SetMaxIdleCon而关闭的连接总数
# TYPE base_boot_gorm_max_idle_closed_connections gauge
base_boot_gorm_max_idle_closed_connections{db="test"} 0

# HELP base_boot_gorm_max_idle_time_closed_connections 由于SetConnMaxIdleTime而关闭的连接总数
# TYPE base_boot_gorm_max_idle_time_closed_connections gauge
base_boot_gorm_max_idle_time_closed_connections{db="test"} 0

# HELP base_boot_gorm_max_lifetime_closed_connections 由于SetConnMaxLifetime而关闭的连接总数
# TYPE base_boot_gorm_max_lifetime_closed_connections gauge
base_boot_gorm_max_lifetime_closed_connections{db="test"} 0

# HELP base_boot_gorm_wait_count_connections 等待的连接总数
# TYPE base_boot_gorm_wait_count_connections gauge
base_boot_gorm_wait_count_connections{db="test"} 0

# HELP base_boot_gorm_wait_duration_connections 等待新连接的总阻塞时间
# TYPE base_boot_gorm_wait_duration_connections gauge
base_boot_gorm_wait_duration_connections{db="test"} 0
```
### gin
```docker
# HELP gin_request_body_total the server received request body size, unit byte
# TYPE gin_request_body_total counter
gin_request_body_total 0

# HELP gin_request_duration the time server took to handle the request.
# TYPE gin_request_duration histogram
gin_request_duration_bucket{uri="/api/orm/query",le="0.1"} 5
gin_request_duration_bucket{uri="/api/orm/query",le="0.3"} 5
gin_request_duration_bucket{uri="/api/orm/query",le="1.2"} 5
gin_request_duration_bucket{uri="/api/orm/query",le="5"} 5
gin_request_duration_bucket{uri="/api/orm/query",le="10"} 5
gin_request_duration_bucket{uri="/api/orm/query",le="+Inf"} 5
gin_request_duration_sum{uri="/api/orm/query"} 0.008236192
gin_request_duration_count{uri="/api/orm/query"} 5

# HELP gin_request_total all the server received request num.
# TYPE gin_request_total counter
gin_request_total 5

# HELP gin_request_uv_total all the server received ip num.
# TYPE gin_request_uv_total counter
gin_request_uv_total 1

# HELP gin_response_body_total the server send response body size, unit byte
# TYPE gin_response_body_total counter
gin_response_body_total 5220

# HELP gin_uri_request_total all the server received request num with every uri.
# TYPE gin_uri_request_total counter
gin_uri_request_total{code="200",method="GET",uri="/api/orm/query"} 5
```
### tcp
#### tcp 接收端
```docker
# HELP base_boot_tcp_connections_open tcp的连接数
# TYPE base_boot_tcp_connections_open gauge
base_boot_tcp_connections_open 0

# HELP base_boot_tcp_request_counter tcp的请求总数
# TYPE base_boot_tcp_request_counter counter
base_boot_tcp_request_counter 1

# HELP base_boot_tcp_transfer_rate_bytes TCP传输速率
# TYPE base_boot_tcp_transfer_rate_bytes histogram
base_boot_tcp_transfer_rate_bytes_bucket{le="1024"} 0
base_boot_tcp_transfer_rate_bytes_bucket{le="2048"} 0
base_boot_tcp_transfer_rate_bytes_bucket{le="4096"} 0
base_boot_tcp_transfer_rate_bytes_bucket{le="16384"} 1
base_boot_tcp_transfer_rate_bytes_bucket{le="65536"} 1
base_boot_tcp_transfer_rate_bytes_bucket{le="262144"} 1
base_boot_tcp_transfer_rate_bytes_bucket{le="1.048576e+06"} 1
base_boot_tcp_transfer_rate_bytes_bucket{le="+Inf"} 1
base_boot_tcp_transfer_rate_bytes_sum 5724
base_boot_tcp_transfer_rate_bytes_count 1
```
### grpc
#### grpc 接收端
```docker
# HELP grpc_server_handled_total Total number of RPCs completed on the server, regardless of success or failure.
# TYPE grpc_server_handled_total counter
grpc_server_handled_total{grpc_code="Aborted",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 4
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 0

# HELP grpc_server_msg_received_total Total number of RPC stream messages received on the server.
# TYPE grpc_server_msg_received_total counter
grpc_server_msg_received_total{grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 4

# HELP grpc_server_msg_sent_total Total number of gRPC stream messages sent by the server.
# TYPE grpc_server_msg_sent_total counter
grpc_server_msg_sent_total{grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 4

# HELP grpc_server_started_total Total number of RPCs started on the server.
# TYPE grpc_server_started_total counter
grpc_server_started_total{grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 4
```
#### grpc 发起端
```docker
# HELP grpc_client_handled_total Total number of RPCs completed by the client, regardless of success or failure.
# TYPE grpc_client_handled_total counter
grpc_client_handled_total{grpc_code="OK",grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 4

# HELP grpc_client_msg_sent_total Total number of gRPC stream messages sent by the client.
# TYPE grpc_client_msg_sent_total counter
grpc_client_msg_sent_total{grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 4

# HELP grpc_client_started_total Total number of RPCs started on the client.
# TYPE grpc_client_started_total counter
grpc_client_started_total{grpc_method="SayHello",grpc_service="demo.Greeter",grpc_type="unary"} 4
```

### http
```docker
# HELP base_boot_http_client_requests_counter 向外部发起的http请求总个数
# TYPE base_boot_http_client_requests_counter counter
base_boot_http_client_requests_counter{method="GET",status_code="200",url="/api/ok"} 1
base_boot_http_client_requests_counter{method="GET",status_code="500",url="/api/err"} 4
```
### nats
#### nats jetstream 发送端
```docker
# HELP base_boot_nats_js_client_request_ok_counter nats jetstream 发布消息成功总数
# TYPE base_boot_nats_js_client_request_ok_counter counter
base_boot_nats_js_client_request_ok_counter{method="PublishMsg",reply="",subject="test.js.demo.req"} 9

# HELP base_boot_nats_js_client_request_err_counter nats jetstream 发布消息异常总数
# TYPE base_boot_nats_js_client_request_err_counter counter
base_boot_nats_js_client_request_err_counter{method="PublishMsg",reply="",subject="test.js.demo.req"} 1
```
#### nats jetstream 接收端
```docker
# HELP base_boot_nats_js_server_handle_ok_counter nats jetstream 消息成功处理总数
# TYPE base_boot_nats_js_server_handle_ok_counter counter
base_boot_nats_js_server_handle_ok_counter{subject="test.js.demo.req"} 9

# HELP base_boot_nats_js_server_handle_err_counter nats jetstream 消息异常处理总数
# TYPE base_boot_nats_js_server_handle_err_counter counter
base_boot_nats_js_server_handle_err_counter{subject="test.js.demo.req"} 1
```

#### nats 非jetstream 发送端
```docker
# HELP base_boot_nats_client_request_ok_counter nats发布消息成功总数
# TYPE base_boot_nats_client_request_ok_counter counter
base_boot_nats_client_request_ok_counter{method="PublishMsg",reply="",subject="test.js.demo.req"} 9

# HELP base_boot_nats_client_request_err_counter nats发布消息异常总数
# TYPE base_boot_nats_client_request_err_counter counter
base_boot_nats_client_request_err_counter{method="PublishMsg",reply="",subject="test.js.demo.req"} 1
```
#### nats 非jetstream 接收端
```docker
# HELP base_boot_nats_server_handle_ok_counter nats消息成功处理总数
# TYPE base_boot_nats_server_handle_ok_counter counter
base_boot_nats_server_handle_ok_counter{subject="test.js.demo.req"} 9

# HELP base_boot_nats_server_handle_err_counter nats消息异常处理总数
# TYPE base_boot_nats_server_handle_err_counter counter
base_boot_nats_server_handle_err_counter{subject="test.js.demo.req"} 1
```
### redis
```docker
# HELP base_boot_redis_single_commands Histogram of single Redis commands
# TYPE base_boot_redis_single_commands histogram
base_boot_redis_single_commands_bucket{command="get",instance="redis-demo",le="0.001"} 11
base_boot_redis_single_commands_bucket{command="get",instance="redis-demo",le="0.005"} 12
base_boot_redis_single_commands_bucket{command="get",instance="redis-demo",le="0.01"} 12
base_boot_redis_single_commands_bucket{command="get",instance="redis-demo",le="+Inf"} 12
base_boot_redis_single_commands_sum{command="get",instance="redis-demo"} 0.0033382660000000008
base_boot_redis_single_commands_count{command="get",instance="redis-demo"} 12
base_boot_redis_single_commands_bucket{command="set",instance="redis-demo",le="0.001"} 4
base_boot_redis_single_commands_bucket{command="set",instance="redis-demo",le="0.005"} 4
base_boot_redis_single_commands_bucket{command="set",instance="redis-demo",le="0.01"} 4
base_boot_redis_single_commands_bucket{command="set",instance="redis-demo",le="+Inf"} 4
base_boot_redis_single_commands_sum{command="set",instance="redis-demo"} 0.000637669
base_boot_redis_single_commands_count{command="set",instance="redis-demo"} 4
```
### tdengine
```docker
# HELP base_boot_tdengine_request_err_counter tdengine处理异常的总数
# TYPE base_boot_tdengine_request_err_counter counter
base_boot_tdengine_request_err_counter{RunType="insert",db="td_orm"} 8

# HELP base_boot_tdengine_request_ok_counter tdengine处理正常的总数
# TYPE base_boot_tdengine_request_ok_counter counter
base_boot_tdengine_request_ok_counter{RunType="insert",db="td_orm"} 13

# HELP base_boot_tdengine_slow_sql_histogram tdengine慢查询的统计（单位秒）
# TYPE base_boot_tdengine_slow_sql_histogram histogram
base_boot_tdengine_slow_sql_histogram_bucket{db="td_orm",sql="insert",le="1"} 13
base_boot_tdengine_slow_sql_histogram_bucket{db="td_orm",sql="insert",le="5"} 13
base_boot_tdengine_slow_sql_histogram_bucket{db="td_orm",sql="insert",le="10"} 13
base_boot_tdengine_slow_sql_histogram_bucket{db="td_orm",sql="insert",le="30"} 13
base_boot_tdengine_slow_sql_histogram_bucket{db="td_orm",sql="insert",le="60"} 13
base_boot_tdengine_slow_sql_histogram_bucket{db="td_orm",sql="insert",le="300"} 13
base_boot_tdengine_slow_sql_histogram_bucket{db="td_orm",sql="insert",le="1800"} 13
base_boot_tdengine_slow_sql_histogram_bucket{db="td_orm",sql="insert",le="+Inf"} 13
base_boot_tdengine_slow_sql_histogram_sum{db="td_orm",sql="insert"} 0.012560939000000002
base_boot_tdengine_slow_sql_histogram_count{db="td_orm",sql="insert"} 13
```
