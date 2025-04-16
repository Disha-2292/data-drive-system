package routes

import (
	"github.com/Disha-2292/data-drive-system/controllers"
	"github.com/Disha-2292/data-drive-system/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Auth routes for login and register
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Protected routes requiring authentication
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware()) // Auth middleware to ensure the user is logged in
	{
		protected.GET("/me", func(c *gin.Context) {
			userID := c.MustGet("user_id").(uint)
			c.JSON(200, gin.H{"user_id": userID})
		})

		// File routes protected by RBAC
		files := protected.Group("/files")
		{
			// Admin-only routes for file/folder management (Create, Update, Delete)
			files.Use(middleware.AdminOnly())
			{
				files.POST("/create", controllers.CreateFileOrFolder) // Admin can create files/folders
				files.PUT("/:id", controllers.UpdateFile)             // Admin can update files
				files.DELETE("/:id", controllers.DeleteFile)          // Admin can delete files
				files.GET("/all", controllers.GetAllFiles)            // Admin can see all files
			}

			// User-specific routes for file management
			files.Use(middleware.UserOnly()) // User-only middleware for managing their own files
			{
				files.GET("/", controllers.GetUserFiles)             // Users can see their own files
				files.GET("/:id", controllers.GetFileByID)           // Users can get their file by ID
				files.POST("/upload", controllers.UploadFile)        // Users can upload files
				files.GET("/download/:id", controllers.DownloadFile) // Users can download their files
			}

			// New Routes for Sharing and Permissions Management
			// Share a file with other users
			files.POST("/:id/share", controllers.ShareFile) // User can share a file with others

			// Check file permissions
			files.GET("/:id/permissions", controllers.GetFilePermissions) // Get permissions for a file

			// Update file permissions
			files.PUT("/:id/permissions", controllers.UpdateFilePermissions) // Update permissions for a shared file

			files.GET("/search", controllers.SearchFiles)

		}
	}
}
