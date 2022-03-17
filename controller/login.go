package controller

import (
	"JD/dao"
	"JD/models"
	"JD/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func Login(c *gin.Context) {
	var u models.Login
	err := c.ShouldBind(&u)
	if err != nil {
		c.JSON(200, gin.H{
			"state": false,
			"msg":   "参数绑定失败",
		})
		return
	}
	Info, err := dao.Login(u)
	if err != nil {
		c.JSON(200, gin.H{
			"state": false,
			"msg":   err.Error(),
		})
		return
	}

	token := utils.MakeToken(*Info)
	t := int(time.Now().Unix())
	utils.SetToken(t, token)
	s, err := utils.HashGet(u.ClientId)
	if err != nil {
		return
	}
	//完成登录 返回参数

	c.Redirect(301, s.RedirectUrl+"?code="+strconv.Itoa(t)+"&state="+s.State)
	return

}

func Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	c.JSON(200, gin.H{
		"state": state,
		"code":  code,
	})
	return
}
func Logout(context *gin.Context) {
	Authorization := context.Request.Header.Get("Authorization")
	ok := utils.DeleteToken(Authorization)
	if !ok {

		context.JSON(200, gin.H{
			"state": false,
			"msg":   "退出登录失败",
		})
		return
	}
	context.JSON(200, gin.H{
		"state": true,
		"msg":   "退出登录成功",
	})
	return
}
func Find(c *gin.Context) {
	var Forget models.Register
	err := c.ShouldBind(&Forget)
	if err != nil {
		c.JSON(200, gin.H{
			"state": false,
			"msg":   "参数绑定失败",
		})
		return
	}
	err = utils.GetConform(Forget.Number, Forget.Code)
	if err != nil {
		c.JSON(200, gin.H{
			"state": false,
			"msg":   err.Error(),
		})
		return
	}
	err = dao.Find(Forget)
	if err != nil {
		c.JSON(200, gin.H{
			"state": false,
			"msg":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"state": true,
		"msg":   "密码找回成功",
	})
	return
}
func Oauth(c *gin.Context) {
	grantType := c.Query("grant_type")
	code := c.Query("code")
	clientId := c.Query("client_id")
	clientCecret := c.Query("client_secret")
	redirectUrl := c.Query("redirect_url")

	s, err := utils.HashGet(clientId)
	if grantType != "authorization_code" {
		c.JSON(200, gin.H{
			"state": false,
			"msg":   "参数错误",
		})
	}
	if err != nil {
		return
	}
	if s.RedirectUrl != redirectUrl || clientCecret != "hello" {
		return
	}
	tokens, err := utils.GET(code)
	if err != nil {
		return
	}
	info, err := utils.ParseToken(tokens)
	iDtoken, err := dao.IdInfo(info.Uid)
	returnToken := utils.MakeToken(iDtoken)
	if err != nil {
		c.JSON(200, gin.H{
			"state": false,
			"err":   err.Error(),
		})
		return
	}
	///没有重定向到里面去 就在这把数据返回了
	c.JSON(200, gin.H{
		"access_token": tokens,
		"id_token":     returnToken,
	})

}

//func Callback(c *gin.Context) {
//	code := c.Query("code")
//	fmt.Println(code)
//	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", "a3112bb967a7bbe3bcf1", "82625129d028e98a671c52c81bd9e45b4b574705", code)
//	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
//	if err != nil {
//
//	}
//	req.Header.Set("Accept", "application/json")
//	httpClient := http.Client{}
//	res, err := httpClient.Do(req)
//	defer res.Body.Close()
//	info, _ := ioutil.ReadAll(res.Body)
//	var token models.Token
//	err = json.Unmarshal(info, &token)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(token)
//	req, err = http.NewRequest("GET", "https://api.github.com/user", nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//	req.Header.Set("Authorization", "token "+token.AccessToken)
//	res, err = httpClient.Do(req)
//	defer res.Body.Close()
//	info, _ = ioutil.ReadAll(res.Body)
//	//fmt.Println(string(info))
//	var basicinfo models.HubBasicInfo
//	err = json.Unmarshal(info, &basicinfo)
//	if err != nil {
//
//	}
//	user, err := dao.HubLogin(basicinfo)
//	if err != nil {
//		c.JSON(200, gin.H{
//			"state": false,
//			"err":   err,
//		})
//		return
//	}
//	newtoken := utils.MakeToken(user)
//	ok := utils.SetToken(newtoken)
//	if !ok {
//		c.JSON(200, gin.H{
//			"state": false,
//		})
//		return
//	}
//	c.JSON(200, gin.H{
//		"state": true,
//		"msg":   "登录成功",
//		"token": newtoken,
//	})
//	return
//}
