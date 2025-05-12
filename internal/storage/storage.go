package storage

import (
	"finalProject/internal/config"
	"finalProject/internal/storage/sqlite"
)

var cfg = config.MustLoad()
var DataBase, _ = sqlite.New(cfg.StoragePath)
