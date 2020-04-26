## MongoDB

### Docker 启动

```
sudo docker run -d --name mongo -p 27017:27017 mongo
```

### 连接MongoDB操作

```
sudo docker exec -it mongo /bin/bash -c "mongo admin"
```

### 简单操作

#### 数据库操作

```
show databases # 列举数据库

use DB_NAME # 选择数据库，数据库无需事先创建，知识一个命名空间
```

#### 创建collection

```
show collections # 列举数据表
db.createCollection("MY_COLLECTION") # 创建数据表，无需定义字段
```

#### 文档document操作

- 新增
```
db.MY_COLLECTION.insertOne({uid: 1, name: "curder", hobbies: ["music", "code"]}) # 文档ID通常不需要自己指定；任意嵌套层级的BSON
```

- 查询
```
db.MY_COLLECTION.find({hobbies: 'music', name: {$in: ['curder', 'xiaoming']}}).sort({uid: 1}) # 可以基于任意BSON层级过滤；支持的功能与MySQL相当
```

- 更新
```
db.MY_COLLECTION.updateMany({hobbies: 'music'}, {$set: {name: "xiaoming"}}) # 第一个参数过滤条件，第二个参数是要更新的字段
```

- 删除
```
db.MY_COLLECTION.deleteMany({name: "curder"}) # 参数是过滤条件
```

- 创建索引index
```
db.MY_COLLECTION.createIndex({uid: 1, name: -1}) # 可以指定建立索引的正反序
```


### MySQL对比

#### 概念类比
| MySQL | MongoDB |
| ---- | ---- |
| database | database |
| table | collection |
| row | document(bson) |
| column | field |
| index | index |
| table joins | $lookup |
| primary key | _id |
| group by | aggregation pipeline |


#### 聚合类比

| MySQL | MongoDB |
| ---- | ---- |
| where | $match |
| group by | $group |
| having | $match |
| select | $project |
| order by | $sort |
| limit | $limit |
| sum | $sum |
| count | $sum |
