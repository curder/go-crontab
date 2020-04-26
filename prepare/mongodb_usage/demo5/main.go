package main

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "time"
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

type TimeBeforeCondition struct {
    Before int64 `bson:"$lt"`
}

type DeleteCondition struct {
    BeforeCondition TimeBeforeCondition `bson:"timePoint.startTime"`
}

// mongoDB根据条件删除操作
func main() {
    var (
        ctx             context.Context
        client          *mongo.Client
        db              *mongo.Database
        collection      *mongo.Collection
        err             error
        deleteCondition *DeleteCondition
        deleteResult    *mongo.DeleteResult
    )
    // 建立连接
    ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
    if client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://192.168.205.10:27017")); err != nil {
        fmt.Printf("Failed to connect mongodb, err: %s", err.Error())
        return
    }

    // 选择数据库
    db = client.Database(DbName)

    // 选择表
    collection = db.Collection(Collection)

    // 删除开始时间早于当前时间的所有日志
    // 构建删除条件
    deleteCondition = &DeleteCondition{
        BeforeCondition: TimeBeforeCondition{Before: time.Now().Unix()},
    }

    // 删除操作
    if deleteResult, err = collection.DeleteMany(context.TODO(), deleteCondition); err != nil {
        fmt.Printf("Failed to delete many row: %s", err.Error())
        return
    }

    fmt.Printf("删除了%d行", deleteResult.DeletedCount)
}
