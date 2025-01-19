## 认证中心
实现token的分发、校验、删除  

## api
```
DeliverToken
对于给定的userid，为它分发一个token
input: 
{
    user_id int32
}
return: 
{
    token string
}

VerifyToken
验证一个token是否有效，以及属于哪个用户
会自动为token续期
input: 
{
    token string
}
return: 
{
    res bool
    user_id int32
}

DeleteToken
删除一个token，即退出该用户在该机器上的登录
input:
{
    token string
}
return: empty

DeleteAllTokens
删除一个用户的所有token，即退出该用户的所有登录
input:
{
    user_id int32
}
return: empty
```