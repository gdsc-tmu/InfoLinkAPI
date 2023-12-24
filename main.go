package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"InfoLinkAPI/cmd/models"
	"InfoLinkAPI/cmd/routes"
	_ "InfoLinkAPI/docs"
	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// データベース接続文字列
    dsn := "root:root@(mysql-container:3306)/demo?charset=utf8&parseTime=True&loc=Local"

	// データベースへの接続
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("データベースへの接続に失敗しました" + err.Error())
    }

    // モデルのマイグレーション
    db.AutoMigrate(&models.SyllabusBaseInfo{})

    router := gin.Default()

	// Swaggerの設定
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // ルートとコントローラーを設定
    routes.SyllabusRoutes(router, db)
	routes.SyllabusRandomRoutes(router, db)
	// ....
	// ....

    router.Run(":8080")
}
