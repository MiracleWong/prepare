package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gookit/color.v1"
	"log"
	"os"
	"time"
)

type Task struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Text      string             `bson:"text"`
	Completed bool               `bson:"completed"`
}

var ctx = context.TODO()
var collection *mongo.Collection
func init() {
	var(
		err error
		clientOptions *options.ClientOptions
	)

	// 1. 建立连接
	url := "localhost:27017"
	mongoUrl := "mongodb://" + url
	//mongoUrl := "mongodb://" + user + ":" + password + "@" + url + "/" + dbname
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
	collection = client.Database("tasker").Collection("tasks")

}

func createTask(task *Task) error {
	_, err := collection.InsertOne(ctx, task)
	return err
}

func printTask(tasks []*Task) {
	for i, v := range tasks {
		if v.Completed {
			color.Green.Printf("%d: %s\n", i+1, v.Text)
		} else {
			color.Bold.Printf("%d: %s\n", i+1, v.Text)
		}
	}
}

func filterTasks(filter interface{}) ([]*Task, error)  {
	var tasks []*Task
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return tasks, err
	}

	for cur.Next(ctx) {
		var t Task
		err := cur.Decode(&t)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, &t)
	}

	if err := cur.Err(); err != nil {
		return tasks, err
	}

	cur.Close(ctx)
	if len(tasks) == 0 {
		return tasks, mongo.ErrNoDocuments
	}
	return tasks, nil
}


func getAll() ([]*Task, error)  {
	filter := bson.D{{}}
	return filterTasks(filter)
}

func getPending() ([]*Task, error) {
	filter := bson.D{
		primitive.E{Key: "completed", Value: false},
	}
	return filterTasks(filter)
}

func getFinished() ([]*Task, error) {
	filter := bson.D{
		primitive.E{Key: "completed", Value: true},
	}
	return filterTasks(filter)
}

func completeTask(text string) error {
	filter := bson.D{
		primitive.E{
			Key: "text",
			Value: text,
		}}
	update := bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{
					Key: "completed",
					Value: true,
				},
			},
		},
	}
	t := &Task{}
	return collection.FindOneAndUpdate(ctx, filter, update).Decode(t)
}

func deleteTask(text string) error {
	filter := bson.D{
		primitive.E{
			Key: "text",
			Value: text,
		}}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("No Task were deleted")
	}
	return nil
}

//func main() {
//	err := cmd.Execute()
//	if err != nil {
//		log.Fatalf("cmd.Execute err: %v", err)
//	}
//}

func main() {
	app := &cli.App{
		Name:     "tasker",
		Usage:    "A simple CLI program to manage your tasks",
		Action: func(c *cli.Context) error {
			tasks, err := getPending()
			if err != nil {
				if err == mongo.ErrNoDocuments {
					fmt.Print("Nothing to see here.\n Run `add 'task'` to add a task")
					return nil
				}
				return err
			}
			printTask(tasks)
			return nil
		},
		Commands: []*cli.Command{
			{
				// 1. 创建add 新增的任务
				Name:	"add",
				Aliases: []string{"a"},
				Usage: "add a task to the list",
				Action: func(c *cli.Context) error {
					str := c.Args().First()
					if str == "" {
						return errors.New("Cannnot add an empty task")
					}

					task := &Task{
						ID:        primitive.NewObjectID(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Text:      str,
						Completed: false,
					}

					return createTask(task)
				},
			},
			{
				// 2. 创建all 查询所有的任务
				Name:	"all",
				Aliases: []string{"l"},
				Usage: "list all tasks",
				Action: func(c *cli.Context) error {
					tasks, err := getAll()
					if err != nil{
						if err == mongo.ErrNoDocuments {
							fmt.Println("Nothing to see here")
							return nil
						}
						return nil
					}
					printTask(tasks)
					return nil
				},
			},
			{
				// 3. 创建all 查询所有的任务
				Name:	"done",
				Aliases: []string{"d"},
				Usage: "complete a task on the list",
				Action: func(c *cli.Context) error {
					text := c.Args().First()
					return completeTask(text)
				},
			},
			{
				// 4. finished 查询所有已完成的任务
				Name:	"finished",
				Aliases: []string{"f"},
				Usage: "list completed tasks",
				Action: func(c *cli.Context) error {
					tasks, err := getFinished()
					if err != nil{
						if err == mongo.ErrNoDocuments {
							fmt.Println("Nothing to see here.\\nRun `done 'task'` to complete a task")
							return nil
						}
						return err
					}
					printTask(tasks)
					return nil
				},
			},
			{
				// 5. rm 删除一个任务
				Name:	"rm",
				Usage: "deletes a task on the list",
				Action: func(c *cli.Context) error {
					text := c.Args().First()
					err := deleteTask(text)
					if err != nil{
						return err
					}
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}