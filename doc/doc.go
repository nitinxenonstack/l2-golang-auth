package doc

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"l2-golang-auth/src/utils"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	HostAddress string = os.Getenv("mongoDB_address")
	Database    string = os.Getenv("mongoDB_database")
	Collection  string = os.Getenv("mongoDB_collection")
)

// Doc structure
type Doc struct {
	Slug          string `json:"id,omitempty" bson:"_id,omitempty" form:"id"`
	Title         string `json:"title" binding:"required" form:"title"`
	Description   string `json:"description" binding:"required" form:"description"`
	Body          string `json:"body" binding:"required" form:"body"`
	Tags          []string
	TagsSlug      []string  `bson:"tags_slug"`
	DateSubmitted time.Time `bson:"date_submitted"`
	MetaKeyword   []string  `bson:"meta_keyword"`
	Category      string    `json:"category" binding:"required" form:"category"`
}

// Docnewpage Create New Doc
func Docnewpage(c *gin.Context) {
	c.HTML(200, "admin/article/new", gin.H{})
}

// Docnewpost Create New Doc Post
func Docnewpost(c *gin.Context) {
	DocTitle := c.Request.FormValue("name")
	DocSlug := c.Request.FormValue("_id")
	DocDescription := c.Request.FormValue("description")
	DocBody := c.Request.FormValue("body")
	Subcategory := c.Request.FormValue("subcategory")
	DocCategory := c.Request.FormValue("category")
	Docdate := c.Request.FormValue("date")

	Tags := strings.Split(Subcategory, ",")
	Tagslug := make([]string, len(Tags))

	for i := range Tags {
		Tagslug[i] = Tags[i]
		Tagslug[i] = utils.SlugOfName(Tagslug[i])
	}

	layout := "2006-01-02"
	time, _ := time.Parse(layout, Docdate)
	MetaKeyWord := Tags
	session, err := mgo.Dial(HostAddress)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	DocData := Doc{Title: DocTitle, Slug: DocSlug, Description: DocDescription, Body: DocBody, Tags: Tags, TagsSlug: Tagslug, MetaKeyword: MetaKeyWord, DateSubmitted: time, Category: DocCategory}

	err = session.DB(Database).C(Collection).Insert(DocData)
	if err != nil {
		fmt.Println("the error", err)
	}

	c.Redirect(302, "/admin/article/view")
}

// Viewdoc Return All Doc Data
func Viewdoc(c *gin.Context) {
	x, _ := GetAllDocPosts()
	c.HTML(200, "admin/article/view", gin.H{
		"names": x,
	})
}

// Deletedoc Delete the Selected Data
func Deletedoc(c *gin.Context) {
	slug := c.Param("slug")
	session, err := mgo.Dial(HostAddress)

	if err != nil {
		log.Println("Could not connect to mongo: ", err.Error())
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	result := []Doc{}
	collection := session.DB(Database).C(Collection)
	err3 := collection.Find(bson.M{}).All(&result)

	if err3 != nil {
		fmt.Println("the error", err3)
	}

	for i := 0; i < len(result); i++ {
		if result[i].Slug == slug {
			err = collection.Remove(bson.M{"_id": result[i].Slug})
			log.Println(err)
		}
	}
	c.Redirect(302, "/admin/article/view")
}

// Doceditpage Edit Page
func Doceditpage(c *gin.Context) {
	slug := c.Param("slug")
	a, err := Fetcharticle(slug)
	if err != nil {
		log.Println("the err", err)
	}
	c.HTML(200, "article-edit-page", gin.H{
		"names":       a[slug],
		"stringsJoin": strings.Join,
	})
	c.Next()
}

// Doceditpagepost Edit the Document
func Doceditpagepost(c *gin.Context) {
	slug := c.Param("slug")

	DocTitle := c.Request.FormValue("name")
	DocSlug := c.Request.FormValue("_id")
	DocDescription := c.Request.FormValue("description")
	DocBody := c.Request.FormValue("body")
	DocCategory := c.Request.FormValue("category")
	Docdate := c.Request.FormValue("date")
	Subcategory := c.Request.FormValue("subcategory")
	Tags := strings.Split(Subcategory, ",")
	Tagslug := make([]string, len(Tags))

	layout := "2006-01-02"
	time, _ := time.Parse(layout, Docdate)

	for i := range Tags {
		Tagslug[i] = Tags[i]
		Tagslug[i] = utils.SlugOfName(Tagslug[i])
	}
	MetaKeyWord := Tags
	DocData := Doc{Title: DocTitle,
		Slug:          DocSlug,
		Description:   DocDescription,
		Body:          DocBody,
		DateSubmitted: time,
		MetaKeyword:   MetaKeyWord,
		Category:      DocCategory}
	Editdoc(&DocData, slug)

	c.Redirect(302, "/admin/article/view")
}

// Searchtext doc
func Searchtext(str string) ([]Doc, error) {
	var result []Doc
	session, err := mgo.Dial(HostAddress)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	collection := session.DB(Database).C(Collection)
	errs := collection.Find(bson.M{"$text": bson.M{"$search": str}}).All(&result)
	if errs != nil {
		return nil, err
	}

	return result, nil
}
