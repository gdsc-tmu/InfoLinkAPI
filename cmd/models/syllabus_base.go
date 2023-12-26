package models

import (
	"time"
)

// [型名] [型についての簡潔な説明]
//
// [より詳細な説明や使用例など]
type SyllabusBaseInfo struct {
	Year       int16
	Season     string
	Day        string
	Period     string
	Teacher    string
	Name       string
	LectureId  string
	Credits    int16
	URL        string
	Type       string
	Faculty    string
	DeletedAt *time.Time
}