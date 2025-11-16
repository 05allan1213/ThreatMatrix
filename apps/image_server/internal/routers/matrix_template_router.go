package routers

import (
	"image_server/internal/api"
	"image_server/internal/api/matrix_template_api"
	"image_server/internal/middleware"
	"image_server/internal/models"

	"github.com/gin-gonic/gin"
)

func MatrixTemplateRouter(r *gin.RouterGroup) {
	app := api.App.MatrixTemplateApi

	// 矩阵模板创键（POST），绑定 JSON 请求体
	r.POST("matrix_template", middleware.BindJsonMiddleware[matrix_template_api.CreateRequest], app.CreateView)

	// 矩阵模板更新（PUT），绑定 JSON 请求体
	r.PUT("matrix_template", middleware.BindJsonMiddleware[matrix_template_api.UpdateRequest], app.UpdateView)

	// 矩阵模板列表（GET），绑定 Query 参数
	r.GET("matrix_template", middleware.BindQueryMiddleware[models.PageInfo], app.ListView)

	// 矩阵模板选项列表（GET）
	r.GET("matrix_template/options", app.OptionsView)

	// 矩阵模板删除（DELETE），绑定 JSON 请求体
	r.DELETE("matrix_template", middleware.BindJsonMiddleware[models.IDListRequest], app.Remove)

}
