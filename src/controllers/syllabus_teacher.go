package controllers

import (
	"InfoLinkAPI/src/models"
	"net/http"
	"unicode"
	"fmt"
	"github.com/gin-gonic/gin"
)
// rune が空白文字であれば true を返却
func IsWhiteSpace(r rune) bool {
	return unicode.IsSpace(r)
}
// GetSyllabusByFaculty 指定した教員名が部分一致するシラバスを返す。
//
// 引数: 教員名 e.g. 山口
// 戻り値: teacherフィールドが引数に部分一致したシラバス
// documentation: 
func (sc *SyllabusController) GetSyllabusByTeacher(c *gin.Context) {
	var syllabus []models.SyllabusBaseInfo
	teacherName := c.Param("name")
	// `都立太郎` -> `%都%立%太%郎%`
	queryTeacherName := "%"
	for _, r := range teacherName {
		if !IsWhiteSpace(r) {
			queryTeacherName += fmt.Sprintf("%s%%", string(r))
		}
	}
	result := sc.DB.Where("teacher LIKE ?", queryTeacherName).Find(&syllabus)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	
	c.JSON(http.StatusOK, syllabus)
}