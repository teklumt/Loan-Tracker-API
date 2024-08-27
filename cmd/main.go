package main

import (
	"loan-tracker-api/config"
	"loan-tracker-api/config/db"
	"loan-tracker-api/delivery/routers"
)

func main() {
    config.InitiEnvConfigs() 
    db.ConnectDB(config.EnvConfigs.MongoURI)
    //printGreen color 



    
    router := routers.SetupRouter()

    router.Run(config.EnvConfigs.LocalServerPort)
}
