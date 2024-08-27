// Delivery/routers/router.go
package routers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    setUpAuthRoutes(router)
    setUpUserRoutes(router)
    setUpAdminRoutes(router)
    setUpLoanRoutes(router)
    


   
    return router
}

