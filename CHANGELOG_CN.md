# Change Logs

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

[TOC]

> Unrelease

## [变更逻辑]

1. 注册成功后，不会自动创建AA钱包，需要调用新增接口 [POST /api/passkey/v1/account/chain](#post-apipasskeyv1accountchain)
2. 同一个邮件，在不同设备允许再次注册（本质是注册该账户的Passkey），注册成功后依然是原用户
3. 登录不再需要输入邮箱，现在会自动发现当前设备在该origin下的Passkeys

## [变更接口]

### GET /api/passkey/v1/account/info

1. query增加参数：alias, 默认为空

### POST /api/passkey/v1/reg/verify

1. query增加参数：alias, 默认为空

### POST /api/passkey/v1/sign

1. body移除email，只保留origin

### POST /api/passkey/v1/sign/verify

1. query参数只保留origin

### POST /api/passkey/v1/tx/sign

1. body参数中nonce改名为ticket
2. body参数中增加network，表示链名称，必填
3. body参数中增加network_alias，表示链别名，非必填参数


### POST /api/passkey/v1/tx/sign/verify

1. query参数中nonce改名为ticket
2. query参数中增加network，表示链名称，必填
3. query参数中增加network_alias，表示链别名，非必填参数
4. 响应结果增加BLS信息，包括bls_sign，bls_pubkey和bls_schema，其中bls_sign和bls_pubkey使用`base64url`编码，分别表示bls的签名和公钥，bls_schema表示bls算法信息

响应结果示例

```json
{
    "code": 200,
    "message": "",
    "data": {
        "code": 200,
        "txdata": "48656c6c6f2c20576f726c6421",
        "sign": "0xce6e001af297172aa5176b2a50200148002e9a9ece9293694fe5374b453f62c30e445a94a4b6ebab0c05631614c934b4eb318ffd4dd5ff159807430aeaff32e51c",
        "bls_sign": "jclHo_HBKZUzYKMgJsCqpjXxGkIrJZZXx2npXHFHHdgNA4fu0VD5wsrLlVUB04OKCFNfiKw7XzGOT4Ob7qhkYkFjXDXxg3l0mG9VYct7KI7EtYd6H3jgGE93j6lYfwaM",
        "bls_pubkey": "rzpZJX8PiAXgsf4BQZJ6hbJkhQUWyM0lrch_YmrqzZ34Z_YeqXa4nKc-lnkm7vqR",
        "address": "0xf75291198A0c549962Db2D4816321e95a0e48fc3",
        "bls_schemal": "BLS12_381:EthModeDraft07"
    },
    "cost": "2562047h47m16.854775807s"
}
```

## [新增接口]

### GET /api/passkey/v1/chains/support

> 获取支持的链名称

入参：无

*data*对象里的key表示链名称，在其他接口参数中的network需符合该键名，value为true时表示支持，*data*中没有列出的链表示不支持

出参：

```json
{
  "code": 200,
  "message": "",
  "data": {
    "base-sepolia": true,
    "ethereum-sepolia": true,
    "optimism-sepolia": true
  },
  "cost": "2562047h47m16.854775807s"
}
```

### POST /api/passkey/v1/account/chain

> 创建指定链的AA钱包

入参：

body:

```json
{
  "alias": "string",
  "network": "ethereum-mainnet"
}
```

alias: 别名，允许为空，同一个链支持不同别名的钱包，目前上限为10个