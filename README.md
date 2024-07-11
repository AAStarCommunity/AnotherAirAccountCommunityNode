[![codecov](https://codecov.io/gh/AAStarCommunity/AnotherAirAccountCommunityNode/graph/badge.svg?token=G741C0D6SR)](https://codecov.io/gh/AAStarCommunity/AnotherAirAccountCommunityNode)

# AnotherAirAccountCommunityNode

A decentration community node for AirAccount

## What's Community Node

> TBD

## Prepare

### Pgsql

```shell
docker run --name community_node -e POSTGRES_PASSWORD=mypassword -d -p 5432:5432 postgres
```

## Quick Start

### 1. Swagger

#### 1.1 install

```shell
go get -u github.com/swaggo/swag
go install github.com/swaggo/swag/cmd/swag@latest
```

#### 1.2 init swag

```shell
swag init -g ./main.go
```

> FAQ: [Unknown LeftDelim and RightDelim in swag.Spec](https://github.com/swaggo/swag/issues/1568)

### 2. Run

```shell
go mod tidy
go run ./cmd/server/main.go
```

### Docker启动

```shell
docker-compose -f ./example/one-click-deploy/docker-compose.yaml up -d
```
