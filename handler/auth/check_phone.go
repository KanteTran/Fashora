package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"fashora-backend/database"
	"fashora-backend/models"
	"fashora-backend/utils"
)

func CheckPhoneNumberExists(c *gin.Context) {
	var input struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	var user models.Users
	err := database.GetDBInstance().DB().Where("phone = ?", input.PhoneNumber).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Phone number does not exist
			utils.SendSuccessResponse(c, http.StatusOK, "Phone number does not exist", gin.H{"exists": false})
			return
		} else {
			// Database query error
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Something went wrong while querying the database")
			return
		}
	}

	// Phone number exists
	utils.SendSuccessResponse(c, http.StatusOK, "Phone number exists", gin.H{"exists": true})
}
