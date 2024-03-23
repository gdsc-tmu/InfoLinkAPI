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

// union
func TestSyllabusCourseRoutes(t *testing.T) {
	// ------------------------------
	// test1: クエリパラメータで科目名が指定され，該当レコードが存在する場合
	// test2: クエリパラメータで科目名が指定され，該当レコードが存在しない場合
	// test3: クエリパラメータが存在しない場合
	// ------------------------------
	type TestCase struct {
		name string
		testFunc func(t *testing.T)
	}

	tests := []TestCase{
		{
			name: "test1",
			testFunc: TestSyllabusCourseRoutesValidNameResultHit,
		},
		{
			name: "test2",
			testFunc: TestSyllabusCourseRoutesValidNameResultUnHit,
		},
		{
			name: "test3",
			testFunc: TestSyllabusCourseRoutesNoQueryParam,
		},
	}

	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T){
			tt.testFunc(t)
		})
	}
}

// test1
// クエリパラメータで科目名が指定され，該当レコードが存在する場合
func TestSyllabusCourseRoutesValidNameResultHit(t *testing.T) {
	// mockDBの開設
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close() // 全部終わったらdbを閉じる
	// DBのカラム名を定義しておく
	cols := []string{"year", "season", "day", "period", "teacher", "name", "lecture_id", "credits", "url", "type", "faculty"}

	// GORMからmockDBに接続する
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a GORM database connection", err)
	}

	// Ginエンジンをtest用に設定
	gin.SetMode(gin.TestMode)

	// ------------------------------
	// case1: name=データ構造とアルゴリズム演習（CS）
	// ------------------------------
	{
		// mockの設定
		expectedRows := sqlmock.NewRows(cols).AddRow(2023, "後期", "水", "2限", "岡本 正吾", "データ構造とアルゴリズム演習（CS）", "L0111", 2, "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0111.html", "専門教育科目", "A6")
		mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE name LIKE \\?").WithArgs("%データ構造とアルゴリズム演習（CS）%").WillReturnRows(expectedRows)

		// 検証用のcontextを設ける
		writer := httptest.NewRecorder() // inspectable http.ResponseWriter
		context, engine := gin.CreateTestContext(writer)

		// routeの設定
		SyllabusCourseRoutes(engine, gormDB)

		// リクエストのシミュレート
		context.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/course", nil)
		// クエリパラメータを設定
		context.Request.URL.RawQuery = "name=データ構造とアルゴリズム演習（CS）"
		engine.ServeHTTP(writer, context.Request)

		// レスポンスをアサート
		assert.Equal(t, http.StatusOK, writer.Code)
		// レスポンスのボディをアサート
		assert.JSONEq(t,
		`[{
			"year": 2023,
			"season": "後期",
			"day": "水",
			"period": "2限",
			"teacher": "岡本 正吾",
			"name": "データ構造とアルゴリズム演習（CS）",
			"lectureId": "L0111",
			"credits": 2,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0111.html",
			"type": "専門教育科目",
			"faculty": "A6"
		}]`,
		writer.Body.String())
	}
	// ------------------------------
	// case2: name=基礎ゼミナール
	// ------------------------------
	{
		// mockの設定
		expectedRows := sqlmock.NewRows(cols).AddRow(2023, "前期", "月", "5限", "會田 雅樹", "基礎ゼミナール", "A0127", 2, "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/0/1/2023_0D01_A0127.html", "基礎科目群", "0D01")
		mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE name LIKE \\?").WithArgs("%基礎ゼミナール%").WillReturnRows(expectedRows)

		// 検証用のcontextを設ける
		writer := httptest.NewRecorder() // inspectable http.ResponseWriter
		context, engine := gin.CreateTestContext(writer)

		// routeの設定
		SyllabusCourseRoutes(engine, gormDB)

		// リクエストのシミュレート
		context.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/course", nil)
		// クエリパラメータを設定
		context.Request.URL.RawQuery = "name=基礎ゼミナール"
		engine.ServeHTTP(writer, context.Request)

		// レスポンスをアサート
		assert.Equal(t, http.StatusOK, writer.Code)
		// レスポンスのボディをアサート
		assert.JSONEq(t,
		`[{
			"year": 2023,
			"season": "前期",
			"day": "月",
			"period": "5限",
			"teacher": "會田 雅樹",
			"name": "基礎ゼミナール",
			"lectureId": "A0127",
			"credits": 2,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/0/1/2023_0D01_A0127.html",
			"type": "基礎科目群",
			"faculty": "0D01"
		}]`,
		writer.Body.String())
	}
}

// test2
// クエリパラメータで科目名が指定され，該当レコードが存在しない場合
func TestSyllabusCourseRoutesValidNameResultUnHit(t *testing.T) {
	// mockDBの開設
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close() // 全部終わったらdbを閉じる
	// DBのカラム名を定義しておく
	cols := []string{"year", "season", "day", "period", "teacher", "name", "lecture_id", "credits", "url", "type", "faculty"}

	// GORMからmockDBに接続する
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a GORM database connection", err)
	}

	// Ginエンジンをtest用に設定
	gin.SetMode(gin.TestMode)

	// mockの設定
	expectedRows := sqlmock.NewRows(cols) // blank result
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE name LIKE \\?").WithArgs("%デタゴリ🦍%").WillReturnRows(expectedRows)

	// 検証用のcontextを設ける
	writer := httptest.NewRecorder() // inspectable http.ResponseWriter
	context, engine := gin.CreateTestContext(writer)

	// routeの設定
	SyllabusCourseRoutes(engine, gormDB)

	// リクエストのシミュレート
	context.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/course", nil)
	// クエリパラメータを設定
	context.Request.URL.RawQuery = "name=デタゴリ🦍"
	engine.ServeHTTP(writer, context.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t, `[]`, writer.Body.String())
}

// test3
// クエリパラメータが存在しない場合
func TestSyllabusCourseRoutesNoQueryParam(t *testing.T) {
	// mockDBの開設
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close() // 全部終わったらdbを閉じる
	// DBのカラム名を定義しておく
	cols := []string{"year", "season", "day", "period", "teacher", "name", "lecture_id", "credits", "url", "type", "faculty"}

	// GORMからmockDBに接続する
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a GORM database connection", err)
	}

	// Ginエンジンをtest用に設定
	gin.SetMode(gin.TestMode)

	// mockの設定
	expectedRows := sqlmock.NewRows(cols).AddRow(2023, "後期", "水", "2限", "岡本 正吾", "データ構造とアルゴリズム演習（CS）", "L0111", 2, "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0111.html", "専門教育科目", "A6")
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE name LIKE \\?").WithArgs("%データ構造とアルゴリズム演習（CS）%").WillReturnRows(expectedRows)

	// 検証用のcontextを設ける
	writer := httptest.NewRecorder() // inspectable http.ResponseWriter
	context, engine := gin.CreateTestContext(writer)

	// routeの設定
	SyllabusCourseRoutes(engine, gormDB)

	// リクエストのシミュレート
	context.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/course", nil)
	engine.ServeHTTP(writer, context.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	// レスポンスのボディをアサート
	assert.Equal(t, `{"error":"Invalid course name. See: https://www.notion.so/42e2fc5ed65a4ba2b6c3ea8bd4dcaad8?pvs=4"}`, writer.Body.String())
}