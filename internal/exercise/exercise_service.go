package exercise

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/raa11dev/course/internal/domain"
	"gorm.io/gorm"
)

type ExerciseService struct {
	db *gorm.DB
}

func NewExerciseService(db *gorm.DB) *ExerciseService {
	return &ExerciseService{
		db: db,
	}
}

func (ex ExerciseService) CreateAnswer(ctx *gin.Context) {
	Idex := ctx.Param("id")
	idEX, err := strconv.Atoi(Idex)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}

	idQue := ctx.Param("id2")
	idQUE, err := strconv.Atoi(idQue)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid question id",
		})
		return
	}

	token := strings.Split(ctx.Request.Header["Authorization"][0], " ")[1]
	data, _ := ex.DecriptJWT(token)
	tokenString := fmt.Sprintf("%v", data["user_id"])
	tokenInt, _ := strconv.Atoi(tokenString)

	var answer domain.Answer
	err = ctx.ShouldBind(&answer)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	if answer.Answer == "" {
		ctx.JSON(400, gin.H{
			"message": "field Answer required",
		})
		return
	}

	answer.QuestionID = idQUE
	answer.ExerciseID = idEX
	answer.UserID = tokenInt

	if err := ex.db.Create(&answer).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "failed when create answer",
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Succes store question",
		"data":    answer,
	})
}

func (ex ExerciseService) CreateQuestions(ctx *gin.Context) {
	paramID := ctx.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}

	var question domain.Question
	err = ctx.ShouldBind(&question)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	if question.Body == "" {
		ctx.JSON(400, gin.H{
			"message": "field body required",
		})
		return
	}

	if question.OptionA == "" {
		ctx.JSON(400, gin.H{
			"message": "field option a required",
		})
		return
	}

	if question.OptionB == "" {
		ctx.JSON(400, gin.H{
			"message": "field option b required",
		})
		return
	}

	if question.OptionC == "" {
		ctx.JSON(400, gin.H{
			"message": "field option c required",
		})
		return
	}

	if question.OptionD == "" {
		ctx.JSON(400, gin.H{
			"message": "field option d required",
		})
		return
	}

	if question.CorrectAnswer == "" {
		ctx.JSON(400, gin.H{
			"message": "field correct answer required",
		})
		return
	}

	question.ExerciseID = id
	question.Score = 10
	token := strings.Split(ctx.Request.Header["Authorization"][0], " ")[1]
	data, _ := ex.DecriptJWT(token)
	tokenString := fmt.Sprintf("%v", data["user_id"])
	tokenInt, _ := strconv.Atoi(tokenString)
	question.CreatorID = tokenInt

	if err := ex.db.Create(&question).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "failed when create question",
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Succes store question",
		"data":    question,
	})
}

func (ex ExerciseService) CreateExercise(ctx *gin.Context) {
	var exercise domain.Exercise
	err := ctx.ShouldBind(&exercise)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	if exercise.Title == "" {
		ctx.JSON(400, gin.H{
			"message": "field title required",
		})
		return
	}

	if exercise.Description == "" {
		ctx.JSON(400, gin.H{
			"message": "field description required",
		})
		return
	}

	if err := ex.db.Create(&exercise).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "failed when create exercise",
		})
		return
	}
	ctx.JSON(201, gin.H{
		"message": "Succes store exercise",
		"data":    exercise,
	})
}

func (ex ExerciseService) GetExercise(ctx *gin.Context) {
	paramID := ctx.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}
	var exercise domain.Exercise
	err = ex.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}
	ctx.JSON(200, exercise)
}

func (ex ExerciseService) GetUserScore(ctx *gin.Context) {
	paramExerciseID := ctx.Param("id")
	exerciseID, err := strconv.Atoi(paramExerciseID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}
	var exercise domain.Exercise
	err = ex.db.Where("id = ?", exerciseID).Preload("Questions").Take(&exercise).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}

	userID := int(ctx.Request.Context().Value("user_id").(float64))
	var answers []domain.Answer
	err = ex.db.Where("exercise_id = ? AND user_id = ?", exerciseID, userID).Find(&answers).Error

	if err != nil {
		ctx.JSON(200, gin.H{
			"score": 0,
		})
		return
	}

	mapQA := make(map[int]domain.Answer)
	for _, answer := range answers {
		mapQA[answer.QuestionID] = answer
	}

	var score int
	for _, question := range exercise.Questions {
		if strings.EqualFold(question.CorrectAnswer, mapQA[question.ID].Answer) {
			score += question.Score
		}
	}
	ctx.JSON(200, gin.H{
		"score": score,
	})
}

var signatureKey = []byte("mySuperSecretSignature")

func (ex ExerciseService) DecriptJWT(token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("auth invalid")
		}
		return signatureKey, nil
	})

	data := make(map[string]interface{})
	if err != nil {
		return data, err
	}
	if !parsedToken.Valid {
		return data, errors.New("token invalid")
	}
	return parsedToken.Claims.(jwt.MapClaims), nil
}
