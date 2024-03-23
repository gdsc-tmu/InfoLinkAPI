package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ------------------------------
// テスト方針
// ------------------------------
// 以下を検証したい．
// 1. 指定した名前の教員が担当する授業が返されていること
// 2. 指定した名前の教員以外が担当する授業が返されていないこと
// 3. 姓名の間のスペースの有無を吸収できているか
// 4. アルファベットの全角半角に対応できているか
//
// for1,2: test1~3
// for2: test4
// for3: test5
// for4: test3, 6
//
// test1: 異なる2日本人教員について正常な絞り込みかどうか検証
// test2: 似た名前の2日本人教員について正常な絞り込みかどうか検証
// test3: 異なる2外国人教員について正常な絞り込みかどうか検証
// test4: 存在しない教員について結果が空になるか検証
// test5: 姓名間のスペースの違いを吸収できているか検証（日本人教員）
// test6: 全角半角の違いに対応できているか検証（外国人教員）

// union
func TestSyllabusTeacherRoutes(t *testing.T) {
	type TestCase struct {
		name     string
		testFunc func(t *testing.T)
	}

	tests := []TestCase{
		{
			name:     "異なる名前の日本人教員2人",
			testFunc: TestSyllabusTeacherRouteJpName,
		},
		{
			name:     "似た名前の日本人教員2人",
			testFunc: TestSyllabusTeacherRouteJpNameClose,
		},
		{
			name:     "異なる名前の外国人教員2人",
			testFunc: TestSyllabusTeacherRouteEnName,
		},
		{
			name:     "存在しない教員",
			testFunc: TestSyllabusTeacherRouteUnknownName,
		},
		{
			name:     "姓名間のスペースの有無への頑健性（日本人教員）",
			testFunc: TestSyllabusTeacherRouteWhitespaceVariation,
		},
		{
			name:     "全角半角への頑健性（外国人教員）",
			testFunc: TestSyllabusTeacherRouteCharacterVariation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// setup func
func TeacherRouteTestSetup(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock, []string) {
	// mockDBの開設
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	// DBのカラム名を定義しておく
	cols := []string{"year", "season", "day", "period", "teacher", "name", "lecture_id", "credits", "url", "type", "faculty"}
	// GORMからmockDBに接続する
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a GORM database connection", err)
	}
	// Ginエンジンをtest用に設定
	gin.SetMode(gin.TestMode)
	return db, gormDB, mock, cols
}

// test1
func TestSyllabusTeacherRouteJpName(t *testing.T) {
	// 異なる名前の日本人教員について正常に絞り込みが行われているか検証
	// 2ケースについて検証することで，指定していない教員のシラバスが検索結果に含まれていないことを検証できる

	db, gormDB, mock, cols := TeacherRouteTestSetup(t)
	defer db.Close()

	// ------------------------------
	// case1: name=岡本 正吾
	// ------------------------------
	// mockの設定
	expectedRows1 := sqlmock.NewRows(cols).AddRow(
		2023,
		"後期",
		"金",
		"3限",
		"岡本 正吾",
		"バーチャルリアリティ",
		"L0148",
		2,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/5/2023_A6_L0148.html",
		"専門教育科目",
		"A6",
	)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%岡%本%正%吾%").WillReturnRows(expectedRows1)

	// 検証用のcontext
	writer1 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context1, engine1 := gin.CreateTestContext(writer1)

	// routeの設定
	SyllabusTeacherRoutes(engine1, gormDB)

	// リクエストのシミュレート
	context1.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=岡本 正吾", nil)
	engine1.ServeHTTP(writer1, context1.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer1.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[{
    "year": 2023,
    "season": "後期",
    "day": "金",
    "period": "3限",
    "teacher": "岡本 正吾",
    "name": "バーチャルリアリティ",
    "lectureId": "L0148",
    "credits": 2,
    "url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/5/2023_A6_L0148.html",
    "type": "専門教育科目",
    "faculty": "A6"
    }]`,
		writer1.Body.String())

	// ------------------------------
	// case2: name=小町 守
	// ------------------------------
	// mockの設定
	expectedRows2 := sqlmock.NewRows(cols).AddRow(
		2023,
		"集中",
		"他",
		"0限",
		"小町 守",
		"自然言語処理",
		"L0161",
		2,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/9/2023_A6_L0161.html",
		"専門教育科目",
		"A6",
	)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%小%町%守%").WillReturnRows(expectedRows2)

	// 検証用のcontextを設ける
	writer2 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context2, engine2 := gin.CreateTestContext(writer2)

	// routeの設定
	SyllabusTeacherRoutes(engine2, gormDB)

	// リクエストのシミュレート
	context2.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=小町 守", nil)
	engine2.ServeHTTP(writer2, context2.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer2.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[{
    "year": 2023,
    "season": "集中",
    "day": "他",
    "period": "0限",
    "teacher": "小町 守",
    "name": "自然言語処理",
    "lectureId": "L0161",
    "credits": 2,
    "url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/9/2023_A6_L0161.html",
    "type": "専門教育科目",
    "faculty": "A6"
    }]`,
		writer2.Body.String())
}

// test2
func TestSyllabusTeacherRouteJpNameClose(t *testing.T) {
	// 似た名前（同じ苗字）の日本人教員について正常に絞り込みが行われているか検証
	// case1(姓のみ)で異なる名前の教員のシラバスが戻ってくることを検証
	// case2(姓名)でもう一方の教員のシラバスが除外されたことを検証
	// case3 is mirror of case2.

	db, gormDB, mock, cols := TeacherRouteTestSetup(t)
	defer db.Close()

	// ------------------------------
	// case1: name=山口
	// ------------------------------
	// mockの設定
	expectedRows1 := sqlmock.NewRows(cols)
	expectedRows1.AddRow(
		2023,
		"集中",
		"他",
		"0限",
		"山口 大輔",
		"情報と職業",
		"L0001",
		2,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/9/2023_A6_L0001.html",
		"専門教育科目",
		"A6",
	)
	expectedRows1.AddRow(
		2023,
		"後期",
		"火",
		"1限",
		"山口 京一郎",
		"実践英語IIb(612)",
		"A0555",
		1,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/0/2/2023_0D02_A0555.html",
		"言語科目",
		"0D02",
	)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%山%口%").WillReturnRows(expectedRows1)

	// 検証用のcontext
	writer1 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context1, engine1 := gin.CreateTestContext(writer1)

	// routeの設定
	SyllabusTeacherRoutes(engine1, gormDB)

	// リクエストのシミュレート
	context1.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=山口", nil)
	engine1.ServeHTTP(writer1, context1.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer1.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[
		{
			"year": 2023,
			"season": "集中",
			"day": "他",
			"period": "0限",
			"teacher": "山口 大輔",
			"name": "情報と職業",
			"lectureId": "L0001",
			"credits": 2,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/9/2023_A6_L0001.html",
			"type": "専門教育科目",
			"faculty": "A6"
	  	},
		{
			"year": 2023,
			"season": "後期",
			"day": "火",
			"period": "1限",
			"teacher": "山口 京一郎",
			"name": "実践英語IIb(612)",
			"lectureId": "A0555",
			"credits": 1,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/0/2/2023_0D02_A0555.html",
			"type": "言語科目",
			"faculty": "0D02"
		}
	]`,
		writer1.Body.String())

	// ------------------------------
	// case2: name=山口 大輔
	// ------------------------------
	// mockの設定
	expectedRows2 := sqlmock.NewRows(cols).AddRow(
		2023,
		"集中",
		"他",
		"0限",
		"山口 大輔",
		"情報と職業",
		"L0001",
		2,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/9/2023_A6_L0001.html",
		"専門教育科目",
		"A6",
	)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%山%口%大%輔%").WillReturnRows(expectedRows2)

	// 検証用のcontextを設ける
	writer2 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context2, engine2 := gin.CreateTestContext(writer2)

	// routeの設定
	SyllabusTeacherRoutes(engine2, gormDB)

	// リクエストのシミュレート
	context2.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=山口 大輔", nil)
	engine2.ServeHTTP(writer2, context2.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer2.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[
		{
			"year": 2023,
			"season": "集中",
			"day": "他",
			"period": "0限",
			"teacher": "山口 大輔",
			"name": "情報と職業",
			"lectureId": "L0001",
			"credits": 2,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/9/2023_A6_L0001.html",
			"type": "専門教育科目",
			"faculty": "A6"
	  	}
	]`,
		writer2.Body.String())

	// ------------------------------
	// case3: name=山口 京一郎
	// ------------------------------
	// mockの設定
	expectedRows3 := sqlmock.NewRows(cols).AddRow(
		2023,
		"後期",
		"火",
		"1限",
		"山口 京一郎",
		"実践英語IIb(612)",
		"A0555",
		1,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/0/2/2023_0D02_A0555.html",
		"言語科目",
		"0D02",
	)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%山%口%京%一%郎%").WillReturnRows(expectedRows3)

	// 検証用のcontextを設ける
	writer3 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context3, engine3 := gin.CreateTestContext(writer3)

	// routeの設定
	SyllabusTeacherRoutes(engine3, gormDB)

	// リクエストのシミュレート
	context3.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=山口 京一郎", nil)
	engine2.ServeHTTP(writer3, context3.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer3.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[
		{
			"year": 2023,
			"season": "後期",
			"day": "火",
			"period": "1限",
			"teacher": "山口 京一郎",
			"name": "実践英語IIb(612)",
			"lectureId": "A0555",
			"credits": 1,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/0/2/2023_0D02_A0555.html",
			"type": "言語科目",
			"faculty": "0D02"
	  	}
	]`,
		writer3.Body.String())
}

// test3
func TestSyllabusTeacherRouteEnName(t *testing.T) {
	// 異なる名前の外国人教員について正常に絞り込みが行われているか検証
	// 2ケースについて検証することで，指定していない教員のシラバスが検索結果に含まれていないことを検証できる

	db, gormDB, mock, cols := TeacherRouteTestSetup(t)
	defer db.Close()

	// ------------------------------
	// case1: name=Thomas Brotherhood
	// ------------------------------
	// mockの設定
	expectedRows1 := sqlmock.NewRows(cols).AddRow(
		2023,
		"後期",
		"月",
		"3限",
		"ThomasBrotherhood",
		"MigrationandJapan",
		"X0185",
		2,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/0/1/2023_0D05_X0185.html",
		"教養/都市・社会・環境",
		"0D05",
	)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%T%h%o%m%a%s%B%r%o%t%h%e%r%h%o%o%d%").WillReturnRows(expectedRows1)

	// 検証用のcontext
	writer1 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context1, engine1 := gin.CreateTestContext(writer1)

	// routeの設定
	SyllabusTeacherRoutes(engine1, gormDB)

	// リクエストのシミュレート
	context1.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=Thomas Brotherhood", nil)
	engine1.ServeHTTP(writer1, context1.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer1.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[
		{
			"year": 2023,
			"season": "後期",
			"day": "月",
			"period": "3限",
			"teacher": "ThomasBrotherhood",
			"name": "MigrationandJapan",
			"lectureId": "X0185",
			"credits": 2,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/0/1/2023_0D05_X0185.html",
			"type": "教養/都市・社会・環境",
			"faculty": "0D05"
    	}
	]`,
		writer1.Body.String())

	// ------------------------------
	// case2: name=James Baldwin
	// ------------------------------
	// mockの設定
	expectedRows2 := sqlmock.NewRows(cols).AddRow(
		2023,
		"前期",
		"木",
		"6限",
		"JamesBaldwin",
		"医科学英語プレゼンテーションスキルⅠ",
		"U0525",
		2,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/1/4/2023_16_U0525.html",
		"専攻科目",
		"16",
	)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%J%a%m%e%s%B%a%l%d%w%i%n%").WillReturnRows(expectedRows2)

	// 検証用のcontextを設ける
	writer2 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context2, engine2 := gin.CreateTestContext(writer2)

	// routeの設定
	SyllabusTeacherRoutes(engine2, gormDB)

	// リクエストのシミュレート
	context2.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=James Baldwin", nil)
	engine2.ServeHTTP(writer2, context2.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer2.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[{
		"year": 2023,
		"season": "前期",
		"day": "木",
		"period": "6限",
		"teacher": "JamesBaldwin",
		"name": "医科学英語プレゼンテーションスキルⅠ",
		"lectureId": "U0525",
		"credits": 2,
		"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/1/4/2023_16_U0525.html",
		"type": "専攻科目",
		"faculty": "16"
    }]`,
		writer2.Body.String())
}

// test4
func TestSyllabusTeacherRouteUnknownName(t *testing.T) {
	// 存在しない教員名を指定したときに空が返されることを検証
	// 複数のケースについて検証

	db, gormDB, mock, cols := TeacherRouteTestSetup(t)
	defer db.Close()

	// ------------------------------
	// case1: name=戌亥 とこ
	// ------------------------------
	// mockの設定
	expectedRows1 := sqlmock.NewRows(cols)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%戌%亥%と%こ%").WillReturnRows(expectedRows1)

	// 検証用のcontext
	writer1 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context1, engine1 := gin.CreateTestContext(writer1)

	// routeの設定
	SyllabusTeacherRoutes(engine1, gormDB)

	// リクエストのシミュレート
	context1.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=戌亥 とこ", nil)
	engine1.ServeHTTP(writer1, context1.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer1.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[]`,
		writer1.Body.String())

	// ------------------------------
	// case2: name=高橋 恵
	// ------------------------------
	// mockの設定
	expectedRows2 := sqlmock.NewRows(cols)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%高%橋%恵%").WillReturnRows(expectedRows2)

	// 検証用のcontextを設ける
	writer2 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context2, engine2 := gin.CreateTestContext(writer2)

	// routeの設定
	SyllabusTeacherRoutes(engine2, gormDB)

	// リクエストのシミュレート
	context2.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=高橋 恵", nil)
	engine2.ServeHTTP(writer2, context2.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer2.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[]`,
		writer2.Body.String())
}

// test5
func TestSyllabusTeacherRouteWhitespaceVariation(t *testing.T) {
	// 姓名の区切り文字の種類に関する頑健性を検証
	// 区切りなし，半角スペース，全角スペースの3通り．
	// 発行されるSQLクエリの時点で違いを吸収していることを想定する．
	// case1: 区切り無し
	// case2: 半角区切り
	// case3: 全角区切り

	db, gormDB, mock, cols := TeacherRouteTestSetup(t)
	defer db.Close()

	// ------------------------------
	// case1: name=塩田さやか
	// ------------------------------
	// mockの設定
	expectedRows1 := sqlmock.NewRows(cols).AddRow(
		2023,
		"前期",
		"水",
		"4限",
		"塩田 さやか",
		"画像処理（CS）",
		"L0133",
		2,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0133.html",
		"専門教育科目",
		"A6",
	)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%塩%田%さ%や%か%").WillReturnRows(expectedRows1)

	// 検証用のcontext
	writer1 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context1, engine1 := gin.CreateTestContext(writer1)

	// routeの設定
	SyllabusTeacherRoutes(engine1, gormDB)

	// リクエストのシミュレート
	context1.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=塩田さやか", nil)
	engine1.ServeHTTP(writer1, context1.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer1.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[
		{
			"year": 2023,
			"season": "前期",
			"day": "水",
			"period": "4限",
			"teacher": "塩田 さやか",
			"name": "画像処理（CS）",
			"lectureId": "L0133",
			"credits": 2,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0133.html",
			"type": "専門教育科目",
			"faculty": "A6"
		}
	]`,
		writer1.Body.String())

	// ------------------------------
	// case2: name=小野 順貴
	// ------------------------------
	// mockの設定
	expectedRows2 := sqlmock.NewRows(cols).AddRow(
		2023,
		"後期",
		"水",
		"4限",
		"小野 順貴",
		"音響・音声信号処理",
		"L0146",
		2,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0146.html",
		"専門教育科目",
		"A6",
	)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%小%野%順%貴%").WillReturnRows(expectedRows2)

	// 検証用のcontext
	writer2 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context2, engine2 := gin.CreateTestContext(writer2)

	// routeの設定
	SyllabusTeacherRoutes(engine2, gormDB)

	// リクエストのシミュレート
	context2.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=小野 順貴", nil)
	engine2.ServeHTTP(writer2, context2.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer2.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[
		{
			"year": 2023,
			"season": "後期",
			"day": "水",
			"period": "4限",
			"teacher": "小野 順貴",
			"name": "音響・音声信号処理",
			"lectureId": "L0146",
			"credits": 2,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0146.html",
			"type": "専門教育科目",
			"faculty": "A6"
		}
	]`,
		writer2.Body.String())

	// ------------------------------
	// case3: name=松田　崇弘
	// ------------------------------
	// mockの設定
	expectedRows3 := sqlmock.NewRows(cols).AddRow(
		2023,
		"後期",
		"木",
		"2限",
		"松田 崇弘",
		"無線ネットワーク（CS）",
		"L0141",
		2,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/4/2023_A6_L0141.html",
		"専門教育科目",
		"A6",
	)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%松%田%崇%弘%").WillReturnRows(expectedRows3)

	// 検証用のcontext
	writer3 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context3, engine3 := gin.CreateTestContext(writer3)

	// routeの設定
	SyllabusTeacherRoutes(engine3, gormDB)

	// リクエストのシミュレート
	context3.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=松田　崇弘", nil)
	engine3.ServeHTTP(writer3, context3.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer3.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[
		{
			"year": 2023,
			"season": "後期",
			"day": "木",
			"period": "2限",
			"teacher": "松田 崇弘",
			"name": "無線ネットワーク（CS）",
			"lectureId": "L0141",
			"credits": 2,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/4/2023_A6_L0141.html",
			"type": "専門教育科目",
			"faculty": "A6"
		}
	]`,
		writer3.Body.String())
}

// test6
func TestSyllabusTeacherRouteCharacterVariation(t *testing.T) {
	// アルファベットの全角半角に関する頑健性を検証
	// case1: 全半角 -> TestSyllabusTeacherRouteEnName
	// case2: 全全角
	// case3: case2からの絞り込み（検索が機能していることの検証）

	db, gormDB, mock, cols := TeacherRouteTestSetup(t)
	defer db.Close()

	// ------------------------------
	// case1: 全半角, Already verified @ TestSyllabusTeacherRouteEnName.
	// ------------------------------

	// ------------------------------
	// case2: 全全角, name=Ａｄａｍ
	// ------------------------------
	// mockの設定
	expectedRows2 := sqlmock.NewRows(cols)
	expectedRows2.AddRow(
		2023,
		"前期",
		"月",
		"2限",
		"AdamLincCronin",
		"Ecology（生態学各論）",
		"I429",
		2,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/0/1/2023_05_I429.html",
		"専門教育科目",
		"05",
	)
	expectedRows2.AddRow(
		2023,
		"後期",
		"水",
		"4限",
		"VerlAdams",
		"SeminarinSpatialDesignⅡ",
		"L0761",
		2,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0761.html",
		"専門教育科目",
		"A6",
	)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%Ａ%ｄ%ａ%ｍ%").WillReturnRows(expectedRows2)

	// 検証用のcontext
	writer2 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context2, engine2 := gin.CreateTestContext(writer2)

	// routeの設定
	SyllabusTeacherRoutes(engine2, gormDB)

	// リクエストのシミュレート
	context2.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=Ａｄａｍ", nil)
	engine2.ServeHTTP(writer2, context2.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer2.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[
		{
			"year": 2023,
			"season": "前期",
			"day": "月",
			"period": "2限",
			"teacher": "AdamLincCronin",
			"name": "Ecology（生態学各論）",
			"lectureId": "I429",
			"credits": 2,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/0/1/2023_05_I429.html",
			"type": "専門教育科目",
			"faculty": "05"
		},
		{
			"year": 2023,
			"season": "後期",
			"day": "水",
			"period": "4限",
			"teacher": "VerlAdams",
			"name": "SeminarinSpatialDesignⅡ",
			"lectureId": "L0761",
			"credits": 2,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0761.html",
			"type": "専門教育科目",
			"faculty": "A6"
		}
	]`,
		writer2.Body.String())

	// ------------------------------
	// case3: 全全角, name=ＶｅｒｌＡｄａｍｓ
	// ------------------------------
	// mockの設定
	expectedRows3 := sqlmock.NewRows(cols)
	expectedRows3.AddRow(
		2023,
		"後期",
		"水",
		"4限",
		"VerlAdams",
		"SeminarinSpatialDesignⅡ",
		"L0761",
		2,
		"http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0761.html",
		"専門教育科目",
		"A6",
	)
	mock.ExpectQuery("^SELECT \\* FROM `syllabus_base_infos` WHERE teacher LIKE \\?").WithArgs("%Ｖ%ｅ%ｒ%ｌ%Ａ%ｄ%ａ%ｍ%ｓ%").WillReturnRows(expectedRows3)

	// 検証用のcontext
	writer3 := httptest.NewRecorder() // inspectable http.ResponseWriter
	context3, engine3 := gin.CreateTestContext(writer3)

	// routeの設定
	SyllabusTeacherRoutes(engine3, gormDB)

	// リクエストのシミュレート
	context3.Request, _ = http.NewRequest(http.MethodGet, "/syllabus/teacher?name=ＶｅｒｌＡｄａｍｓ", nil)
	engine3.ServeHTTP(writer3, context3.Request)

	// レスポンスをアサート
	assert.Equal(t, http.StatusOK, writer3.Code)
	// レスポンスのボディをアサート
	assert.JSONEq(t,
		`[
		{
			"year": 2023,
			"season": "後期",
			"day": "水",
			"period": "4限",
			"teacher": "VerlAdams",
			"name": "SeminarinSpatialDesignⅡ",
			"lectureId": "L0761",
			"credits": 2,
			"url": "http://www.kyouikujouhou.eas.tmu.ac.jp/syllabus/2023/A/3/2023_A6_L0761.html",
			"type": "専門教育科目",
			"faculty": "A6"
		}
	]`,
		writer3.Body.String())
}
