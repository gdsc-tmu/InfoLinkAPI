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

func (syllabus *SyllabusBaseInfo) ToSyllabusViewModel() SyllabusViewModel {
	return SyllabusViewModel{
        Year:   syllabus.Year,
        Season: syllabus.Season,
        Day:    syllabus.Day,
		Period: syllabus.Period,
        Teacher: syllabus.Teacher,
		Name:   syllabus.Name,
		LectureId: syllabus.LectureId,
		Credits: syllabus.Credits,
		URL: syllabus.URL,
		Type: syllabus.Type,
		Faculty: syllabus.Faculty,
    }
}