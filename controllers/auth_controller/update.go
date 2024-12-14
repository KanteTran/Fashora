package auth_controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"fashora-backend/models"
	"fashora-backend/services/auth_service"
	"fashora-backend/services/user_service"
	"fashora-backend/utils"
)

func Update(c *gin.Context) {
	var input models.UserInfo
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	user, err := auth_service.GetAuthenticatedUser(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if user.Phone != input.PhoneNumber {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid Token")
		return
	}

	if err := user_service.UpdateUserByPhoneNumber(input); err != nil {
		handleUpdateError(c, err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "User updated successfully", nil)
}

func handleUpdateError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.SendErrorResponse(c, http.StatusNotFound, "User does not exist")
	} else {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update user")
	}
}
