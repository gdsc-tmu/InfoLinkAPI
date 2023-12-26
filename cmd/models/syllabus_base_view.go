package models

import (
	"time"
)

// [型名] [型についての簡潔な説明]
//
// [より詳細な説明や使用例など]
type SyllabusBaseInfoViewModel struct {
	Year       int16  `json:"year"`
	Season     string `json:"season"`
	Day        string `json:"day"`
	Period     string `json:"period"`
	Teacher    string `json:"teacher"`
	Name       string `json:"name"`
	LectureID  string `json:"lectureId"primary_key:"true"`
	Credits    int16  `json:"credits"`
	URL        string `json:"url"`
	Type       string `json:"type"`
	Faculty    string `json:"faculty"`
}