## 用户服务
实现用户的注册、登录、删除、修改密码  
用户登出由后端实现  

## api

Register  
注册一个用户  
```
input:
{
    username string
    password string
}
return:
{
    userid uint32
}
```

Login  
用户登录  
```
input:
{
    username string
    password string
}
return:
{
    userid uint32
}
```

Delete  
删除一个用户  
```
input:
{
    username string
    password string
}
return:
{
    success bool
}
```

Update  
更新一个用户的密码  
```
input:
{
    username string
    old_password string
    ner_password string
}
return:
{
    success bool
}
```