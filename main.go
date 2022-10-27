package main

import (
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"

	database "example/web-service-gin/database"
	domain "example/web-service-gin/domain"
	"example/web-service-gin/models"
	"example/web-service-gin/services"
)

var albums = []domain.Album{
	{Id: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{Id: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{Id: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum domain.Album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.Id == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func getUsers(c *gin.Context) {
	db := database.Connect()

	result, err := db.Query("Select Id, Sign, TotalDailyLoginCount, RegisteredDate From users")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var users []domain.User

	for result.Next() {
		var user domain.User
		err := result.Scan(&user.Id, &user.Sign, &user.TotalDailyLoginCount, &user.RegisteredDate)
		if err != nil {
			db.Close()
			panic(err.Error())
		}
		users = append(users, user)
	}

	c.IndentedJSON(http.StatusOK, users)
}

func saveUser(c *gin.Context) {
	db := database.Connect()

	stmt, err := db.Prepare("INSERT INTO users(Sign,TotalDailyLoginCount,RegisteredDate,LastLoginDate) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	var newUser domain.User

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	newUser.RegisteredDate = time.Now()
	newUser.LastLoginDate = time.Now()

	_, err = stmt.Exec(newUser.Sign, newUser.TotalDailyLoginCount, newUser.RegisteredDate, newUser.LastLoginDate)
	if err != nil {
		panic(err.Error())
	}

	c.IndentedJSON(http.StatusCreated, newUser)
}

func login(c *gin.Context) {
	var loginModel models.LoginModel

	if err := c.BindJSON(&loginModel); err != nil {
		return
	}

	if loginModel.Username == "ihsan" && loginModel.Password == "password" {
		token, err := services.GenerateJWT(loginModel.Username)

		if err != nil {
			panic(err.Error())
		}

		response := map[string]interface{}{
			"token":      token,
			"loginModel": loginModel,
		}

		c.IndentedJSON(http.StatusOK, response)
		return
	}

	c.IndentedJSON(http.StatusNotFound, "User not found")
}

func authTest(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")

	_, err := services.ValidateJWT(strings.Replace(tokenStr, "Bearer ", "", 1))

	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, "Your token is valid")
}

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		if c.Request.URL.Path != "/login" {
			token := c.GetHeader("Authorization")
			_, err := services.ValidateJWT(strings.Replace(token, "Bearer ", "", 1))

			if err != nil {
				c.AbortWithError(401, err)
				c.IndentedJSON(http.StatusUnauthorized, err.Error())
				return
			}
		}

		c.Next()

		latency := time.Since(t)
		print(latency)

		status := c.Writer.Status()
		print(status)
	}
}

func main() {
	ginEngine := gin.New()

	ginEngine.Use(GinMiddleware())

	ginEngine.GET("/albums", getAlbums)
	ginEngine.POST("/albums", postAlbums)
	ginEngine.GET("/albums/:id", getAlbumById)
	ginEngine.GET("/users", getUsers)
	ginEngine.POST("/users", saveUser)

	ginEngine.POST("/login", login)
	ginEngine.GET("/auth-test", authTest)

	ginEngine.Run("localhost:8080")
}
