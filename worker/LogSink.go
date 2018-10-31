package worker

import (
	"context"
	"github.com/gothicrush/crush-scheduler/common"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

// MongoDB存储
type LogSink struct {
	client         *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}

var (
	// 存储单例
	G_logSink *LogSink
)

func InitLogSink() error {

	// 与MongoDB进行连接
	client, err := mongo.Connect(context.TODO(), G_config.MongodbUri,
		clientopt.ConnectTimeout(time.Duration(G_config.MongodbConnectTimeout)*time.Millisecond))

	if err != nil {
		return err
	}

	// 选择db和collection
	G_logSink = &LogSink{
		client:         client,
		logCollection:  client.Database("cron").Collection("log"),
		logChan:        make(chan *common.JobLog, 1000),
		autoCommitChan: make(chan *common.LogBatch, 1000),
	}

	go G_logSink.writeLoop()

	return nil
}

// 日志存储协程
func (logSink *LogSink) writeLoop() {

	var logBatch *common.LogBatch // 当前批次
	var commitTimer *time.Timer
	var timeoutBatch *common.LogBatch // 超时批次

	for {
		select {
		case log := <-logSink.logChan:
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				commitTimer = time.AfterFunc(
					time.Duration(G_config.JobLogCommitTimeout)*time.Millisecond,
					func(batch *common.LogBatch) func() {
						return func() {
							logSink.autoCommitChan <- batch
						}
					}(logBatch),
				)
			}

			//把新的日志追加到批次中去
			logBatch.Logs = append(logBatch.Logs, log)

			if len(logBatch.Logs) >= G_config.JobLogBatchSize {
				// 批量写入日志
				logSink.saveLogs(logBatch)
				// 情况logBatch
				logBatch = nil
				// 取消定时器
				commitTimer.Stop()
			}
		case timeoutBatch = <-logSink.autoCommitChan:
			// 判断过期批次是否还是当前的批次
			if timeoutBatch != logBatch {
				continue // 跳过已经提交的批次
			}
			// 将该批次写入mongodb中
			logSink.saveLogs(timeoutBatch)
			//情况logBatch
			logBatch = nil
		}
	}
}

// 批量写入日志
func (logSink *LogSink) saveLogs(batch *common.LogBatch) {
	logSink.logCollection.InsertMany(context.TODO(), batch.Logs)
}

// 发送日志
func (logSink *LogSink) Append(jobLog *common.JobLog) {
	select {
	case logSink.logChan <- jobLog:
	default: //队列满了就丢弃
	}

}
