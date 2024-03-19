package controllers

import (
	"InfoLinkAPI/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSyllabusByFaculty 指定した教員名が部分一致するシラバスを返す。
//
// 引数: 教員名 e.g. 山口
// 戻り値: teacherフィールドが引数に部分一致したシラバス
// documentation: 
func (sc *SyllabusController) GetSyllabusByTeacher(c *gin.Context) {
	var syllabus []models.SyllabusBaseInfo
	teacherName := c.Param("name")
	teacherName = "%" + teacherName + "%" //部分一致
	result := sc.DB.Where("teacher LIKE ?", teacherName).Find(&syllabus)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	
	c.JSON(http.StatusOK, syllabus)
}