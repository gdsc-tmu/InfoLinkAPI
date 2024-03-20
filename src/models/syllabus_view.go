package models

// [型名] [型についての簡潔な説明]
//
// [より詳細な説明や使用例など]
type SyllabusViewModel struct {
	Year       int16  `json:"year"`
	Season     string `json:"season"`
	Day        string `json:"day"`
	Period     string `json:"period"`
	Teacher    string `json:"teacher"`
	Name       string `json:"name"`
	LectureId  string `json:"lectureId"`
	Credits    int16  `json:"credits"`
	URL        string `json:"url"`
	Type       string `json:"type"`
	Faculty    string `json:"faculty"`
}

func GetSyllabusViewModelBySyllabusBaseInfo(syllabusBaseInfo SyllabusBaseInfo) SyllabusViewModel {
	return SyllabusViewModel{
		Year:       syllabusBaseInfo.Year,
		Season:     syllabusBaseInfo.Season,
		Day:        syllabusBaseInfo.Day,
		Period:     syllabusBaseInfo.Period,
		Teacher:    syllabusBaseInfo.Teacher,
		Name:       syllabusBaseInfo.Name,
		LectureId:  syllabusBaseInfo.LectureId,
		Credits:    syllabusBaseInfo.Credits,
		URL:        syllabusBaseInfo.URL,
		Type:       syllabusBaseInfo.Type,
		Faculty:    syllabusBaseInfo.Faculty,
	}
}