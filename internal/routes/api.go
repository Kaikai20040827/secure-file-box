package routes

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/Kaikai20040827/graduation/internal/config"
	"github.com/Kaikai20040827/graduation/internal/handler"
	"github.com/Kaikai20040827/graduation/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

// RegisterAPIRoutes 注册所有 API 路由
// 由 main.go 负责把 handler 和 config 注入进来
func RegisterAPIRoutes(
	r *gin.Engine,
	authH *handler.AuthHandler,
	userH *handler.UserHandler,
	fileH *handler.FileHandler,
	jwtCfg *config.JWTConfig,
) {
	api := r.Group("/api/v1")

	// 公共 API
	{
		api.GET("/ping", handler.Ping)

		// 公共文件上传（无需认证，供前端测试或匿名上传使用）
		api.POST("/files/public/upload", fileH.UploadFilePublic)

		auth := api.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
		}
	}

	// 需要认证
	authRequired := api.Group("")
	authRequired.Use(middleware.JWTAuthMiddleware(jwtCfg))
	{
		// 用户
		authRequired.GET("/user/profile", userH.GetProfile)
		authRequired.PUT("/user/profile", userH.UpdateProfile)
		authRequired.PUT("/user/password", userH.ChangePassword)

		// 文件
		authRequired.POST("/files/upload", fileH.UploadFile)
		authRequired.GET("/files", fileH.ListFiles)
		authRequired.GET("/files/download/:id", fileH.DownloadFile)
		authRequired.DELETE("/files/:id", fileH.DeleteFile)
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://127.0.0.1:8080", "http://localhost:8080"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge: 24 * time.Hour,
    }))
	
	// 获取项目根目录（更可靠的方式）
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	webStaticPath := filepath.Join(projectRoot, "web", "static")
	webTemplatesPath := filepath.Join(projectRoot, "web", "templates")
	webImagesPath := filepath.Join(projectRoot, "web", "static", "images")

	// 静态文件服务
	r.Static("/static", webStaticPath)

	// HTML 页面路由
	r.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(webTemplatesPath, "login.html"))
	})

	r.GET("/register", func(c *gin.Context) {
		c.File(filepath.Join(webTemplatesPath, "signup.html"))
	})

	r.GET("/logo", func(c *gin.Context) {
		c.File(filepath.Join(webImagesPath, "logo.png"))
	})

	r.GET("/index", func(c *gin.Context) {
		c.File(filepath.Join(webTemplatesPath, "index.html"))
	})

	r.GET("/exam", func(c *gin.Context) {
		c.File(filepath.Join(webTemplatesPath, "exam.html"))
	})

	r.GET("/timetable", func(c *gin.Context) {
		c.File(filepath.Join(webTemplatesPath, "timetable.html"))
	})

	r.GET("/password", func(c *gin.Context) {
		c.File(filepath.Join(webTemplatesPath, "password.html"))
	})

	r.GET("/settings", func(c *gin.Context) {
		c.File(filepath.Join(webTemplatesPath, "settings.html"))
	})

	r.GET("/register_result", func(c *gin.Context) {
		c.File(filepath.Join(webTemplatesPath, "register_result.html"))
	})

	// API 路由
	r.GET("/ping", handler.Ping)

	// 可以继续添加其他API路由

	return r
}
