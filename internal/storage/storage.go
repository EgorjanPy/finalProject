package storage

import (
	"finalProject/internal/storage/sqlite"
)

var DataBase, _ = sqlite.New()
