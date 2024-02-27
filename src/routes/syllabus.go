package routes

import (
	_ "InfoLinkAPI/src/models"
	"InfoLinkAPI/src/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary シラバス全データ取得
// @Description シラバス全データを取得します。重すぎてswaggerで表示できないので注意。
// @Tags tags
// @Accept  json
// @Produce  json
// @Success 200 {object} models.SyllabusViewModel
// @Router /syllabus [get]
func SyllabusRoutes(router *gin.Engine, db *gorm.DB) {
	syllabusController := controllers.SyllabusController{DB: db}

	router.GET("/syllabus", syllabusController.GetAll)
}