package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type Note struct {
	gorm.Model
	Question string `json:"question" gorm:"not null;default:null"`
	Answer   string `json:"answer" gorm:"not null;default:null"`
}

type ErrorMessage struct {
	ErrocCode uint
	Message   string
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/note", GetNotes)
	router.GET("/note/", GetNotes)
	router.GET("/note/:id", GetNote)
	router.GET("/note/:id/", GetNote)
	router.POST("/note", CreateNote)
	router.POST("/note/", CreateNote)
	router.PUT("note/:id", UpdateNote)
	router.PUT("note/:id", UpdateNote)
	router.DELETE("note/:id", DeleteNote)
	return router
}

func getNoteID(ctx *gin.Context) (int, error) {
	note_id, err := strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		ctx.JSON(
			http.StatusNotFound,
			ErrorMessage{http.StatusNotFound, fmt.Sprintf("Malformed id: %v", err)},
		)
		fmt.Println(err)
		return 0, err
	}
	return note_id, nil
}

func GetNotes(ctx *gin.Context) {
	var notes []Note
	if err := db.Find(&notes).Error; err != nil {
		ctx.JSON(
			http.StatusNotFound,
			ErrorMessage{http.StatusNotFound, fmt.Sprintf("%v", err)},
		)
		fmt.Println(err)
	} else {
		ctx.JSON(http.StatusOK, notes)
	}
}

func GetNote(ctx *gin.Context) {
	note_id, err := getNoteID(ctx)
	if err != nil {
		return
	}
	var note Note
	if err := db.Where("id = ?", note_id).First(&note).Error; err != nil {
		ctx.JSON(
			http.StatusNotFound,
			ErrorMessage{http.StatusNotFound, fmt.Sprintf("%v", err)},
		)
		fmt.Println(err)
	} else {
		ctx.JSON(http.StatusOK, note)
	}
}

func CreateNote(ctx *gin.Context) {
	var new_note Note
	if err := ctx.BindJSON(&new_note); err != nil {
		fmt.Println(err)
		ctx.JSON(
			http.StatusBadRequest,
			ErrorMessage{http.StatusBadRequest, fmt.Sprintf("JSON error: %v", err)},
		)
		return
	}

	if err := db.Create(&new_note).Error; err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			ErrorMessage{http.StatusBadRequest, fmt.Sprintf("DB error: %v", err)},
		)
		fmt.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, new_note)
}

func UpdateNote(ctx *gin.Context) {
	note_id, err := getNoteID(ctx)
	if err != nil {
		return
	}

	var updated_note Note
	if err := ctx.BindJSON(&updated_note); err != nil {
		fmt.Println(err)
		ctx.JSON(
			http.StatusBadRequest,
			ErrorMessage{http.StatusBadRequest, fmt.Sprintf("JSON error: %v", err)},
		)
		return
	}
	updated_note.ID = uint(note_id)

	var existing_note Note
	if err := db.First(&existing_note, updated_note.ID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			ctx.JSON(
				http.StatusBadRequest,
				ErrorMessage{http.StatusBadRequest, fmt.Sprintf("Note with id %v does not exist: %v", note_id, err)},
			)
		} else {
			ctx.JSON(
				http.StatusInternalServerError,
				ErrorMessage{http.StatusInternalServerError, fmt.Sprintf("Unknown DB error: %v", err)},
			)
		}
		return
	}

	if err := db.Save(&updated_note).Error; err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			ErrorMessage{http.StatusBadRequest, fmt.Sprintf("DB error: %v", err)},
		)
		fmt.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, updated_note)
}

func DeleteNote(ctx *gin.Context) {
	note_id, err := strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		ctx.JSON(
			http.StatusNotFound,
			ErrorMessage{http.StatusNotFound, fmt.Sprintf("Malformed id: %v", err)},
		)
		fmt.Println(err)
		return
	}

	var existing_note Note
	if err := db.First(&existing_note, note_id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			ctx.JSON(
				http.StatusBadRequest,
				ErrorMessage{http.StatusBadRequest, fmt.Sprintf("Note with id %v does not exist: %v", note_id, err)},
			)
		} else {
			ctx.JSON(
				http.StatusInternalServerError,
				ErrorMessage{http.StatusInternalServerError, fmt.Sprintf("Unknown DB error: %v", err)},
			)
		}
		return
	}

	if err := db.Delete(&existing_note).Error; err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			ErrorMessage{http.StatusBadRequest, fmt.Sprintf("DB error: %v", err)},
		)
		fmt.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func main() {
	db, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// initialise database
	db.AutoMigrate(&Note{})

	router := SetupRouter()
	router.Run("127.0.0.1:8080")
}
