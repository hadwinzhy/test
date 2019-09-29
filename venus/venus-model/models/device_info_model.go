package models

const (
	STATUS_NETWORK_OFFLINE     = "network_offline"
	STATUS_UPGRADE_START       = "upgrade_start"
	STATUS_UPGRADE_DOWNLOADING = "upgrade_downloading"
	STATUS_UPGRADE_INSTALLING  = "upgrade_installing"
	STATUS_UPGRADE_RESTARTING  = "upgrade_restarting"
	STATUS_NETWORK_PAIRING     = "network_pairing"
	STATUS_NETWORK_ONLINE      = "network_online"
	STATUS_XCLOUD_ONLINE       = "xcloud_online"
	STATUS_XCLOUD_OFFLINE      = "xlcoud_offline"
)

// DeviceInfo ...
type DeviceInfo struct {
	BaseModel
	DeviceID      uint   `gorm:"not null;index;unique_index" json:"device_id"`
	XcloudStatus  string `gorm:"type:varchar(50);not null" json:"xcloud_status"`
	NetworkName   string `gorm:"type:varchar(50);not null" json:"network_name"`
	NetwordStatus string `gorm:"type:varchar(50);not null" json:"network_status"`
	TmateKey      string `gorm:"type:varchar(50)" json:"tmate_key"`
	IpAddr        string `gorm:"type:varchar(50)" json:"ip_addr"`
	// UpgradeAt     time.Time `gorm:"type:timestamp with time zone" json:"upgrade_at"`
	// UpgradeVersionCode int       `json:"upgrade_version_code"`
}
