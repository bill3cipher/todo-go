package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//Model for respresent table and send data to front-end
type Todo struct {
	//frontend sent text
	Title string `json:"text"`
	gorm.Model
}

//manual define tablename
func (Todo) TableName() string {
	return "todos"
}

//dependency gormDB
type TodoHandler struct {
	db *gorm.DB
}

// for frontend
func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

func (t *TodoHandler) NewTask(c *gin.Context) {
	// s := c.Request.Header.Get("Authorization")
	// tokenString := strings.TrimPrefix(s, "Bearer ")

	// if err := auth.Protect(tokenString); err != nil {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }

	var todo Todo
	//shouldblind manual error handle
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	r := t.db.Create(&todo)
	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ID": todo.Model.ID,
	})
}
