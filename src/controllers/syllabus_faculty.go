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
// 学部コードは https://github.com/tenk-9/tmuSyllabus_scrapingに一覧があります．
func (sc *SyllabusController) GetSyllabusByFaculty(c *gin.Context) {
	var syllabus []models.SyllabusBaseInfo
	facultyCode := c.Param("code")
	result := sc.DB.Where("faculty = ?", facultyCode).Find(&syllabus)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	
	c.JSON(http.StatusOK, syllabus)
}