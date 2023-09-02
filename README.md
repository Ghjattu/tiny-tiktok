# Tiny-Tiktok

## Development Environment
- OS: Ubuntu 20.04
- Golang: 1.20.7
- MySQL: 8.0.33
- Redis: 7.0.12
- RabbitMQ: 3.8.2

## Get Started
### Download Dependencies
```go
go mod download
```

### Environment Variables
在项目根目录下创建 `.env` 文件，变量列表如下：
#### 阿里云OSS
- `OSS_BUCKET_NAME`（必需）: 阿里云 OSS 的 Bucket 名称
- `OSS_ENDPOINT`（必需）: 阿里云 OSS 的 Endpoint(地域节点)
- `OSS_ACCESS_KEY_ID`（必需）: 阿里云的 AccessKey ID
- `OSS_ACCESS_KEY_SECRET`（必需）: 阿里云的 AccessKey Secret
#### MySQL
- `MYSQL_USERNAME`（可选）: MySQL 的用户名，默认为 `root`
- `MYSQL_PASSWORD`（可选）: MySQL 的密码，默认为空
- `MYSQL_IP`（可选）: MySQL 的 IP，默认为 `127.0.0.1`
- `MYSQL_PORT`（可选）: MySQL 的端口，默认为 `3306`
- `MYSQL_DB_NAME`（可选）: MySQL 的数据库名，默认为 `tiktok`
- `MYSQL_DB_NAME_TEST`（可选）: MySQL 的测试数据库名，默认为 `tiktok_test`
#### JSON Web Token
- `SECRET_KEY`（可选）: JWT 签名的密钥，默认为 `secret-key`
- `TOKEN_LIFESPAN`（可选）: 整数，表示 JWT 的有效期，单位为小时，默认为 `24`
#### Redis
- `REDIS_IP`（可选）: Redis 的 IP，默认为 `127.0.0.1`
- `REDIS_PORT`（可选）: Redis 的端口，默认为 `6379`
- `REDIS_PASSWORD`（可选）: Redis 的密码，默认为空
#### RabbitMQ
- `RABBITMQ_USERNAME`（可选）: RabbitMQ 的用户名，默认为 `guest`
- `RABBITMQ_PASSWORD`（可选）: RabbitMQ 的密码，默认为 `guest`
- `RABBITMQ_IP`（可选）: RabbitMQ 的 IP，默认为 `127.0.0.1`
- `RABBITMQ_PORT`（可选）: RabbitMQ 的端口，默认为 `5672`

### Run
```go
go run main.go
```

## Architecture
![architecture](./images/architecture.png)