package main

import (
	"html/template"
	"log"
	"net/http"
	"snsmod/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/navbar.html", "templates/index.html"))
		r.SetHTMLTemplate(html)
		ctx.HTML(http.StatusOK, "navbar.html", gin.H{})
	})

	r.GET("/settings", func(ctx *gin.Context) {
		m := util.MySnsData{}
		m, err := util.GetUserItem()
		if err != nil {
			log.Print("An error has occurred: %s\n", err)
			return
		}
		html := template.Must(template.ParseFiles("templates/navbar.html", "templates/settings.html"))

		r.SetHTMLTemplate(html)
		ctx.HTML(http.StatusOK, "navbar.html", gin.H{"mySnsData": m})
	})

	r.GET("/mkqrcode", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/navbar.html", "templates/mkqrcode.html"))

		r.SetHTMLTemplate(html)
		ctx.HTML(http.StatusOK, "navbar.html", gin.H{})
	})

	r.POST("/save", func(ctx *gin.Context) {
		m := util.MySnsData{}

		m.Facebook = strings.TrimSpace(ctx.PostForm("facebook"))
		m.Twitter = strings.TrimSpace(ctx.PostForm("twitter"))
		m.Instagram = strings.TrimSpace(ctx.PostForm("instagram"))
		m.Line = strings.TrimSpace(ctx.PostForm("line"))

		_, err := util.SaveUserItem(m)
		if err != nil {
			log.Print("An error has occurred: %s\n", err)
			return
		}
		ctx.Redirect(302, "/settings")
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
