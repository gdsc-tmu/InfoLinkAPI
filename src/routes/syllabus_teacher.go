package routes

import (
	"InfoLinkAPI/src/controllers"
	_ "InfoLinkAPI/src/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary 教員名でシラバスを検索します．
// @Description パラメータに与えた教員名を，syllabus_base_infos.teacherとの部分一致で検索します．
// @Tags tags
// @Accept  json
// @Produce  json
// @Param	name	path	string	true	"teacher name"
// @Success 200 {object} models.SyllabusViewModel
// @Router /syllabus/teacher/{name} [get]
func SyllabusTeacherRoutes(router *gin.Engine, db *gorm.DB){
	syllabusController := controllers.SyllabusController{DB: db}
	router.GET("/syllabus/teacher/:name", syllabusController.GetSyllabusByTeacher)
}
