package routes

import (
)

type SyllabusResponse struct {
	Year       int16  `json:"year"`
	Season     string `json:"season"`
	Day        string `json:"day"`
	Period     string `json:"period"`
	Teacher    string `json:"teacher"`
	Name       string `json:"name"`
	LectureID  string `json:"lectureId"`
	Credits    int16  `json:"credits"`
	URL        string `json:"url"`
	Type       string `json:"type"`
	Faculty    string `json:"faculty"`
}