package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type (
	washupModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}
	transformedWashup struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

func init() {
	//open connection
	var err error
	db, err = gorm.Open("mysql", "root:@/washup?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect")
	}
	//migrate database
	db.AutoMigrate(&washupModel{})
}

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1/washup")
	{
		v1.POST("/", createWashup)
		v1.GET("/", fetchAllWashup)
		v1.GET("/:id", fetchSingleWashup)
		v1.PUT("/:id", updateWashup)
		v1.DELETE("/:id", deleteWashup)
	}
	router.Run()
}

func createWashup(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := washupModel{Title: c.PostForm("title"), Completed: completed}
	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "messange": "Todo item created succesfully!", "resourceId": todo.ID})
}

func fetchAllWashup(c *gin.Context) {
	var todos []washupModel
	var _todos []transformedWashup

	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "messange": "no todo found!"})
	}
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, transformedWashup{ID: item.ID, Title: item.Title, Completed: completed})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
}

func fetchSingleWashup(c *gin.Context) {
	var todo washupModel
	todoID := c.Param("id")

	db.First(todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "messange": "not todo found!"})
		return
	}
	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}
	_todo := transformedWashup{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}
func updateWashup(c *gin.Context) {
	var todo washupModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "messange": "Not found todo"})
		return
	}
	db.Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "messange": "todo updated successfully"})
}

func deleteWashup(c *gin.Context) {
	var todo washupModel
	todoID := c.Param("id")

	db.First(&todo, todoID)
}
