package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/logrusorgru/aurora"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var err error
var DB *gorm.DB

var (
	addressPort = flag.String("addr", ":"+"8080", "")
	cert        = flag.String("cert", "", "")
	key         = flag.String("key", "", "")
)

type User struct {
	gorm.Model

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func main() {

	fmt.Println(aurora.Blink(aurora.Green("server is has started running")))

	/*CONNECT TO THE DATABASE*/
	OpenDB()

	/*HNADLE REQUESTS*/
	if err := handleRoutes(); err != nil {
		log.Fatal(err.Error())
	}

}

func handleRoutes() error {

	if *addressPort == ":" {
		*addressPort = ":8080"
	}

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000",
		"http://192.168.0.106:6080",
		"http://192.168.0.106:4000"}
	config.AllowMethods = []string{"PUT", "PATCH", "POST", "GET", "OPTIONS", "HEAD"}
	config.AllowHeaders = []string{
		"Authorization",
		"X-Requested-With",
		"Content-Length",
		"Content-Type",
		"Accept",
		"Accept-Encoding",
		"Authorization",
		"X-CSRF-Token",
		"Cache-Control",
		"Origin",
	}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	// config.AllowOriginFunc = func(origin string) bool {
	// 	return origin == "https://github.com"
	// }
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))

	router.POST("/", func(c *gin.Context) {

		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err,
			})
			return
		}

		n := User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}

		createdUSer := DB.Create(&n)
		err := createdUSer.Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("errror is %v", err.Error()),
			})
			return
		}

		c.JSON(http.StatusOK, n)
	})

	router.GET("/:uuid", func(c *gin.Context) {

		var users []User

		uuid := c.Param("uuid")
		if uuid == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "User id is required",
			})
			return
		}

		DB.Find(&users, "email LIKE  ? OR first_name LIKE ? OR last_name LIKE ?", "%"+uuid+"%", "%"+uuid+"%", "%"+uuid+"%")

		c.JSON(http.StatusOK, users)
	})

	srv := &http.Server{
		Addr:         "0.0.0.0" + *addressPort,
		WriteTimeout: time.Second * 15, /*Good practice to set timeouts to avoid Slowloris attacks.*/
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, /*Pass our instance of gorilla/mux in.*/
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

	return srv.ListenAndServe()

}

func OpenDB() *gorm.DB {

	host := "localhost"
	dbport := "5432"
	user := "postgres"
	dbName := "bg"
	password := "admin"

	/*Database connection string for local posgress*/
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, dbport, user, dbName, password)

	/*opening database connection*/
	DB, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal(err)

	} else {

		fmt.Println(aurora.Blink(aurora.Green("Successfully connected to database")))

		//AUTOMIGRATE MODELS
		DB.AutoMigrate(&User{})

	}

	return DB

}
