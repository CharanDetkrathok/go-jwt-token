package main

import (
	"context"
	"fmt"
	"net/http"
	"survey/controller"
	middleware "survey/middlewere"
	"time"

	"survey/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})

	err := rdb.SetNX(ctx, "test", "redis-cache", 20*time.Second).Err()
	if err != nil {
		panic(err)
	}

	value, err := rdb.Get(ctx,"test").Result()
	if err != nil {
		panic(err)
	}
	
	fmt.Println("Value of Redis cache", value)

	var loginService service.LoginService = service.StaticLoginService()
	var jwtService service.JWTService = service.AuthService()
	var loginController controller.LoginController = controller.LoginHandler(loginService, jwtService)

	server := gin.New()

	server.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

    //จัดกลุ่มการทำงานเพื่อง่ายต่อการจัดการ
	student := server.Group("/student")
	//ดักจับทุก Request ที่ต้องการใช้ Resource ของ PATH URL(/student)
	student.Use(middleware.AuthorizeJWT())
	fmt.Println(student)
	{
		student.GET("/students", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message":"success"})
		})
	}

	port := "8888"
	server.Run(":" + port)
}
