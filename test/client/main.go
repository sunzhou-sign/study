package main

import (
	gin_session "client/session"
	"client/utils"
	"client/utils/mzip"
	"context"
	"database/sql"
	"fmt"
	"github.com/casbin/casbin"
	xormadapter "github.com/casbin/xorm-adapter"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"image/color"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql"
)

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0
	blackIndex = 1
)

var temp = struct {
	count int
	mu    sync.RWMutex
}{}

func AuthMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		if cookie, err := context.Cookie("cookie"); err == nil {
			if cookie == "success" {
				context.Next()
				return
			}
		}

		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "err",
		})
		context.Abort()
		return
	}
}

//自定义一个字符串
var jwtkey = []byte("www.topgoer.com")
var str string

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func main01() {
	root, _ := os.Getwd()
	fmt.Println(root)

	a := xormadapter.NewAdapter("mysql", "root:12345678@tcp(127.0.0.1:3306)/test?charset=utf8", true)
	e := casbin.NewEnforcer("./client/rbac_models.conf", a)

	err := e.LoadPolicy()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	v1 := r.Group("policy")
	{
		v1.GET("/add", func(context *gin.Context) {
			fmt.Println("增加权限！")

			sub := context.DefaultQuery("user", "admin")
			uri := context.DefaultQuery("uri", "ok")
			act := context.DefaultQuery("method", "GET")

			if !e.AddPolicy(sub, uri, act) {
				context.String(http.StatusOK, "权限已经存在！")
			} else {
				context.String(http.StatusOK, "增加成功！")
			}
		})

		v1.GET("/delete", func(context *gin.Context) {
			fmt.Println("删除权限！")

			sub := context.DefaultQuery("user", "admin")
			uri := context.DefaultQuery("uri", "ok")
			act := context.DefaultQuery("method", "GET")

			fmt.Println(sub, uri, act)

			if e.RemovePolicy(sub, uri, act) {
				context.String(http.StatusOK, "删除成功！")
			} else {
				context.String(http.StatusOK, "权限不存在！")
			}
		})

		v1.GET("/get", func(context *gin.Context) {
			fmt.Println("查看权限！")
			list := e.GetPolicy()
			for _, vList := range list {
				for _, v := range vList {
					fmt.Printf("value: %s, ", v)
				}
			}

			context.JSON(http.StatusOK, list)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(Authorize(e))
	{
		v2.GET("hello", func(context *gin.Context) {
			context.String(http.StatusOK, "v2 Hello!!!")
		})

		v2.GET("hi", func(context *gin.Context) {
			context.String(http.StatusOK, "v2 Hi!!!")
		})
	}

	v3 := r.Group("/test")
	{
		v3.GET("/json", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"html": "<b>Hello, world!</b>",
			})
		})

		// Serves literal characters
		v3.GET("/purejson", func(c *gin.Context) {
			c.PureJSON(200, gin.H{
				"html": "<b>Hello, world!</b>",
			})
		})

		v3.GET("/somejson", func(c *gin.Context) {
			// Will output  :   while(1);["lena","austin","foo"]
			c.SecureJSON(200, gin.H{
				"html": "<b>Hello, world!</b>",
			})
		})
	}

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}

func Authorize(e *casbin.Enforcer) gin.HandlerFunc {
	return func(context *gin.Context) {
		sub := context.DefaultQuery("user", "admin")
		uri := context.Request.URL.Path
		act := context.Request.Method

		fmt.Println(sub, uri, act)

		if e.Enforce(sub, uri, act) {
			fmt.Println("权限验证通过！")
			context.Next()
		} else {
			fmt.Println("权限验证失败！")
			context.String(http.StatusUnauthorized, "无权限！！！")
			context.Abort()
		}
	}
}

//颁发token
func setting(ctx *gin.Context) {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: 2,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "127.0.0.1",  // 签名颁发者
			Subject:   "user token", //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// fmt.Println(token)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		fmt.Println(err)
	}
	str = tokenString
	ctx.JSON(200, gin.H{"token": tokenString})
}

//解析token
func getting(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")

	fmt.Println("auth:", tokenString)
	//vcalidate token formate
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
		ctx.Abort()
		return
	}

	token, claims, err := ParseToken(tokenString)
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
		ctx.Abort()
		return
	}
	fmt.Println(111)
	fmt.Println(claims.UserId)
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtkey, nil
	})
	return token, Claims, err
}

