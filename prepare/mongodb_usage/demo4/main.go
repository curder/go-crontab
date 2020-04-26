package main

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const DbName = `cron`
const Collection = `logs`

// 日志结构体
type LogRecord struct {
    JobName   string    `bson:"jobName"`   // 任务名称
    Command   string    `bson:"command"`   // 任务名称
    Error     string    `bson:"error"`     // 脚本错误
    Content   string    `bson:"content"`   // 脚本输出
    TimePoint TimePoint `bson:"timePoint"` // 执行时间点
}

// 任务执行时间点
type TimePoint struct {
    StartTime int64 `bson:"startTime"` // 开始时间
    EndTime   int64 `bson:"endTime"`   // 结束时间
}

// jobName过滤条件
type FindByJobName struct {
    JobName string `bson:"jobName"`
}

// mongoDB根据条件读取数据
func main() {
    var (
        client     *mongo.Client
        err        error
        database   *mongo.Database
        collection *mongo.Collection
        condition  *FindByJobName
        cursor     *mongo.Cursor
        skip       int64
        limit      int64
        record     *LogRecord
    )

    // 建立连接
    if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://192.168.205.10:27017")); err != nil {
        fmt.Printf("Failed to connect mongo, err: %s", err.Error())
        return
    }

    // 选择数据库
    database = client.Database(DbName)

    // 选择集
    collection = database.Collection(Collection)

    // 根据条件查找
    // 按照jobName过滤，查找jobName = job1的document
    condition = &FindByJobName{JobName: "job1"}
    limit = 2
    skip = 0
    // 查询（过滤+翻页参数）
    if cursor, err = collection.Find(context.TODO(), condition, &options.FindOptions{Limit: &limit, Skip: &skip}); err != nil {
        fmt.Printf("Failed to find err: %s", err.Error())
        return
    }

    // 延迟关闭游标
    defer cursor.Close(context.TODO())

    // 遍历结果集
    for cursor.Next(context.TODO()) {
        // 初始化一个空结构体变量
        record = &LogRecord{}
        // 反序列化bson对象到结构体
        if err = cursor.Decode(record); err != nil {
            fmt.Printf("Decode cursor err: %s", err.Error())
            return
        }

        fmt.Printf("%#v\n", *record)
    }
}
