package service

import (
	"immudblog/immudb"
	"immudblog/model"

	log "github.com/sirupsen/logrus"
)

// service package
// This layer is between API call and the persists layer. When more action is needed within a service
// we call them here. Now it seems unnecessary since we don't have complex services...

// AddLog service for adding new log
func AddLog(l model.Log) error {
	log.Infof("AddLog service. Param: %v", l)
	return immudb.ImmuDB.AddLog(l)
}

// AddLogs service for adding new logs
func AddLogs(logs []model.Log) error {
	log.Infof("AddLogs service. Param: %v", logs)
	return immudb.ImmuDB.AddLogs(logs)
}

// CountLogs service counts lines in Logs table
func CountLogs() (uint64, error) {
	log.Infof("CountLogs service")
	return immudb.ImmuDB.CountLogs()
}

// GetLogs service return  max. cou items from table Log. Ordered by ID desc
func GetLogs(cou uint64, app string) (ret []model.Log, err error) {
	log.Infof("GetLogs service. Param:%v, %v", cou, app)
	return immudb.ImmuDB.GetLogs(cou, app)
}
