package controllers

import (
	"InfoLinkAPI/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSyllabusByCourseName 指定した科目名のシラバスを取得し，HTTPレスポンスを行う。
//
// 部分一致で検索する。
// APIドキュメント: https://www.notion.so/42e2fc5ed65a4ba2b6c3ea8bd4dcaad8?pvs=4
func (sc *SyllabusController) GetSyllabusByCourseName(c *gin.Context) {
	// シラバス配列
	var syllabusBaseInfos []models.SyllabusBaseInfo
	// クエリパラメータから科目名を取得
	courseName := c.Query("name")
	// クエリパラメータがない場合エラー
	if courseName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course name. See: https://www.notion.so/42e2fc5ed65a4ba2b6c3ea8bd4dcaad8?pvs=4"})
		return
	}

	// 部分一致で検索
	result := sc.DB.Where("name LIKE ?", "%"+courseName+"%").Find(&syllabusBaseInfos)
	// エラーハンドリング
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// ビューモデルに変換
	syllabusViewModels := models.GetSyllabusViewModelBySyllabusBaseInfo(syllabusBaseInfos)

	// HTTPレスポンス
	c.JSON(http.StatusOK, syllabusViewModels)
}