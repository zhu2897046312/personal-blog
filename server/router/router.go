package router

import (
	"github.com/gin-gonic/gin"
	"github.com/personal-blog/handler"
	"github.com/personal-blog/middleware"
	"github.com/personal-blog/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter initializes the router and sets up all routes
func SetupRouter(factory service.Factory) *gin.Engine {
	r := gin.Default()

	// Add middleware
	r.Use(middleware.CORSMiddleware())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create handlers
	userHandler := handler.NewUserHandler(factory.GetUserService())
	postHandler := handler.NewPostHandler(factory.GetPostService())
	categoryHandler := handler.NewCategoryHandler(factory.GetCategoryService())
	tagHandler := handler.NewTagHandler(factory.GetTagService())

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public routes
		// User routes
		users := v1.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
		}

		// Post routes (public)
		posts := v1.Group("/posts")
		{
			posts.GET("", postHandler.List)                     // 获取文章列表
			posts.GET("/:id", postHandler.Get)                  // 获取文章详情
			posts.GET("/:id/tags", tagHandler.GetPostTags) // 获取文章标签
		}

		// Category routes (public)
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.List)    // 获取分类列表
			categories.GET("/:id", categoryHandler.Get) // 获取分类详情
		}

		// Tag routes (public)
		tags := v1.Group("/tags")
		{
			tags.GET("", tagHandler.List)    // 获取标签列表
			tags.GET("/:id", tagHandler.Get) // 获取标签详情
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			// User routes (authenticated)
			authUsers := protected.Group("/users")
			{
				authUsers.GET("/profile", userHandler.GetProfile)                // 获取个人信息
				authUsers.PUT("/profile", userHandler.UpdateProfile)            // 更新个人信息
				authUsers.PUT("/password", userHandler.ChangePassword)          // 修改密码
				authUsers.GET("", middleware.AdminAuthMiddleware(), userHandler.ListUsers) // 获取用户列表（管理员）
			}

			// Post routes (authenticated)
			authPosts := protected.Group("/posts")
			{
				authPosts.POST("", postHandler.Create)      // 创建文章
				authPosts.PUT("/:id", postHandler.Update)   // 更新文章
				authPosts.DELETE("/:id", postHandler.Delete) // 删除文章
			}

			// Category routes (admin only)
			authCategories := protected.Group("/categories")
			authCategories.Use(middleware.AdminAuthMiddleware())
			{
				authCategories.POST("", categoryHandler.Create)      // 创建分类
				authCategories.PUT("/:id", categoryHandler.Update)   // 更新分类
				authCategories.DELETE("/:id", categoryHandler.Delete) // 删除分类
			}

			// Tag routes (admin only)
			authTags := protected.Group("/tags")
			authTags.Use(middleware.AdminAuthMiddleware())
			{
				authTags.POST("", tagHandler.Create)           // 创建标签
				authTags.POST("/batch", tagHandler.CreateBatch) // 批量创建标签
				authTags.PUT("/:id", tagHandler.Update)        // 更新标签
				authTags.DELETE("/:id", tagHandler.Delete)      // 删除标签
			}
		}
	}

	return r
}
