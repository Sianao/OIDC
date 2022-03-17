package controller

import (
	"github.com/gin-gonic/gin"
)

func Client(c *gin.Context) {

	//req, err := http.NewRequest("GET", "http://localhost:8080/oauth_login", nil)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//url := req.URL.Query()
	//url.Add("response_type", "code")
	//url.Add("client_id", "id123")
	//state := strconv.FormatInt(time.Now().Unix(), 10)
	//url.Add("state", state)
	//url.Add("redirect_url", "http://localhost:8080/callback")
	//url.Add("scope", "openid")
	//req.URL.RawQuery = url.Encode()
	//resp, err := http.DefaultClient.Do(req)
	//c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/login?client_id=id123&state="+state)
	//resp.Body.Close()

	return
}
