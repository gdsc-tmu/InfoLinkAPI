package controllers

import (
	"net/http"
	"InfoLinkAPI/src/models"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// [型名] [型についての簡潔な説明]。
//
// [より詳細な説明や使用例など]
type SyllabusController struct {
	DB *gorm.DB
}

// [メソッド名] [動詞で始まる簡潔な説明]。
//
// [引数についての詳細説明（必要な場合）]
// [戻り値についての詳細説明（必要な場合）]
// [その他特記事項があれば記述]
func (sc *SyllabusController) GetAll(c *gin.Context) {
	var syllabus []models.SyllabusBaseInfo
	result := sc.DB.Find(&syllabus)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	res := make([]models.SyllabusViewModel, 0)
	for _, s := range syllabus {
		res = append(res, s.ToSyllabusViewModel())
	}
	
	c.JSON(http.StatusOK, res)
}