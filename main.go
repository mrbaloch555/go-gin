package main

import (
	"context"
	"encoding/xml"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mrbaloch555/go-gin/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func IndexHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"messgae": "Hello world",
	})
}

func IndexHanlderByParam(ctx *gin.Context) {
	name := ctx.Params.ByName("name")
	ctx.JSON(200, gin.H{
		"messgae": "Hello " + name,
	})
}

type Person struct {
	XMLName   xml.Name `xml:"person"`
	FirstName string   `xml:"firstName,attr"`
	LastName  string   `xml:"lastName,attr"`
}

func SendXMLHandler(ctx *gin.Context) {
	ctx.XML(200, Person{
		FirstName: "Durrah",
		LastName:  "Khan",
	})
}

// var recipes []Recipe

// func NewRecipeHandler(ctx *gin.Context) {
// 	var recipe Recipe

// 	if err := ctx.ShouldBindJSON(&recipe); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	recipe.ID = primitive.NewObjectID().String()
// 	recipe.PublishedAt = time.Now()
// 	recipes = append(recipes, recipe)

// 	_, err = collection.InsertOne(ctx, recipe)

// 	if err != nil {
// 		fmt.Println(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error while inserting a new recipe",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, recipe)
// }

// func GetRecipesHandler(ctx *gin.Context) {
// 	// ctx.JSON(200, gin.H{
// 	// 	"recipes": recipes,
// 	// })

// 	cur, err := collection.Find(ctx, bson.M{})

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	defer cur.Close(ctx)

// 	recipes := make([]Recipe, 0)

// 	for cur.Next(ctx) {
// 		var recipe Recipe
// 		cur.Decode(&recipe)
// 		recipes = append(recipes, recipe)
// 	}

// 	ctx.JSON(http.StatusOK, recipes)
// }

// func GetSingleRecipeHandler(ctx *gin.Context) {
// 	id := ctx.Param("id")

// 	objectId, _ := primitive.ObjectIDFromHex(id)
// 	var recipe Recipe
// 	err := collection.FindOne(ctx, bson.M{
// 		"_id": objectId,
// 	}).Decode(&recipe)
// 	if err != nil {
// 		fmt.Println(err)
// 		ctx.JSON(http.StatusInternalServerError,
// 			gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, recipe)
// }

// func SearchRecipesHandler(c *gin.Context) {
// 	tag := c.Query("tag")
// 	listOfRecipes := make([]Recipe, 0)
// 	for i := 0; i < len(recipes); i++ {
// 		found := false
// 		for _, t := range recipes[i].Tags {
// 			if strings.EqualFold(t, tag) {
// 				found = true
// 			}
// 		}
// 		if found {
// 			recipes = append(recipes,
// 				recipes[i])
// 		}
// 	}
// 	c.JSON(http.StatusOK, listOfRecipes)
// }

// func UpdateRecipeHandler(ctx *gin.Context) {
// 	id := ctx.Params.ByName("id")
// 	var recipe Recipe

// 	if err := ctx.ShouldBindJSON(&recipe); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 	}

// 	objectId, _ := primitive.ObjectIDFromHex(id)

// 	_, err = collection.UpdateOne(ctx, bson.M{
// 		"_id": objectId,
// 	}, bson.D{{"$set", bson.D{
// 		{"name", recipe.Name},
// 		{"instructions", recipe.Instructions},
// 		{"ingredients", recipe.Ingredients},
// 		{"tags", recipe.Tags},
// 	}}})
// 	if err != nil {
// 		fmt.Println(err)
// 		ctx.JSON(http.StatusInternalServerError,
// 			gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, recipe)

// }

// func DeleteRecipeHandler(ctx *gin.Context) {
// 	id := ctx.Param("id")

// 	objecId, _ := primitive.ObjectIDFromHex(id)

// 	_, err = collection.DeleteOne(ctx, bson.M{
// 		"_id": objecId,
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 		ctx.JSON(http.StatusInternalServerError,
// 			gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"messgae": "Recipe has been deleted",
// 	})

// }

var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection
var recipesHandler *handlers.RecipeHandler

func init() {
	// recipes = make([]Recipe, 0)
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Database connected")

	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection)
}
func main() {

	router := gin.Default()
	// router.GET("/", IndexHandler)
	// router.GET("/:name", IndexHanlderByParam)
	// router.GET("/xml", SendXMLHandler)
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.GetRecipesHandler)
	router.GET("/recipes/:id", recipesHandler.GetSingleRecipeHandler)
	router.PATCH("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	// router.GET("/recipes/search", SearchRecipesHandler)
	router.Run(":3000")
}
