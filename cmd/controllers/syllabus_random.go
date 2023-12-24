package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"InfoLinkAPI/cmd/models"
)

// GetRandom - シラバスデータをランダムに1つ取得
func (sc *SyllabusController) GetRandom(c *gin.Context) {
	var syllabus models.SyllabusBaseInfo
	// データベースからランダムなレコードを取得
	err := sc.DB.Order("RAND()").First(&syllabus).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データの取得に失敗しました"})
		return
	}
	// ランダムに取得したレコードを返す
	c.JSON(http.StatusOK, syllabus)
}
