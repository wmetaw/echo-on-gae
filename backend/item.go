package backend

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"google.golang.org/appengine"
	"net/http"
	"os"
	"strconv"
)

/*

  CloudSQLのテスト


  GAEでテストする場合はsecret.yamlに環境変数を記述する

env_variables:
  CLOUDSQL_CONNECTION_NAME: "hoge"
  CLOUDSQL_USER: "hoge"
  CLOUDSQL_PASSWORD: "hoge"

*/

type Item struct {
	Id          uint `gorm:primary_key`
	Category    uint16
	Name        string
	Description string
	Price       uint
	IsShop      bool
}

var rdb *gorm.DB

func init() {

	rdb, _ = gorm.Open("mysql", rdbDsn())
}

// getItems 全件取得
func setItems(c echo.Context) error {

	items := []*Item{
		{Id: 1, Category: 1, Name: "ポーション", Description: "HP回復薬", Price: 100, IsShop: true},
		{Id: 2, Category: 2, Name: "エーテル", Description: "MP回復薬", Price: 300, IsShop: true},
		{Id: 3, Category: 3, Name: "エリクサー", Description: "HP,MP全回復", Price: 500, IsShop: false},
	}
	for _, i := range items {
		rdb.Create(&i)
	}

	return c.JSON(http.StatusCreated, items)
}

// getItems 全件取得
func getItems(c echo.Context) error {

	var allItem []Item
	rdb.Find(&allItem)

	return c.JSON(http.StatusCreated, allItem)
}

// getItem 1件取得
func getItem(c echo.Context) error {

	pid := c.Param("id")
	id, err := strconv.Atoi(pid)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusInternalServerError, "パラメータが不正です")
	}

	var i Item
	rdb.First(&i, id)

	return c.JSON(http.StatusCreated, i)
}

// rdbDsn 本番またはローカルのDSNを取得
func rdbDsn() string {

	// local サンドボックス環境
	if appengine.IsDevAppServer() {

		user := getEnv("MYSQL_USER", "root")
		pwd := getEnv("MYSQL_PASSWORD", "root")
		host := getEnv("MYSQL_HOST", "127.0.0.1")
		port := getEnv("MYSQL_PORT", "3306")
		db := getEnv("MYSQL_DB", "dev")
		opt := "charset=utf8&parseTime=True&loc=Asia%2FTokyo"

		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, pwd, host, port, db, opt)
	}

	// GAE
	connectionName := getEnv("CLOUDSQL_CONNECTION_NAME", "")
	user := getEnv("CLOUDSQL_USER", "")
	password := getEnv("CLOUDSQL_PASSWORD", "")

	return fmt.Sprintf("%s:%s@cloudsql(%s)/dev", user, password, connectionName)
}

// getEnv 環境変数から値を取得。なければデフォルト値をセット
func getEnv(name, def string) string {
	env := os.Getenv(name)
	if len(env) != 0 {
		return env
	}
	return def
}
