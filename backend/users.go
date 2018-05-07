package backend

import (
	"net/http"

	"github.com/labstack/echo"

	"cloud.google.com/go/spanner"
	"github.com/satori/go.uuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"strconv"
)

/*

  Datastoreのテスト

*/

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func init() {

	//e.Use(middleware.CORS())

	e.POST("/users", createUser)
	e.GET("/users", getUsers)
	e.GET("/users/:id", getUser)
	e.GET("/users/logtest", logtest)

	e.GET("/spa", spa)
}

func spa(c echo.Context) error {

	// コンテキストの取得
	ctx := appengine.NewContext(c.Request())

	// This database must exist.
	databaseName := ""

	client, err := spanner.NewClient(ctx, databaseName)
	if err != nil {
		log.Errorf(ctx, "Failed to create client %v", err)
	}
	defer client.Close()

	stmt := spanner.Statement{SQL: "SELECT 1"}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	row, err := iter.Next()
	if err != nil {
		log.Errorf(ctx, "Query failed with %v", err)
	}

	var i int64
	if row.Columns(&i) != nil {
		log.Errorf(ctx, "Failed to parse row %v", err)
	}

	return c.JSON(http.StatusOK, i)
}

// API reference
// https://cloud.google.com/appengine/docs/standard/go/datastore/reference?hl=ja#Query

func createUser(c echo.Context) error {

	// ダミーデータの作成
	user := User{ID: uuid.NewV4().String(), Name: "山田"}

	// コンテキストの取得
	ctx := appengine.NewContext(c.Request())

	// keyを作成し、DataStoreへinsert
	// keyは文字列ではなく、Datastoreのエンティティを管理する構造体 https://cloud.google.com/datastore/docs/concepts/entities?hl=ja#key
	key := datastore.NewIncompleteKey(ctx, "User", nil)

	// insert or update
	if _, err := datastore.Put(ctx, key, &user); err != nil {
		log.Errorf(ctx, "datastore.Put: %v", err)
		return c.JSON(http.StatusInternalServerError, "error")
	}

	return c.JSON(http.StatusCreated, user)
}

func getUsers(c echo.Context) error {

	// コンテキストの取得
	ctx := appengine.NewContext(c.Request())

	// エンティティを分類した種類(kind)で取得
	q := datastore.NewQuery("User")
	var users []User
	if key, err := q.GetAll(ctx, &users); err != nil {
		log.Errorf(ctx, "%v:%v", key, err)
	}

	//// エンティティIDを指定して取得
	//var keys []*datastore.Key
	//keys = append(keys, datastore.NewKey(ctx, "User", "", 5328783104016384, nil)) // IDは適宜変更
	//keys = append(keys, datastore.NewKey(ctx, "User", "", 5891733057437696, nil)) // IDは適宜変更
	//
	//users := make([]*User, len(keys))
	//err := datastore.GetMulti(ctx, keys, users)
	//if err != nil {
	//	log.Errorf(ctx, "datastore.GetMulti: %v", err)
	//	return c.JSON(http.StatusInternalServerError, "error")
	//}

	return c.JSON(http.StatusOK, users)
}

func getUser(c echo.Context) error {

	// コンテキストの取得
	ctx := appengine.NewContext(c.Request())

	pid := c.Param("id")
	id, err := strconv.Atoi(pid)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusInternalServerError, "パラメータが不正です")
	}

	user := new(User)
	key := datastore.NewKey(ctx, "User", "", int64(id), nil) // IDは適宜変更

	// Keyで取得
	// Get http://localhost:8080/users/5328783104016384
	err = datastore.Get(ctx, key, user)
	if err != nil {
		log.Errorf(ctx, "datastore.Get: %v", err)
		return c.JSON(http.StatusInternalServerError, "getUser")
	}

	return c.JSON(http.StatusOK, user)
}

// ログテスト
func logtest(c echo.Context) error {

	// コンテキストの取得
	ctx := appengine.NewContext(c.Request())

	// 通常ログ
	log.Infof(ctx, "Info!!")

	// エラーログ
	log.Errorf(ctx, "Error")

	// 警告ログ
	log.Warningf(ctx, "Warning!!")

	// デバッグログ
	log.Debugf(ctx, "Debug!!")

	return c.JSON(http.StatusOK, "logged")
}
