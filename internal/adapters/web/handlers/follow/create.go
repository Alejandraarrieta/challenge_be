package follow

import (
	"challenge_be/internal/usecases/follow"
	types "challenge_be/pkg/types/follow/options"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// CreateFollow godoc
// @Summary Seguir a un usuario
// @Description Crea una relación de follow entre el usuario origen y el destino
// @Tags follows
// @Accept json
// @Produce json
// @Param follow body types.InputCreateFollow true "Datos del follow"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /follows [post]
func NewCreateFollowHandler(usecase follow.CreateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input types.InputCreateFollow
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
		if err := usecase.Execute(c.Request.Context(), input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create follow"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Follow created successfully"})
	}
}
