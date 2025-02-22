## 认证中心
实现token的分发、校验、删除  

DeliverToken  
对于给定的userid，为它分发一个token  

VerifyToken  
验证一个token是否有效，以及属于哪个用户  
会自动为token续期  

DeleteToken  
删除一个token，即退出该用户在该机器上的登录  

DeleteAllTokens  
删除一个用户的所有token，即退出该用户的所有登录  