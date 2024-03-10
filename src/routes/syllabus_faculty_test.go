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

// ------------------------------
// テスト方針
// ------------------------------
// SyllabusFacultyRoutes()では
// 0. if code is invalid -> response 400 w/ err msg
// 1. 検索クエリを生成する
// 2. DBに対して検索を実行する
// 3. response 200 with results
// の流れで処理をする．
// このテストではDBをmockして，生成されるクエリと検索結果は矛盾の無いものを仮定する．
// 主に検証するのは正しいロジック，すなわち
// - 無効な学部コードが指定された時にBadRequestを返す(test3)
// - 有効なコードであればOKを返す(test1,2)
// - 指定した学部のシラバスのみが返される(test1)
// を満たしているか否か．

// union
func TestSyllabusFacultyRoutes(t *testing.T) {
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
			testFunc: TestSyllabusFacultyRoutesValidCode,
		},
		{
			// 学部のみの絞り込みではありえないが，複数条件で絞り込むと起こりうるケース
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

// test1
func TestSyllabusFacultyRoutesValidCode(t *testing.T) {
	// 複数の学部コードを検証する
	// 1つだと指定した学部"のみ"が返ってくるかわからない．e.g.常にある1つの学部を返すふるまい

	// mockDBの開設
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
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
	// case1: code=A6
	// ------------------------------
	// mockの設定
	expectedRowsA6 := sqlmock.NewRows(cols).AddRow(2023, "後期", "水", "2限", "岡本 正吾", "データ構造とアルゴリズム演習（CS）", "L0111", 2, "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0111.html", "専門教育科目", "A6")
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE faculty = \\?").WithArgs("A6").WillReturnRows(expectedRowsA6)

	// 検証用のcontextを設ける
	writerA6 := httptest.NewRecorder() // inspectable http.ResponseWriter
	contextA6, engineA6 := gin.CreateTestContext(writerA6)

	// routeの設定
	SyllabusFacultyRoutes(engineA6, gormDB)

	// リクエストのシミュレート
	contextA6.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/faculty/A6", nil)
	engineA6.ServeHTTP(writerA6, contextA6.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writerA6.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
	`[{
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
	}]`,
	writerA6.Body.String())

	// ------------------------------
	// case2: code=0D01
	// ------------------------------
	// mockの設定
	expectedRows0D01 := sqlmock.NewRows(cols).AddRow(2023, "前期", "月", "5限", "會田 雅樹", "基礎ゼミナール", "A0127", 2, "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/0/1/2023_0D01_A0127.html", "基礎科目群", "0D01")
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE faculty = \\?").WithArgs("0D01").WillReturnRows(expectedRows0D01)

	// 検証用のcontextを設ける
	writer0D01 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context0D01, engine0D01 := gin.CreateTestContext(writer0D01)

	// routeの設定
	SyllabusFacultyRoutes(engine0D01, gormDB)

	// リクエストのシミュレート
	context0D01.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/faculty/0D01", nil)
	engine0D01.ServeHTTP(writer0D01, context0D01.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer0D01.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
	`[{
    "Year": 2023,
    "Season": "前期",
    "Day": "月",
    "Period": "5限",
    "Teacher": "會田 雅樹",
    "Name": "基礎ゼミナール",
    "LectureId": "A0127",
    "Credits": 2,
    "URL": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/0/1/2023_0D01_A0127.html",
    "Type": "基礎科目群",
    "Faculty": "0D01",
    "DeletedAt": null
	}]`,
	writer0D01.Body.String())
}

// test2
func TestSyllabusFacultyRoutesValidCodeResultUnHit(t *testing.T) {
	// mockDBの開設
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
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
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE faculty = \\?").WithArgs("A6").WillReturnRows(expectedRows)

	// 検証用のcontextを設ける
	writer := httptest.NewRecorder() // inspectable http.ResponseWriter
	context, engine := gin.CreateTestContext(writer)

	// routeの設定
	SyllabusFacultyRoutes(engine, gormDB)

	// リクエストのシミュレート
	context.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/faculty/A6", nil)
	engine.ServeHTTP(writer, context.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
	`[]`,
	writer.Body.String())
}

// test3
func TestSyllabusFacultyRoutesInValidCode(t *testing.T) {
	// mockDBの開設
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
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
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE faculty = \\?").WithArgs("111111TMU").WillReturnRows(expectedRows)

	// 検証用のcontextを設ける
	writer := httptest.NewRecorder() // inspectable http.ResponseWriter
	context, engine := gin.CreateTestContext(writer)

	// routeの設定
	SyllabusFacultyRoutes(engine, gormDB)

	// リクエストのシミュレート
	context.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/faculty/111111TMU", nil)
	engine.ServeHTTP(writer, context.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t, `{"error": "Invalid faculty code. See: https://www.notion.so/24f67335e99344d0b454168b722af1ae?pvs=4#8ae439dc15f84d9297cf4ef1731e1dea"}`, writer.Body.String())
}
