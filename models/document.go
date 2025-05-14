package models

import "time"

type Document struct {
	ID        string `gorm:"primaryKey;size:36"`
	FileName  string `gorm:"size:255"`
	FilePath  string `gorm:"size:255"`
	Signature string `gorm:"type:text"`
	FileType  string `gorm:"size:20"`
	CreatedAt time.Time
}
