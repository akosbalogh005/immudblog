package model

import "time"

type Log struct {

	// unique id (in DB)
	ID int64 `json:"id,omitempty gorm:"primarykey"`

	// logline version
	Version int32 `json:"version,omitempty"`

	// Hostame
	Hostname string `json:"hostname,omitempty" example:"hostname"`

	// name of application
	Application string `json:"application,omitempty" example:"app1" `

	// process id
	Pid string `json:"pid,omitempty"`

	// priority (facility*8 + severity)
	Pri int32 `json:"pri,omitempty"`

	// timestamp of logline (RFC3339)
	Timestamp time.Time `json:"timestamp,omitempty" format:"date-time" gorm:"column:ts"`

	// Message ID
	Messageid int64 `json:"messageid,omitempty"`

	// log message
	Message string `json:"meaasge,omitempty"`
}
