package controller

import (
	"JD/models"
	"JD/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Service(c *gin.Context) {
	var re models.Request
	err := c.ShouldBindQuery(&re)
	if err != nil {
		return
	}
	utils.HashSet(re)
	c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/login?client_id="+re.ClientId+"&state="+re.State)
}
func RedirectService() {

}
