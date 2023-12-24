package routes

import (
	"InfoLinkAPI/cmd/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary シラバスデータをランダム取得
// @Description シラバスデータ1つをランダムに取得します
// @Tags tags
// @Accept  json
// @Produce  json
// @Success 200 {object} SyllabusResponse
// @Router /syllabus/random [get]
func SyllabusRandomRoutes(router *gin.Engine, db *gorm.DB) {
	syllabusController := controllers.SyllabusController{DB: db}

	router.GET("/syllabus/random", syllabusController.GetRandom)
}
