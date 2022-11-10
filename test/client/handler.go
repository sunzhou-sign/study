package main

import (
	gin_session "client/session"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//用户信息
type UserInfo struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

// 编写一个校验用户是否登录的中间件
// 其实就是从上下文中取到session data,从session data取到isLogin
func AuthMiddleware(c *gin.Context) {
	// 1. 从上下文中取到session data
	// 1. 先从上下文中获取session data
	fmt.Println("in Auth")
	tmpSD, _ := c.Get(gin_session.SessionContextName)
	sd := tmpSD.(gin_session.SessionData)
	// 2. 从session data取到isLogin
	fmt.Printf("%#v\n", sd)
	value, err := sd.Get("isLogin")
	if err != nil {
		fmt.Println(err)
		// 取不到就是没有登录
		c.Redirect(http.StatusFound, "/login")
		return
	}
	fmt.Println(value)
	isLogin, ok := value.(bool) //类型断言
	if !ok {
		fmt.Println("!ok")
		c.Redirect(http.StatusFound, "/login")
		return
	}
	fmt.Println(isLogin)
	if !isLogin {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.Next()
}

//这个是最主要的，因此涉及到表单数据的提取，cookie的设置等。
func loginHandler(c *gin.Context) {
	if c.Request.Method == "POST" { //判断请求的方法，先判是否为post
		toPath := c.DefaultQuery("next", "/index") //一个路径，用于后面的重定向
		var u UserInfo
		//绑定，并解析参数
		err := c.ShouldBind(&u)
		if err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"err": "用户名或密码不能为空",
			})
			return
		}
		//解析成功
		//验证输入的账号密码受否正确
		//这里再生产中应该去数据区取信息进行比对,但这里直接写死
		if u.Username == "zhouzheng" && u.Password == "123" {
			//接下来是核心代码
			//验证成功，,在当前sessiondata设置islogin=true
			// 登陆成功，在当前这个用户的session data 保存一个键值对：isLogin=true
			// 1. 先从上下文中获取session data
			tmpSD, ok := c.Get(gin_session.SessionContextName)
			if !ok {
				panic("session middleware")
			}
			sd := tmpSD.(gin_session.SessionData)
			// 2. 给session data设置isLogin = true
			sd.Set("isLogin", true)
			//调用Save，存储到数据库
			sd.Save()
			//跳转到index界面
			c.Redirect(http.StatusMovedPermanently, toPath)
		} else { //验证失败，重新登陆
			//返回错误和重新登陆界面
			c.HTML(http.StatusOK, "login.html", gin.H{
				"err": "用户名或密码错误",
			})
			return
		}
	} else { //get

		c.HTML(http.StatusOK, "login.html", nil)
	}
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", nil)
}

func vipHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "vip.html", nil)
}
