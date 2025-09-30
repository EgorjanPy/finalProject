package storage

import (
	"finalProject/internal/config"
	"finalProject/internal/storage/sqlite"
)

var DataBase, _ = sqlite.New(config.Cfg.StoragePath)
