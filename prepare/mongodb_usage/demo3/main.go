package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// mongoDB批量写入多条
func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		record     *LogRecord
		logArr     []interface{}
		result     *mongo.InsertManyResult
		insertId   interface{}
		docId      string
	)

	// 创建连接
	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://192.168.205.10:27017")); err != nil {
		fmt.Printf("Failed to connect mongo, err: %s", err.Error())
		return
	}

	// 选择数据库
	database = client.Database(DbName)

	// 选择集
	collection = database.Collection(Collection)

	// 插入记录
	record = &LogRecord{
		JobName: "job1",
		Command: "echo hello",
		Error:   "",
		Content: "hello",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	logArr = []interface{}{record, record, record} // 批量插入多条document

	// 发起写入
	if result, err = collection.InsertMany(context.TODO(), logArr); err != nil {
		fmt.Printf("Failed to insert mongo, err: %s", err.Error())
		return
	}

	for _, insertId = range result.InsertedIDs {
		docId = insertId.(primitive.ObjectID).Hex()
		fmt.Printf("自增ID：%s \n", docId)
	}
}
