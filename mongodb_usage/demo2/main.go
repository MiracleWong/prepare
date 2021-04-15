package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// 声明一个MongoConn，存储MongoDB连接相关的信息

type MongoConn struct {
	clientOptions *options.ClientOptions
	client        *mongo.Client
	collections   *mongo.Collection
}

// 执行的任务时间点

type TimePoint struct {
	StartTime int64 `bson:"StartTime"`
	EndTime int64 `bson:"endTime"`
}

// 1条日志

type LogRecord struct {
	JobName string `bson:"jobName"`
	Command string `bson:"command"`
	Err string	`bson:"err"`
	Content string `bson:"content"`
	TimePoint TimePoint `bson:"timePoint"`
}



var ctx = context.TODO()
var err error
var clientOptions *options.ClientOptions
var collection *mongo.Collection
var insertOneResult *mongo.InsertOneResult
var docId primitive.ObjectID

//func InitMongoConn(url, dbname string) error {
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	// 1. 建立连接
//	// construct url: mongodb://username:password@127.0.0.1:27017/dbname
//	url = "localhost:27017"
//	mongoUrl := "mongodb://" + url
//	mongoConn.clientOptions = options.Client().ApplyURI(mongoUrl)
//
//	// Connect to MongoDB
//	var err error
//	mongoConn.client, err = mongo.Connect(ctx, mongoConn.clientOptions)
//	if err != nil {
//		println(err)
//	}
//
//	// 2. 测试连接
//	err = mongoConn.client.Ping(context.TODO(), nil)
//	if err != nil {
//		println(err)
//	}
//	// 2. 选择数据库和表
//	mongoConn.collections = mongoConn.client.Database(dbname).Collection("tests")
//
//	return err
//}

//func CloseMongoConn() {
//	err := mongoConn.client.Disconnect(context.TODO())
//	if err != nil {
//		println("disconnect mongo connect is error: %v", err)
//		return
//	}
//	println("connection to MongoDB closed.")
//}

func main() {

	// 1. 建立连接
	url := "localhost:27017"
	mongoUrl := "mongodb://" + url
	clientOptions = options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 2. Ping 方法
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 选择数据库和Collection
	collection = client.Database("cron").Collection("log")

	// 4. 插入记录
	record := &LogRecord{
		JobName: "job10",
		Command: "echo Hello",
		Err: "",
		Content: "Hello",
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix()+10},
	}
	if insertOneResult,err = collection.InsertOne(ctx,record); err != nil {
		println(err)
	}

	fmt.Println("InsetedID: ", insertOneResult.InsertedID)
	docId = insertOneResult.InsertedID.(primitive.ObjectID)
	fmt.Println("自增ID：",docId.Hex())
}