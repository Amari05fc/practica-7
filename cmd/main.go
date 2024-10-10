package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Amari05fc/practica-7/database"
	"github.com/Amari05fc/practica-7/servicios"
	"github.com/gin-gonic/gin"
)

type User struct {
	Id    int    `json: "id"`
	Name  string `json: "name"`
	Email string `json: "email"`
}

func main() {
	//Base de datos conectada
	//El NewDatabaseDriver es para crear un nueva BD
	db, err := database.NewDatabaseDriver()
	if err != nil {
		fmt.Println("Error al conectar a la base de datos: ", err)
		return
	}
	fmt.Println("Base de datos")
	db.AutoMigrate(&database.User{})
	fmt.Println("AutoMigrate")
	userService := servicios.NewUserService(db)

	router := gin.Default()
	users := []User{}
	indexUser := 1
	fmt.Println("Running")

	//Tomas los archivos de la carpeta template
	router.LoadHTMLGlob("templates/*")
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong fer",
		})
	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title":       "Main website",
			"total_users": len(users),
			"users":       users,
		})
	})
	//API URLs
	//Obtener usuarios
	router.GET("/api/users", func(c *gin.Context) {
		c.JSON(200, users)
	})

	//Creacion de usuarios
	router.POST("/api/users", func(c *gin.Context) {
		var user User
		if c.BindJSON(&user) == nil {
			user.Id = indexUser
			users = append(users, user)
			indexUser++
			c.JSON(200, user)
			fmt.Println(user.Id, user.Email, user.Name)
			userService.CreateUser(database.User{
				Id:    user.Id,
				Name:  user.Name,
				Email: user.Email,
			})
			fmt.Println("Insertando en la Base de Datos")
		} else {
			c.JSON(400, gin.H{
				"error": "Invalid payload",
			})
		}
	})
	
	//Eliminacion de usuarios
	router.DELETE("/api/users/:id", func(c *gin.Context) {
		fmt.Println("Delete")
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid id",
			})
			return
		}
		fmt.Println("Id a borrar: ", id)
		for i, user := range users {
			if user.Id == idParsed {
				users = append(users[:i], users[i+1:]...)
				c.JSON(200, gin.H{
					"message": "User Deleted",
				})
				return
			}
		}
		c.JSON(201, gin.H{})
	})

	//Actualizar usuarios
	router.PUT("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid id",
			})
			return
		}
		var user User
		err = c.BindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid payload",
			})
			return
		}
		fmt.Println("Id a actualizar: ", id)
		for i, u := range users {
			if u.Id == idParsed {
				users[i] = user
				users[i].Id = idParsed
				c.JSON(200, users[i])
				return
			}
		}
		c.JSON(201, gin.H{})
	})

	router.Run(":8001") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
