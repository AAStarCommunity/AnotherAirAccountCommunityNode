[![codecov](https://codecov.io/gh/AAStarCommunity/AnotherAirAccountCommunityNode/graph/badge.svg?token=G741C0D6SR)](https://codecov.io/gh/AAStarCommunity/AnotherAirAccountCommunityNode)

# AnotherAirAccountCommunityNode

A decentration community node for AirAccount

## What's Community Node

> TBD

## Quick Start

### 1. Swagger

#### 1.1 install

```shell
go get -u github.com/swaggo/swag
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