package master

import (
	"context"
	"github.com/gothicrush/crush-scheduler/common"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"time"
)

type LogManager struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	//单例
	G_logManager *LogManager
)

func InitLogManager() error {

	client, err := mongo.Connect(context.TODO(), G_config.MongodbUri,
		clientopt.ConnectTimeout(time.Duration(G_config.MongodbConnectTimeout)*time.Millisecond))

	if err != nil {
		return err
	}

	G_logManager = &LogManager{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
	}

	return nil
}

func (logManager *LogManager) ListLog(name string, skip int, limit int) ([]*common.JobLog, error) {

	// 过滤条件
	filter := common.JobLogFilter{
		JobName: name,
	}

	// 排序条件
	logSort := &common.SortLogByStartTime{
		StartOrder: -1,
	}

	// 游标
	cursor, err := logManager.logCollection.Find(context.TODO(), filter,
		findopt.Sort(logSort), findopt.Skip(int64(skip)), findopt.Limit(int64(limit)))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var logArr []*common.JobLog = make([]*common.JobLog, 0)

	var jobLog *common.JobLog
	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}

		if err := cursor.Decode(jobLog); err != nil {
			continue
		}

		logArr = append(logArr, jobLog)
	}

	return logArr, nil
}
