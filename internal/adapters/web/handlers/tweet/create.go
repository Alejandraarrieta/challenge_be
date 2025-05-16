package tweet

import (
	"challenge_be/internal/usecases/tweet"
	types "challenge_be/pkg/types/tweet/options"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// CreateTweet godoc
// @Summary Crear un tweet
// @Description Crea un nuevo tweet de hasta 280 caracteres
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet body types.InputCreateTweet true "Datos del tweet"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tweets/ [post]
func NewCreateTweetHandler(usecase tweet.CreateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input types.InputCreateTweet
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validate.Struct(input); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			errs := make(map[string]string)
			for _, e := range validationErrors {
				errs[e.Field()] = "El campo es requerido"
			}
			c.JSON(http.StatusBadRequest, gin.H{"validation_errors": errs})
			return
		}

		if len(input.Content) > 280 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El tweet no puede superar los 280 caracteres"})
			return
		}
		if err := usecase.Execute(c.Request.Context(), input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tweet"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Tweet created successfully"})
	}
}
