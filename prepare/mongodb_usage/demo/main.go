package main

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "time"
)

const DbName = `my_db`
const Collection = `my_collection`

// mongoDB 的连接
func main() {
    var (
        ctx        context.Context
        client     *mongo.Client
        db         *mongo.Database
        collection *mongo.Collection
        err        error
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

    collection = collection
}
