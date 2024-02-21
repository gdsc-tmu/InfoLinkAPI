package routes

import (
	_ "InfoLinkAPI/cmd/models"
	"InfoLinkAPI/cmd/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary 学部コードでシラバスを検索します．
// @Description パラメータ引数に与えた学部コードに一致するシラバスを返します．
// @Tags tags
// @Accept  json
// @Produce  json
// @Success 200 {object} models.SyllabusViewModel
// @Router /syllabus/random [get]
func SyllabusFacultyRoutes(router *gin.Engine, db *gorm.DB) {
	syllabusController := controllers.SyllabusController{DB: db}

	router.GET("/syllabus/faculties/{}", syllabusController.GetSyllabusByFaculty)
}
