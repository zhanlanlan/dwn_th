package handers

import (
	"dwn_th/services"

	"github.com/gin-gonic/gin"
)

func Who(c *gin.Context) *services.Claim {
	user, ok := c.Get("user")
	if !ok {
		panic("no user stored")
	}

	claim, ok := user.(*services.Claim)
	if !ok {
		panic("assert c.Get(\"user\") as model.User failed")
	}

	return claim
}
