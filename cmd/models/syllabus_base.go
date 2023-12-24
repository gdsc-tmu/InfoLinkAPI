package models

import (
	"time"
)

// [型名] [型についての簡潔な説明]
//
// [より詳細な説明や使用例など]
type SyllabusBaseInfo struct {
	LectureID  string `gorm:"primaryKey;size:5"`
	Year       int16  `gorm:"type:smallint"`
	Season     string `gorm:"size:8"`
	Day        string `gorm:"size:30"`
	Period     string `gorm:"size:30"`
	Teacher    string `gorm:"size:50"`
	Name       string `gorm:"size:100"`
	Credits    int16  `gorm:"type:smallint"`
	URL        string `gorm:"size:100"`
	Type       string `gorm:"size:20"`
	Faculty    string `gorm:"size:4"`
	DeletedAt  *time.Time
}