package main

import (
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/db"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/utils"
)

func main() {
	utils.Logger.Info("Service Started!!")
	defer db.DisconnectDB()

}
