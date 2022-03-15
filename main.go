package main

import (
	"l2-golang-auth/doc"
	"l2-golang-auth/routes"
	"l2-golang-auth/src/api"
	"l2-golang-auth/src/ginjwt"
	"l2-golang-auth/src/noroute"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {

	// Set the router as the default one provided by Gin
	router = gin.Default()

	router.NoRoute(noroute.NotFoundHandler)
	staticRoutes := router.Group("/")
	staticRoutes.Use(setCacheControlHeaderMw)
	staticRoutes.Static("/admin/static", "./assets/admin/static")
	staticRoutes.Static("/static/default", "./assets/static/default")
	router.LoadHTMLGlob("templates/*.tmpl")

	// Get Doc
	router.GET("/v1/doc_content", api.DocContent)

	// Get Doc with Id
	router.GET("/v1/doc_content/:slug", func(c *gin.Context) {
		name := c.Param("slug")
		api.DocIDContent(c, name)
	})

	// Search Doc Data
	router.GET("/v1/search/:slug", func(c *gin.Context) {
		name := c.Param("slug")
		api.DocIDContent(c, name)
	})

	v1 := router.Group("/")

	authMiddleware := ginjwt.MwInitializer()
	v1.Use(get_cookie, authMiddleware.MiddlewareFunc())
	{
		// Doc Backend
		v1.GET("admin/article/new", doc.Docnewpage)
		v1.POST("admin/article/new", doc.Docnewpost)
		v1.GET("admin/article/view", doc.Viewdoc)
		v1.GET("article-edit-page/:slug", doc.Doceditpage)
		v1.POST("article-edit-page/:slug", doc.Doceditpagepost)
		v1.GET("deletedoc/:slug", doc.Deletedoc)

	}

	template := multitemplate.New()
	extroute := routes.CreateMyRenderwebsite(template)
	introute := CreateMyRender(template)

	addCommonTemplates(extroute)
	router.HTMLRender = extroute
	routes.Setuprouterwebsite(router)
	router.HTMLRender = introute

	router.Run()

}

func addCommonTemplates(r multitemplate.Render) {
	for k := range r {
		tmpl, _ := r[k].ParseGlob("./templates/commons/*.tmpl")
		tmpl.New("popup.tmpl").Parse("<div></div>") // + popup.PopConvert().Body)
		r[k] = tmpl
	}

}

func setCacheControlHeaderMw(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "max-age=2592000, public")
	c.Next()
}

func get_cookie(c *gin.Context) {
	k1, err := c.Request.Cookie("token")
	if err != nil {
		c.Next()
		return
	}
	c.Request.Header.Set("Authorization", "Bearer "+k1.Value)
	c.Next()
}

func CreateMyRender(r multitemplate.Render) multitemplate.Render {

	//
	r.AddFromFiles("admin/article/new", "templates/article-new-page.tmpl")
	r.AddFromFiles("admin/article/view", "templates/article.tmpl")
	r.AddFromFiles("article-edit-page", "templates/article-edit-page.tmpl")

	return r
}
