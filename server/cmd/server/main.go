package main

import (
	"log"
		
	"github.com/personal-blog/config"
	"github.com/personal-blog/database"
	"github.com/personal-blog/repository/mysql"
	"github.com/personal-blog/repository/redis"
	"github.com/personal-blog/router"
	"github.com/personal-blog/service"
	
	// 确保导入你的docs包
	_ "github.com/personal-blog/docs"  
)

// @title Personal Blog API
// @version 1.0
// @description 个人博客系统API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database
	if err := database.InitMySQL(); err != nil {
		log.Fatalf("Error initializing MySQL: %v", err)
	}

	if err := database.InitRedis(); err != nil {
		log.Fatalf("Error initializing Redis: %v", err)
	}

	// Create MySQL factory
	mysqlFactory := mysql.NewFactory(database.DB)

	// Create Redis factory
	redisFactory := redis.NewFactory(database.RedisClient)

	// Create service factory
	factory := service.NewFactory(mysqlFactory, redisFactory)

	// Set up the router
	r := router.SetupRouter(factory)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
