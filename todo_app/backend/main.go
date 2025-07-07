package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var (
	todos  = make(map[int]Todo)
	nextID = 1
	mu     sync.Mutex
)

func main() {
	r := gin.Default()

	r.GET("/todos", getTodos)
	r.POST("/todos", addTodo)
	r.PUT("/todos/:id", updateTodo)
	r.DELETE("/todos/:id", deleteTodo)

	r.Run(":8092")
}

func getTodos(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	var todoList []Todo
	for _, todo := range todos {
		todoList = append(todoList, todo)
	}

	c.JSON(http.StatusOK, todoList)
}

func addTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	newTodo.ID = nextID
	todos[nextID] = newTodo
	nextID++

	c.JSON(http.StatusCreated, newTodo)
}

func updateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedTodo Todo
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := todos[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	updatedTodo.ID = id
	todos[id] = updatedTodo

	c.JSON(http.StatusOK, updatedTodo)
}

func deleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := todos[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	delete(todos, id)
	c.Status(http.StatusNoContent)
}
