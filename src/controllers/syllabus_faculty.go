package controllers

import (
	"InfoLinkAPI/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSyllabusByFaculty 指定した学部コードのシラバスを返す。
//
// 引数: 学部コード e.g.A6
// 戻り値: Facultyフィールドが指定した学部コードであるレコード
// documentation: https://www.notion.so/24f67335e99344d0b454168b722af1ae?pvs=4#8ae439dc15f84d9297cf4ef1731e1dea
func (sc *SyllabusController) GetSyllabusByFaculty(c *gin.Context) {
	var syllabus []models.SyllabusBaseInfo
	facultyCode := c.Query("code")
	result := sc.DB.Where("faculty = ?", facultyCode).Find(&syllabus)
	
	// handle invalid faculty code
	_, valid := models.FacultyMap[facultyCode]
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty code. See: https://www.notion.so/24f67335e99344d0b454168b722af1ae?pvs=4#8ae439dc15f84d9297cf4ef1731e1dea"})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	
	c.JSON(http.StatusOK, syllabus)
}