package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSyllabusRandomRoutes(t *testing.T) {
	// データベースのモックを設定
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a GORM database connection", err)
	}

	// テスト用のGinエンジンを設定
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// ルートの設定
	SyllabusRandomRoutes(r, gormDB)

	// データベース呼び出しを模擬
	rows := sqlmock.NewRows([]string{"year", "season", "day", "period", "teacher", "name", "lectureId", "credits", "url", "type", "faculty"}).
		AddRow(2023,"集中","他","0限", "福田 公子", "生命科学特別講義", "R414", 1, "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/1/9/2023_13_R414.html", "大学院科目",  "13")

	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos`.*").WillReturnRows(rows)

	// リクエストをシミュレート
	c.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/random", nil)
	r.ServeHTTP(w, c.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, w.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t, `{"year":2023,"season":"集中","day":"他","period":"0限","teacher":"福田 公子","name":"生命科学特別講義","lectureId":"R414","credits":1,"url":"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/1/9/2023_13_R414.html","type":"大学院科目","faculty":"13"}`, w.Body.String())
}
