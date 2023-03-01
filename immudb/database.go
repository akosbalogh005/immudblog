package immudb

import (
	"context"
	"database/sql"
	"fmt"
	"immudblog/config"
	"immudblog/model"

	"github.com/codenotary/immudb/pkg/client"

	immugorm "github.com/codenotary/immugorm"
	gormlogrus "github.com/onrik/gorm-logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ImmuDBStruct struct {
	DB *gorm.DB
	ImmuDBIF
}

//go:generate mockery --name=ImmuDBIF
type ImmuDBIF interface {
	AddLog(l model.Log) error
	AddLogs(logs []model.Log) error
	CountLogs() (uint64, error)
	GetLogs(cou uint64, app string) (ret []model.Log, err error)
}

var ImmuDB ImmuDBIF
var immuDBImpl ImmuDBStruct

func Init() {
	log.Infof("Initialize Immudb...")
	opts := client.DefaultOptions()

	opts.Username = config.ImmudbFlags.Username
	opts.Password = config.ImmudbFlags.Password
	opts.Database = config.ImmudbFlags.Database
	opts.Address = config.ImmudbFlags.Host
	opts.Port = config.ImmudbFlags.Port
	opts.HealthCheckRetries = 10

	db, err := gorm.Open(immugorm.OpenWithOptions(opts, &immugorm.ImmuGormConfig{Verify: false}), &gorm.Config{
		Logger: gormlogrus.New(),
	})

	if err != nil {
		log.Fatalf("Cannot initialize IMMUDB! Error: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&model.Log{})
	if err != nil {
		log.Fatalf("Cannot AutoMigrate Log table in IMMUDB! Error: %v", err)
	}

	immuDBImpl.DB = db
	ImmuDB = immuDBImpl

}

// AddLog adds new log into database
func (d ImmuDBStruct) AddLog(l model.Log) error {
	var x []model.Log
	x = append(x, l)
	return d.AddLogs(x)
}

// AddLog adds new logs into database in one transaction
func (d ImmuDBStruct) AddLogs(logs []model.Log) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {

		for i := range logs {

			log.Debugf("Adding : %v", logs[i])
			err := tx.Create(&logs[i]).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Count Rows in table Logs
func (d ImmuDBStruct) CountLogs() (uint64, error) {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return 0, err
	}
	var cou uint64
	err = sqlDB.QueryRow("SELECT COUNT(*) AS COU FROM Logs").Scan(&cou)
	if err != nil {
		return 0, err
	}
	return cou, nil
}

// GetLogs get max. cou items from table Log. Ordered by ID desc
func (d ImmuDBStruct) GetLogs(cou uint64, app string) (ret []model.Log, err error) {
	var sqlDB *sql.DB
	sqlDB, err = d.DB.DB()
	if err != nil {
		return
	}
	whereFilter := ""
	if app != "" {
		// FIXME: SQL Injection
		whereFilter = fmt.Sprintf(" WHERE APPLICATION = '%s'", app)
	}
	rows, err := sqlDB.QueryContext(
		context.TODO(),
		fmt.Sprintf("SELECT ID, VERSION, HOSTNAME, APPLICATION, PID, PRI, TS, MESSAGEID, MESSAGE FROM Logs %s ORDER BY id DESC", whereFilter),
	)
	if err != nil {
		return
	}
	defer rows.Close()
	max := cou
	for rows.Next() && max > 0 {
		var l model.Log
		err = rows.Scan(&l.ID, &l.Version, &l.Hostname, &l.Application, &l.Pid, &l.Pri, &l.Timestamp, &l.Messageid, &l.Message)
		if err != nil {
			return
		}
		ret = append(ret, l)
		max--
	}
	return
}
