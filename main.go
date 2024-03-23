package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"InfoLinkAPI/src/models"
	"InfoLinkAPI/src/routes"
	"net/http"
	_ "InfoLinkAPI/docs"
	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// データベース接続文字列
    dsn := "root:root@(mysql-container:3306)/demo?charset=utf8&parseTime=True&loc=Local"

	// データベースへの接続
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
    if err != nil {
        panic("データベースへの接続に失敗しました" + err.Error())
    }

    // モデルのマイグレーション
    db.AutoMigrate(&models.SyllabusBaseInfo{})

    router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
	
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
	
		c.Next()
	})

	// Swaggerの設定
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // ルートとコントローラーを設定
    routes.SyllabusRoutes(router, db)
	routes.SyllabusRandomRoutes(router, db)
	routes.SyllabusFacultyRoutes(router, db)
	routes.SyllabusTeacherRoutes(router, db)
	// ....
	// ....

    router.Run(":8080")
}
