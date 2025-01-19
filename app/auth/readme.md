## 认证中心
实现token的分发、校验、删除  

## api
```
DeliverTokenByRPC
input: 
{
    user_id int32
}
return: 
{
    token string
}

VerifyTokenByRPC
input: 
{
    token string
}
return: 
{
    res bool
    user_id int32
}
```