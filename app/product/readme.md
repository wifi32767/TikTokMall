## 商品服务
实现商品信息的增删改查  

## api

### CreateProduct  
创建一个商品  
```
input:
{
    name string
    description string
    picture string
    price string
    categories []string
}
return:
{
    id uint32
}
```

### UpdateProduct  
更新商品信息  
```
input:
{
    id uint32
    name string
    description string
    picture string
    price float
    categories []string
}
return:
{
    success bool
}
```

### DeleteProduct  
删除商品  
```
input:
{
    id uint32
}
return:
{
    success bool
}
```

### ListProducts  
列出一页商品列表  
```
input:
{
    page int32
    pageSize int32
    categoryName string
}
return:
{
    products []Product
}
```

### GetProduct  
获取单个商品信息  
```
input:
{
    id uint32
}
return:
{
    product Product
}
```

### SearchProducts  
搜索名字或描述中含有指定字符串的商品  
```
input:
{
    query string
}
return:
{
    results []Product
}
```