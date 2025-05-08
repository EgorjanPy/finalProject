package sqlite_test

import (
	"finalProject/internal/config"
	"finalProject/internal/storage/sqlite"
	"fmt"
	"testing"
)

func TestDataBase(t *testing.T) {
	//t.Run("New", TestNew)
	_, err := sqlite.New(config.MustLoad().StoragePath)
	if err != nil {
		fmt.Errorf("Error %v", err)
		//return
	}
	//id, err := storage.AddUser("egor", "1234")
	//if err != nil {
	//	fmt.Errorf("Error %v", err)
	//	return
	//}
	//exID, err := storage.AddExpression(&sqlite.Expression{
	//	ID:         id,
	//	Expression: "2+2",
	//})
	//if err != nil {
	//	fmt.Errorf("Error %v", err)
	//	return
	//}
	//ex := storage.GetExpressionById(exID, id)
	//fmt.Println(ex.Expression)
	//err = storage.SetResult(exID, "4")
	//if err != nil {
	//	fmt.Errorf("Error %v", err)
	//	return
	//}
	//storage.UpdateUserPassword(id, "qwerty")
	//storage.UpdateUserPassword(id, "qwerty")

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
