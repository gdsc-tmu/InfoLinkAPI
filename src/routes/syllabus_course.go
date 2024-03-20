package routes

import (
	"InfoLinkAPI/src/controllers"
	_ "InfoLinkAPI/src/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary 科目名でシラバスを検索します．
// @Description クエリパラメータnameと部分一致するシラバスを返します．
// @Tags tags
// @Accept  json
// @Produce  json
// @Param	name query string true "course name"
// @Success 200 {object} models.SyllabusViewModel
// @failure 400 {object} string "invalid course name exception"
// @Router /syllabus/course [get]
func SyllabusCourseRoutes(router *gin.Engine, db *gorm.DB) {
	syllabusController := controllers.SyllabusController{DB: db}

	// this router needs url params to search syllabus by faculty code.
	router.GET("/syllabus/course", syllabusController.GetSyllabusByCourseName)
}
