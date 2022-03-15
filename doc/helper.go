package doc

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GetAllDocPosts Fetch All Doc data
func GetAllDocPosts() ([]Doc, error) {

	session, err := mgo.Dial(HostAddress)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(Database).C(Collection)
	result := []Doc{}
	err3 := collection.Find(bson.M{}).All(&result)
	if err3 != nil {
		return nil, err
	}

	return result, nil

}

// Fetcharticle Fetches the single article
func Fetcharticle(id string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	session, err := mgo.Dial(HostAddress)
	if err != nil {
		log.Println(err)
		return m, err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(Database).C(Collection)
	result1 := Doc{}

	err3 := collection.Find(bson.M{"_id": id}).One(&result1)
	if err3 != nil {
		fmt.Println("the error", err3)
		return m, err3
	}

	m[id] = result1
	m["category"] = PutTogetherNameAndSlug(result1.Tags, result1.TagsSlug)
	return m, nil
}

// Editdoc Modified the document field values
func Editdoc(P *Doc, id string) error {
	session, err := mgo.Dial(HostAddress)
	if err != nil {
		log.Println("Could not connect to mongo: ", err.Error())
		//return false
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	collection := session.DB(Database).C(Collection)
	result1 := Doc{}
	err3 := collection.Find(bson.M{"_id": id}).One(&result1)
	if err3 != nil {
		log.Println("the error", err3)
		return err3
	}
	colQuerier := bson.M{"_id": id}
	log.Println(colQuerier)
	change := bson.M{"$set": bson.M{"_id": P.Slug,

		"title":          P.Title,
		"description":    P.Description,
		"date_submitted": P.DateSubmitted,
		"category":       P.Category,
		"body":           P.Body}}
	log.Println(change)
	err1 := collection.Update(colQuerier, change)
	if err1 != nil {
		log.Println(err1, "here error")
		return err1
	}
	return nil
}

// GetDocPostsCount Count Total Document Count
func GetDocPostsCount() int {

	session, err := mgo.Dial(HostAddress)
	if err != nil {
		return 0
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(Database).C(Collection)
	result := []Doc{}
	err3 := collection.Find(bson.M{}).All(&result)
	if err3 != nil {
		return 0
	}

	return len(result)
}

// SearchByCat search doc by cat
func SearchByCat(category string) map[int]interface{} {
	mapd := make(map[int]interface{})
	session, err := mgo.Dial(HostAddress)
	if err != nil {
		return mapd
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	result := []Doc{}
	collection := session.DB(Database).C(Collection)
	err3 := collection.Find(bson.M{}).All(&result)

	if err3 != nil {
		return mapd
	}

	for i := 0; i < len(result); i++ {
		if result[i].Category == category {

			mapd[i] = result[i]
		}
	}

	return mapd

}

type SubCategory struct {
	Name string
	Slug string
}

// PutTogetherNameAndSlug puts name and slug together
func PutTogetherNameAndSlug(names, slugs []string) []SubCategory {
	var subcategories []SubCategory
	if len(names) != len(slugs) {
		return subcategories
	}
	for i := 0; i < len(names); i++ {
		subcategories = append(subcategories, SubCategory{Name: names[i], Slug: slugs[i]})
	}
	return subcategories
}

// FetchArticlebyTag Fetch articles with same tags
func FetchArticlebyTag(id string) ([]Doc, error) {
	session, err := mgo.Dial(HostAddress)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	queryMap := make(map[string]interface{})

	jsonQueryStr := `{"tags_slug":{"$in":["` + id + `"]}}`
	json.Unmarshal([]byte(jsonQueryStr), &queryMap)
	var doc []Doc
	query := session.DB(Database).C(Collection).Find(queryMap)

	query.All(&doc)

	return doc, nil
}

// GetAllTag returns all the tags
func GetAllTag() map[string]string {

	session, err := mgo.Dial(HostAddress)
	if err != nil {

	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	result := []Doc{}
	collection := session.DB(Database).C(Collection)
	err3 := collection.Find(bson.M{}).All(&result)

	if err3 != nil {

	}
	ListSlug := make(map[string]string)
	for i := 0; i < len(result); i++ {
		for j, _ := range result[i].Tags {
			ListSlug[strings.Trim(result[i].TagsSlug[j], " ")] = strings.Trim(result[i].Tags[j], " ")
		}

	}

	// for i := 0; i < len(result); i++ {
	// 	mapd[i] = PutTogetherNameAndSlug(List[i], ListSlug[i])
	// }

	return ListSlug
}

// SearchByCat2 search doc by cat
func SearchByCat2(category string) map[int]interface{} {
	mapd := make(map[int]interface{})
	session, err := mgo.Dial(HostAddress)
	if err != nil {
		return mapd
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	result := []Doc{}
	collection := session.DB(Database).C(Collection)
	err3 := collection.Find(bson.M{}).All(&result)

	if err3 != nil {
		return mapd
	}

	limit := 0

	for i := 0; i < len(result); i++ {
		if result[i].Category == category {
			if limit < 5 {
				mapd[i] = result[i]
				limit++
			}
		}
	}

	return mapd

}
