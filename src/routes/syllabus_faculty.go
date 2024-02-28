package routes

import (
	"InfoLinkAPI/src/controllers"
	_ "InfoLinkAPI/src/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary 学部コードでシラバスを検索します．
// @Description パラメータ引数に与えた学部コードに一致するシラバスを返します．
// @Tags tags
// @Accept  json
// @Produce  json
// @Param	faculty	query	string	true	"faculty code"
// @Success 200 {object} models.SyllabusViewModel
// @failure 400 {object} string "invalid faculty code exception"
// @Router /syllabus/search [get]
func SyllabusFacultyRoutes(router *gin.Engine, db *gorm.DB) {
	syllabusController := controllers.SyllabusController{DB: db}

	// this router needs url params to search syllabus by faculty code.
	router.GET("/syllabus/search", syllabusController.GetSyllabusByFaculty)
}