func server() {

	http.HandleFunc("/ok", func(writer http.ResponseWriter, request *http.Request) {
		temp.mu.Lock()
		temp.count++
		temp.mu.Unlock()
		fmt.Fprintln(writer, "url: ", request.Host+request.URL.Path)
	})

	http.HandleFunc("/count", func(writer http.ResponseWriter, request *http.Request) {
		temp.mu.RLock()
		fmt.Fprintf(writer, "count %d\n", temp.count)
		temp.mu.RUnlock()
	})

	http.HandleFunc("/info", func(writer http.ResponseWriter, request *http.Request) {
		resInfo := fmt.Sprintf("%s\t%s\t%s\n", request.Method, request.URL, request.Proto)
		for k, v := range request.Header {
			resInfo += fmt.Sprintf("Header[%q] = %q\n", k, v)
		}

		resInfo += fmt.Sprintf("Host = %q\n", request.Host)
		resInfo += fmt.Sprintf("RemoteAddr = %q\n", request.RemoteAddr)

		if err := request.ParseForm(); err != nil {
			fmt.Println("err: ", err.Error())
		}

		for k, v := range request.Form {
			resInfo += fmt.Sprintf("Form[%q] = %q\n", k, v)
		}

		writer.Write([]byte(resInfo))
	})

	log.Fatal(http.ListenAndServe("localhost:8888", nil))
}

type tree struct {
	value       int
	left, right *tree
}

func main02() {
	r := gin.Default()
	r.LoadHTMLGlob("client/templates/*")

	gin_session.InitMgr("memory", "")
	r.Use(gin_session.SessionMiddleware(gin_session.MgrObj))

	r.Any("/login", loginHandler)
	r.GET("/index", indexHandler)
	r.GET("/home", homeHandler)
	r.GET("/vip", AuthMiddleware, vipHandler)

	r.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "404.html", nil)
	})

	r.Run(":8888")
}

type User struct {
	gorm.Model
	Name         string
	Age          sql.NullInt64
	Birthday     *time.Time
	Email        string  `gorm:"type:varchar(100);unique_index"`
	Role         string  `gorm:"size:255"`         // 设置字段的大小为255
	MemberNumber *string `gorm:"unique; not null"` // 设置memberNumber 字段唯一且不为空
	Num          int     `gorm:"AUTO_INCREMENT"`   // 设置 Num 字段自增
	Address      string  `gorm:"index:addr"`       // 给 Address 创建一个名字是`addr`的索引
	IgnoreMe     int     `gorm:"-"`                // 忽略这个字段
}

type Animal struct {
	gorm.Model
	Name string `gorm:"default:'galeone'"`
	Age  sql.NullInt64
}

func (a Animal) TableName() string {
	return "animals_t"
}

func main03() {
	db, err := gorm.Open("mysql", "root:12345678@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	animal := &Animal{Name: ""}

	db.AutoMigrate(Animal{})

	db.Create(animal)

	var a Animal
	err = db.First(&a).Error
	if err != nil {
		fmt.Println(err)
	}

	db.Delete(Animal{})

	fmt.Println(utils.MyJson(a))

	fmt.Println(strconv.FormatFloat(3.1413444, 'f', -1, 32))
}

func main04() {
	HttpHandler()
}

func HttpHandler() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	deal(ctx)
}

func deal(ctx context.Context) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		default:
			t, d := ctx.Deadline()
			fmt.Printf("deal time is %#v%#v\n", t.Format("2006-01-02 15:04:05"), d)
		}
	}
}

type p struct {
	name string
}

func (h *p) get() {
	if h != nil {
		fmt.Println("get")
	} else {
		fmt.Println("nil")
	}

}

func main() {
	//err := mapk.GetNewApkInfo("https://cdn.acmemob.com/202210/13/383e70fd2289db1178d77f708609484c.apk")
	//if err != nil {
	//	fmt.Println(err)
	//}

	zipPath := "/Users/sunzhou/Downloads/tt/mojiweather-V9.0609.02-20221102-release-500247.apk"
	unzipPath := "/Users/sunzhou/Downloads/tt/unzip"

	//获取当前目录下的文件或目录名(包含路径)
	err := mzip.Unzip(zipPath, unzipPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	libPath := unzipPath + "/lib"
	fileNames, err := GetPathFileOrDirName(libPath)
	if err != nil {
		fmt.Println("get name err: ", err)
	}

	fmt.Println(fileNames)

	fmt.Println(GetCpuModel(fileNames))
}

func GetPathFileOrDirName(path string) ([]string, error) {
	filepathNames, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, 4)
	for i := range filepathNames {
		fileName := strings.Split(filepathNames[i], "/")
		names = append(names, fileName[len(fileName)-1])
	}

	return names, err
}

func GetCpuModel(libFileNames []string) int {
	res := -1

	var support32, support64 bool
	for _, name := range libFileNames {
		if strings.Contains(strings.ToLower(name), "v8") {
			support64 = true
		}
		if strings.Contains(strings.ToLower(name), "v7") ||
			strings.Contains(strings.ToLower(name), "v6") ||
			strings.Contains(strings.ToLower(name), "v5") {
			support32 = true
		}
	}

	if support64 && support32 {
		res = 0
	} else if support32 && !support64 {
		res = 1
	} else if !support32 && support64 {
		res = 2
	}

	return res
}
