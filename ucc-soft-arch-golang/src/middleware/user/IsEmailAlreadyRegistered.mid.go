package user

import (
	customError "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/errors"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/services"

	"github.com/gin-gonic/gin"
)

func IsEmailAvailable(service services.IUserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.GetString("Email")
		_, err := service.GetUserByEmail(email)
		if err == nil {
			// email already exists -> return a proper error (avoid nil error panic)
			c.Error(customError.NewError("EMAIL_TAKEN", "Email already registered", 409))
			c.Abort()
			return
		}
		c.Next()
	}
}
