package models

type VersionApp struct {
	Id          int    `json:"id"`
	VersionCode int    `json:"version_code"`
	VersionName string `json:"version_name"`
	Type        int    `json:"type"`
}
