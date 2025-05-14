package sqlite_test

import (
	"finalProject/internal/storage/sqlite"
	"fmt"
	"testing"
)

func TestDataBase(t *testing.T) {
	//t.Run("New", TestNew)
	storage, err := sqlite.New("store.db")
	if err != nil {
		fmt.Errorf("Error %v", err)
		//return
	}
	userId, err := storage.AddUser("egor", "1234")
	if err != nil {
		fmt.Errorf("Error %v", err)
		return
	}
	//id1, _ := strconv.Atoi(userId)
	exID, err := storage.AddExpression(&sqlite.Expression{
		ID:         userId,
		Expression: "2+2",
	})
	if err != nil {
		fmt.Errorf("Error %v", err)
		return
	}
	//storage.SetResult(exID, "4")
	ex := storage.GetExpressionById(exID, userId)
	fmt.Println(ex.Expression)
	err = storage.SetResult(exID, "4")
	if err != nil {
		fmt.Errorf("Error %v", err)
		return
	}
	storage.UpdateUserPassword(userId, "qwerty")
	storage.UpdateUserPassword(userId, "qwerty")

}

//	func TestNew(t *testing.T) {
//		storage, err := sqlite.New(config.MustLoad().StoragePath)
//		if err != nil {
//			fmt.Errorf("Error")
//		}
//	}
func TestStorage_AddExpression(t *testing.T) {

}
func TestStorage_AddUser(t *testing.T) {

}
func TestStorage_GetExpressionById(t *testing.T) {

}
func TestStorage_GetExpressions(t *testing.T) {

}
