package router

import (
	"github.com/VolkHackVH/todo-list.git/internal/db"
	"github.com/VolkHackVH/todo-list.git/internal/handlers"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func InitRouter(db *db.Queries) *gin.Engine {
	r := gin.Default()

	r.Use(
		gin.Recovery(),
		requestid.New(),
	)

	registerAPIRouter(r, db)

	r.SetTrustedProxies([]string{"127.0.0.1"})

	return r
}

func registerAPIRouter(r *gin.Engine, db *db.Queries) {
	hand := handlers.NewHandler(db)
	api := r.Group("/api")

	//? User routes
	user := api.Group("/user")
	{
		user.POST("/", hand.User.CreateUser)
		user.GET("/:id", hand.User.GetUserInfo)
		user.PUT("/:id", hand.User.UpdateUser)
		user.DELETE("/:id", hand.User.RemoveUser)
	}

	//? Task routes
	task := api.Group("/task")
	{
		task.POST("/", hand.Task.CreateTask)
		task.GET("/:id", hand.Task.GetTaskInfo)
		task.PUT("/:id", hand.Task.UpdateTask)
		task.DELETE("/:id", hand.Task.RemoveTask)
	}
}
