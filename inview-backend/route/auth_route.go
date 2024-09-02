package route

import (
	"insense-local/config"
	"insense-local/database"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRouter(env *config.Env, timeout time.Duration, db database.Database, group *gin.RouterGroup) {
	//TODO
}
