package backend

import (
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"strconv"
)

type Task struct {
	Name        string `datastore:"name"`
	Category    string `datastore:"category"`
	Description string `datastore:"description"`
}

func init() {

}

// エンティティグループは1秒間に1回しか書き込みできない制約のテスト。トランザクションあり
func testtranslimit(c echo.Context) error {

	// コンテキストの取得
	ctx := appengine.NewContext(c.Request())

	//parentKey := datastore.NewKey(ctx, "User", "6e013fcf-f48f-4d8c-8ed8-8fbc384ff11b", 0, nil)
	parentKey := datastore.NewKey(ctx, "User", "767437d8-7850-49be-af12-20fd259acbb8", 0, nil)

	// トランザクション開始
	err := datastore.RunInTransaction(ctx, func(c context.Context) error {
		key := datastore.NewIncompleteKey(c, "Task", parentKey)
		if _, err := datastore.Put(c, key, &Task{Name: "ほげ", Category: "ふが", Description: "ぴよ"}); err != nil {
			log.Errorf(c, "datastore.Put: %v", err)
			return err
		}

		return nil
	}, nil)

	if err != nil {
		log.Errorf(ctx, "Transaction failed: %v", err)
		return c.JSON(http.StatusInternalServerError, "Transaction failed")
	}

	return c.JSON(http.StatusOK, "OK")
}

// エンティティグループは1秒間に1回しか書き込みできない。制約のテスト
func testlimit(c echo.Context) error {

	// コンテキストの取得
	ctx := appengine.NewContext(c.Request())

	parentKey := datastore.NewKey(ctx, "User", "6e013fcf-f48f-4d8c-8ed8-8fbc384ff11b", 0, nil)

	// insert
	key := datastore.NewIncompleteKey(ctx, "Task", parentKey)
	if _, err := datastore.Put(ctx, key, &Task{Name: "ほげ", Category: "ふが", Description: "ぴよ"}); err != nil {
		log.Errorf(ctx, "datastore.Put: %v", err)
		return c.JSON(http.StatusInternalServerError, "error2")
	}

	return c.JSON(http.StatusOK, "OK")
}

func createTask(c echo.Context) error {

	// コンテキストの取得
	ctx := appengine.NewContext(c.Request())

	// stringIDを指定した完全キーを作成（Root Entity）
	parentKey := datastore.NewKey(ctx, "User", uuid.NewV4().String(), 0, nil)

	// ユーザーの作成
	user := User{ID: uuid.NewV4().String(), Name: "鈴木"}
	if _, err := datastore.Put(ctx, parentKey, &user); err != nil {
		return c.JSON(http.StatusInternalServerError, "error1")
	}

	// Userを親としたKeyを作成 (Entity Groupとなる)
	key := datastore.NewIncompleteKey(ctx, "Task", nil)

	// タスクの作成
	task := Task{Name: "ほげタスク", Category: "ふがカテゴリ", Description: "ぴよ概要"}
	if _, err := datastore.Put(ctx, key, &task); err != nil {
		log.Errorf(ctx, "datastore.Put: %v", err)
		return c.JSON(http.StatusInternalServerError, "error2")
	}

	return c.JSON(http.StatusCreated, "createTask")
}

func createTasks(c echo.Context) error {

	// コンテキストの取得
	ctx := appengine.NewContext(c.Request())

	// stringIDを指定した完全キーを作成（Root Entity）
	parentKey := datastore.NewKey(ctx, "User", uuid.NewV4().String(), 0, nil)

	// ユーザーの作成
	user := User{ID: uuid.NewV4().String(), Name: "伊藤"}
	if _, err := datastore.Put(ctx, parentKey, &user); err != nil {
		return c.JSON(http.StatusInternalServerError, "error1")
	}

	// Userを親としたKeyを作成 (User->TaskのEntity Groupとなる)
	var keys []*datastore.Key
	keys = append(keys, datastore.NewIncompleteKey(ctx, "Task", parentKey))
	keys = append(keys, datastore.NewIncompleteKey(ctx, "Task", parentKey))
	keys = append(keys, datastore.NewIncompleteKey(ctx, "Task", parentKey))

	// タスクの作成
	var tasks []Task
	tasks = append(tasks, Task{Name: "ほげタスク1", Category: "ふがカテゴリ1", Description: "ぴよ概要1"})
	tasks = append(tasks, Task{Name: "ほげタスク2", Category: "ふがカテゴリ2", Description: "ぴよ概要2"})
	tasks = append(tasks, Task{Name: "ほげタスク3", Category: "ふがカテゴリ3", Description: "ぴよ概要3"})

	// Bulk Insert
	if _, err := datastore.PutMulti(ctx, keys, tasks); err != nil {
		log.Errorf(ctx, "datastore.Put: %v", err)
		return c.JSON(http.StatusInternalServerError, "error2")
	}

	// Userに紐づくTask一覧を取得
	var userTask []Task
	query := datastore.NewQuery("Task").Ancestor(parentKey)

	// AncestorにParentKeyを指定したクエリはAncestor(祖先クエリ)と呼ばれ 強整合性(悲観的ロック)を確保したクエリが実行可能
	if _, err := query.GetAll(ctx, &userTask); err != nil {
		log.Errorf(ctx, "query.GetAll: %v", err)
		return c.JSON(http.StatusInternalServerError, "error3")
	}
	/*
		userTask => [
				{"Name":"ほげタスク1","Category":"ふがカテゴリ1","Description":"ぴよ概要1"},
				{"Name":"ほげタスク3","Category":"ふがカテゴリ3","Description":"ぴよ概要3"},
				{"Name":"ほげタスク2","Category":"ふがカテゴリ2","Description":"ぴよ概要2"}
			]
	*/

	// userTaskをイテレートする
	result := map[int64]string{}
	it := query.Run(ctx)
	for {
		var task Task
		key, err := it.Next(&task)

		// 終端まで到達したか
		if err == datastore.Done {
			break
		}

		// 何らかのエラー
		if err != nil {
			log.Errorf(ctx, "500 error: %v", err)
			break
		}

		// key番目にタスク名を代入
		result[key.IntID()] = task.Name
	}

	return c.JSON(http.StatusOK, result)
	// result => {"5048407638933504":"ほげタスク1","5611357592354816":"ほげタスク3","6174307545776128":"ほげタスク2"}
}

func getTask(c echo.Context) error {

	// コンテキストの取得
	ctx := appengine.NewContext(c.Request())

	pid := c.Param("id")
	id, err := strconv.Atoi(pid)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusInternalServerError, "パラメータが不正です")
	}

	task := new(Task)
	key := datastore.NewKey(ctx, "Task", "", int64(id), nil)

	// Keyで取得
	// Get http://localhost:8080/task/5066549580791808
	err = datastore.Get(ctx, key, task)
	if err != nil {
		log.Errorf(ctx, "datastore.Get: %v", err)
		return c.JSON(http.StatusInternalServerError, "datastore.Get")
	}

	for i := 0; i < 500; i++ {
		// update
		if _, err := datastore.Put(ctx, key, &Task{Name: "ほげ", Category: "ふが", Description: "ぴよ"}); err != nil {
			log.Errorf(ctx, "datastore.Put: %v", err)
			return c.JSON(http.StatusInternalServerError, "error2")
		}
	}

	return c.JSON(http.StatusCreated, task)
}
