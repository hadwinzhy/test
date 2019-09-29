package models

import (
	"fmt"
)

// Version ...
type Version struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Latest        bool   `json:"latest"`
	UpdateMessage string `json:"update_message"`
	VersionCode   int64  `gorm:"index;" json:"version_code"`
	VersionType   string `gorm:"type: varchar(30);" json:"version_type"`
	VersionTime   string `json:"version_time"`
	URL           string `json:"url"`
	Md5           string `json:"md5"`
	Lot           string `gorm:"index;" json:"lot"`
}

// VersionSerializer ...
type VersionSerializer struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Latest         bool   `json:"latest"`
	UpdateMessage  string `json:"update_message"`
	VersionCode    int64  `json:"version_code"`
	VersionCodeStr string `json:"version_code_str"`
	VersionTime    string `json:"version_time"`
	VersionType    string `json:"version_type"`
	URL            string `json:"url"`
	Md5            string `json:"md5"`
	Lot            string `json:"lot"`
}

// VersionSerialize ...
func (v *Version) VersionSerialize() VersionSerializer {
	return VersionSerializer{
		ID:             v.ID,
		Name:           v.Name,
		Latest:         v.Latest,
		UpdateMessage:  v.UpdateMessage,
		VersionCode:    v.VersionCode,
		VersionCodeStr: fmt.Sprintf("%d", v.VersionCode),
		VersionType:    v.VersionType,
		VersionTime:    v.VersionTime,
		URL:            v.URL,
		Md5:            v.Md5,
		Lot:            v.Lot,
	}
}
