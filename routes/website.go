package routes

import (
	"encoding/json"
	"fmt"
	htmltempl "html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"l2-golang-auth/config"
	"l2-golang-auth/doc"
	"l2-golang-auth/src/ginjwt"
	"l2-golang-auth/src/noroute"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

var (
	HostAddress string = os.Getenv("HOST")
	PortalLink  string = os.Getenv("PORTALREDIRECT")
	logOut      string = os.Getenv("LOGOUT")
)

func Setuprouterwebsite(v4 *gin.Engine) {

	{
		v4.GET("/login", checklogin, login)
		v4.POST("/login", loginResp)
		v4.GET("/logout", logout)
		v4.GET("/", checklogin, index)
		v4.GET("/article", checklogin, article)
		v4.GET("/article/:slugs", checklogin, singleArticle)
		v4.POST("/search", checklogin, searchPost)
		v4.GET("/articles/tag/:slug", checklogin, searchTag)
	}

}

type accounts struct {
	Profile struct {
		Fname string `json:"fname"`
		Lname string `json:"lname"`
	} `json:"profile"`
}

func checklogin(c *gin.Context) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys["username"] = ""
	logss, err := c.Request.Cookie("accessXtdToken")
	if err != nil {
		fmt.Println(err)
		c.Redirect(302, HostAddress)
		return
	}
	if logss.Value == "" {
		c.Redirect(302, HostAddress)
		return
	} else {
		req, err := http.NewRequest("GET", PortalLink+"api/auth/v1/myaccount", nil)
		if err != nil {
			fmt.Println(err)
			c.Redirect(302, HostAddress)
			return
		}
		token := logss.Value
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			fmt.Println(err)
			c.Redirect(302, HostAddress)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var m accounts
		json.Unmarshal(body, &m)

		c.Keys["username"] = m.Profile.Fname + " " + m.Profile.Lname
	}

}

// func username(c *gin.Context) string {
// 	name, err := c.Request.Cookie("xtendUserName")
// 	if err != nil || name.Value == "" {
// 		c.Redirect(302, "http://dup.vikashkhichar.xyz:39142/login?continueto=http://nipun.vikashkhichar.xyz:8080")
// 		return ""
// 	}
// 	return name.Value

// }

func searchTag(c *gin.Context) {
	slug := c.Param("slug")
	a, err := doc.FetchArticlebyTag(slug)
	if err != nil {
		noroute.NotFoundHandler(c)
	} else {
		c.HTML(200, "articles-category", gin.H{
			"name":   "BROWSE BY TAGS",
			"portal": PortalLink,
			"logout": logOut,
			"names":  a,
		})
	}
}

func conHTML(str string) htmltempl.HTML {
	return htmltempl.HTML(str)
}

func searchPost(c *gin.Context) {
	// Obtain the Posted query values
	// username := username(c)
	searchQuery := c.PostForm("search")
	searchResult, _ := doc.Searchtext(searchQuery)
	c.HTML(200, "/search", gin.H{
		"username": c.Keys["username"],
		"portal":   PortalLink,
		"logout":   logOut,
		"posts":    searchResult,
		"query":    searchQuery,
	})
}

func singleArticle(c *gin.Context) {
	// username := username(c)
	slug := c.Param("slugs")
	BIGDATA := doc.SearchByCat2("Documentation")
	a, err := doc.Fetcharticle(slug)
	if err != nil {
		noroute.NotFoundHandler(c)
	} else {
		c.HTML(200, "article-details", gin.H{
			"username":      c.Keys["username"],
			"names":         a[slug],
			"htmlEscape":    conHTML,
			"portal":        PortalLink,
			"logout":        logOut,
			"related":       BIGDATA,
			"subcategories": a["category"],
		})
	}
}

func article(c *gin.Context) {
	// username := username(c)
	BIGDATA := doc.SearchByCat("Documentation")
	c.HTML(200, "articles-category", gin.H{
		"username": c.Keys["username"],
		"name":     "Documentation",
		"portal":   PortalLink,
		"logout":   logOut,
		"names":    BIGDATA,
	})
}

func index(c *gin.Context) {
	// username := username(c)
	count := doc.GetDocPostsCount()
	data := doc.GetAllTag()
	c.HTML(200, "index", gin.H{
		"username": c.Keys["username"],
		"count":    count,
		"portal":   PortalLink,
		"logout":   logOut,
		"data":     data,
	})
}

func login(c *gin.Context) {
	c.HTML(200, "login", gin.H{
		"title": "login page",
	})
}

func loginResp(c *gin.Context) {
	Name := c.Request.FormValue("emails")
	Pass := c.Request.FormValue("passwords")

	if !config.CheckAdminUseridAndPass(Name, Pass) {
		log.Println("there is some error for the admin")
		c.Redirect(302, "/login")
	} else {
		jwt, _ := ginjwt.GinJwtToken(Name)
		cookie := &http.Cookie{Name: "token",
			Value: jwt["token"].(string)}
		http.SetCookie(c.Writer, cookie)
		c.Redirect(302, "/admin/article/view")
	}
}

func logout(c *gin.Context) {
	k1, _ := c.Request.Cookie("token")
	cookie := &http.Cookie{Name: "token",
		Value: k1.Value, Expires: time.Now()}
	http.SetCookie(c.Writer, cookie)
	c.Redirect(302, "/login")

}

func CreateMyRenderwebsite(r multitemplate.Render) multitemplate.Render {

	r.AddFromFiles("/search", "templates/search-results.tmpl")
	r.AddFromFiles("login", "templates/login.tmpl")
	r.AddFromFiles("404-not-found", "templates/page-not-found.tmpl")
	r.AddFromFiles("index", "templates/index.tmpl")
	r.AddFromFiles("articles-category", "templates/articles-category.tmpl")
	r.AddFromFiles("article-details", "templates/article-details.tmpl")
	return r
}
