package models

type AgentInfo struct {
	AgentID       string `gorm:"type:varchar(256)"`
	Timestamp     uint64 `gorm:"type:BIGINT UNSIGNED"`
	OsType        string `gorm:"type:varchar(24)"`
	OSVersion     string `gorm:"type:varchar(256)"`
	KernelVersion string `gorm:"type:varchar(256)"`
	Arch          string `gorm:"type:varchar(256)"`
	HostName      string `gorm:"type:varchar(256)"`
	UserName      string `gorm:"type:varchar(256)"`
	Pid           uint32 `gorm:"type:BIGINT UNSIGNED"`
	UptimeSec     uint64 `gorm:"type:BIGINT UNSIGNED"`
	WanIP         string `gorm:"type:varchar(256)"`
	LanIPs        string `gorm:"type:varchar(256)"`
	Mac           string `gorm:"type:varchar(256)"`
	PrimaryIface  string `gorm:"type:varchar(256)"`
}
