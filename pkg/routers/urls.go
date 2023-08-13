package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gomessage/pkg/api"
	"gomessage/pkg/api/vClient"
	"gomessage/pkg/authorization"
	"gomessage/pkg/middleware"
	"gomessage/pkg/utils/log/loggers"
	"net/http"
)

func AddStatic(g *gin.Engine) {
	g.StaticFile("/favicon.ico", "./assets/favicon.ico")
	g.Static("/static", "assets/vue/static")
	g.LoadHTMLGlob("assets/vue/*.html")
}

func Path(g *gin.Engine) {

	//加载静态文件
	AddStatic(g)

	//中间件
	g.Use(middleware.Cors())
	g.Use(middleware.AccessLog())

	//健康检测
	g.GET("/ok", api.Hello)
	g.GET("/health", api.Health)

	//首页
	g.GET("/", api.Index)

	//swagger文档
	g.GET("/docs", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/swagger/index.html") })
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) //Swagger页面

	//数据入口Get方法
	g.GET("/go/:namespace", middleware.CheckNamespace(), api.GoMessageByGet)                             //给单个路由追加中间件middleware.CheckNamespace()
	g.GET("/go", func(c *gin.Context) { c.Request.URL.Path = "/go/message"; g.HandleContext(c) })        //把"/go"重定向到"/go/message"的路由上
	g.GET("/gomessage", func(c *gin.Context) { c.Request.URL.Path = "/go/message"; g.HandleContext(c) }) //把"/gomessage"重定向到"/go/message"的路由上
	//数据入口Post方法
	g.POST("/go/:namespace", middleware.CheckNamespace(), api.GoMessageByTransport)                       //给单个路由追加中间件middleware.CheckNamespace()
	g.POST("/go", func(c *gin.Context) { c.Request.URL.Path = "/go/message"; g.HandleContext(c) })        //把"/go"重定向到"/go/message"的路由上
	g.POST("/gomessage", func(c *gin.Context) { c.Request.URL.Path = "/go/message"; g.HandleContext(c) }) //把"/gomessage"重定向到"/go/message"的路由上

	//登录注册相关
	g.POST("/auth/login", authorization.Login)
	g.POST("/auth/logout", authorization.Logout)
	g.POST("/auth/register", authorization.Register)

	//=======================
	// 用户操作：v1版本
	//=======================
	v1View := g.Group("/api/v1")
	v1View.Use(middleware.CheckNamespace())
	v1View.Use(middleware.AuthMiddleware())
	{
		//命名空间健康检测
		v1View.GET("/:namespace/health", api.Health)

		//数据劫持
		v1View.GET("/:namespace/json", api.GetNamespaceJson)

		//用户变量
		v1View.GET("/:namespace/vars", api.ListVariables)  //获取变量映射
		v1View.POST("/:namespace/vars", api.PostVariables) //添加变量映射

		//TODO: 这个功能虽然后端完成了，但是前端ui层面没有启动这个接口对应的功能，有空时再过来修改一下
		v1View.GET("/:namespace/flattening", api.GetNamespaceFlatteningJson)

		//消息模板
		v1View.GET("/:namespace/template", api.ListTemplate)
		v1View.POST("/:namespace/template", api.PostTemplate)

		//客户端
		v1View.GET("/:namespace/client", vClient.ListClient)             //获取所有客户端
		v1View.POST("/:namespace/client", vClient.PostClient)            //新增一个客户端
		v1View.GET("/:namespace/client/:id", vClient.GetClient)          //获取客户端详情
		v1View.PUT("/:namespace/client/:id", vClient.PutClientActive)    //更新一个客户端
		v1View.PUT("/:namespace/client-info/:id", vClient.PutClientInfo) //更新一个客户端
		v1View.DELETE("/:namespace/client/:id", vClient.DeleteClient)    //删除一个客户端

	}

	//=======================
	// 命名空间操作：v1版本
	//=======================
	v1Namespace := g.Group("/api/v1")
	v1Namespace.Use(middleware.AuthMiddleware())
	{
		//命名空间
		v1Namespace.GET("/namespace", api.ListNamespace)
		v1Namespace.POST("/namespace", api.PostNamespace)
		v1Namespace.GET("/namespace/:id", api.GetNamespace)
		v1Namespace.PUT("/namespace/:id", api.PutNamespace)
		v1Namespace.DELETE("/namespace/:id", api.DeleteNamespace) //删除命名空间的时候，需要把当前命名空间下的（变量映射、模板、客户端）全都删除掉
	}

	loggers.DefaultLogger.Info("路由表加载完成...")

}
