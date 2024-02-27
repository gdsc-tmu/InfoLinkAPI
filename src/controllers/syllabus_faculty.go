package controllers

import (
	"InfoLinkAPI/cmd/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// [メソッド名] [動詞で始まる簡潔な説明]。
//
// [引数についての詳細説明（必要な場合）]
// [戻り値についての詳細説明（必要な場合）]
// [その他特記事項があれば記述]
func (sc *SyllabusController) GetSyllabusByFaculty(c *gin.Context) {
	var syllabus []models.SyllabusBaseInfo
	facultyCode := c.Query("code")
	result := sc.DB.Where("Faculty = ?", facultyCode).Find(&syllabus)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	
	c.JSON(http.StatusOK, syllabus)
}