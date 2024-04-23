package handler

import (
	"net/http"

	"github.com/SergeyButGusaim/go_final_project-main/pkg/service"
	"github.com/gin-gonic/gin"
)

const DATE_FORMAT = "20060102"

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) Handler {
	return Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	/*router.GET("/api/nextdate", h.getNextDate)
	router.POST("/task", h.createTask)
	router.GET("/task", h.getTaskById)*/
	api := router.Group("/api")
	{
		api.GET("/nextdate", h.getNextDate)

		api.POST("/task", h.createTask)
		api.GET("/task", h.getTaskById)
	}

	static := router.Group("/")
	{
		static.StaticFS("./css", http.Dir("./web/css"))
		static.StaticFS("./js", http.Dir("./web/js"))
	}

	router.GET("/", h.indexPage)
	router.StaticFile("/index.html", "./web/index.html")
	router.StaticFile("/login.html", "./web/login.html")
	router.StaticFile("/favicon.ico", "./web/favicon.ico")

	return router
}

func (h *Handler) indexPage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "./web/index.html")
}
