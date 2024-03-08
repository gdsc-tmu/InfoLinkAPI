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

func TestSyllabusFacultyRoutesValidCodeResultHit(t *testing.T) {
	// valid faculty code and result exist
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
	SyllabusFacultyRoutes(r, gormDB)

	// setting up mock DB
	rows := sqlmock.NewRows([]string{"year", "season", "day", "period", "teacher", "name", "lecture_id", "credits", "url", "type", "faculty"})
	// A6 example
	rows.AddRow(2023, "後期", "水", "2限", "岡本 正吾", "データ構造とアルゴリズム演習（CS）", "L0111", 2, "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0111.html", "専門教育科目", "A6")

	// データベース呼び出しを模擬
	// 生成されるクエリのシミュレート
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE faculty = \\?").WithArgs("A6").WillReturnRows(rows)

	// APIリクエストと，それに対するレスポンスの検証
	// リクエストをシミュレート
	c.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/faculties/A6", nil)
	r.ServeHTTP(w, c.Request)
	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, w.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t, `[{
    "Year": 2023,
    "Season": "後期",
    "Day": "水",
    "Period": "2限",
    "Teacher": "岡本 正吾",
    "Name": "データ構造とアルゴリズム演習（CS）",
    "LectureId": "L0111",
    "Credits": 2,
    "URL": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0111.html",
    "Type": "専門教育科目",
    "Faculty": "A6",
    "DeletedAt": null
  }]`, w.Body.String())
}

func TestSyllabusFacultyRoutesValidCodeResultUnHit(t *testing.T) {
	// valid faculty code and result NOT exist
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
	SyllabusFacultyRoutes(r, gormDB)

	// setting up mock DB
	rows := sqlmock.NewRows([]string{"year", "season", "day", "period", "teacher", "name", "lecture_id", "credits", "url", "type", "faculty"})

	// データベースの呼び出しを模擬
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE faculty = \\?").WithArgs("0D01").WillReturnRows(rows)
	c.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/faculties/0D01", nil)
	r.ServeHTTP(w, c.Request)
	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, w.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t, `[]`, w.Body.String())
}

func TestSyllabusFacultyRoutesInValidCode(t *testing.T) {
	// Invalid faculty code
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
	SyllabusFacultyRoutes(r, gormDB)

	// setting up mock DB
	rows := sqlmock.NewRows([]string{"year", "season", "day", "period", "teacher", "name", "lecture_id", "credits", "url", "type", "faculty"})

	// データベースの呼び出しを模擬
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE faculty = \\?").WithArgs("111111TMU").WillReturnRows(rows)
	c.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/faculties/111111TMU", nil)
	r.ServeHTTP(w, c.Request)
	// レスポンスをアサート
	assert.Equal(t, http.StatusBadRequest, w.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t, `{"error": "Invalid faculty code. See: https://www.notion.so/24f67335e99344d0b454168b722af1ae?pvs=4#8ae439dc15f84d9297cf4ef1731e1dea"}`, w.Body.String())
}

func TestSyllabusFacultyRoutes(t *testing.T) {
	// "DBに対するクエリとその返り値"を仮定し，"APIのリクエストとそのレスポンス"に対して検証を行う．

	// ------------------------------
	// test1: code有効 and 該当レコードが存在する
	// test2: code有効 and 該当レコードが存在しない
	// test3: code無効
	// ------------------------------
	type TestCase struct {
		name string
		testFunc func(t *testing.T)
	}

	tests := []TestCase{
		{
			name: "有効な各部コード，該当レコードが存在する",
			testFunc: TestSyllabusFacultyRoutesValidCodeResultHit,
		},
		{
			name: "有効な各部コード，該当レコードが存在しない",
			testFunc: TestSyllabusFacultyRoutesValidCodeResultUnHit,
		},
		{
			name: "無効な学部コード",
			testFunc: TestSyllabusFacultyRoutesInValidCode,
		},
	}

	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T){
			tt.testFunc(t)
		})
	}
}
