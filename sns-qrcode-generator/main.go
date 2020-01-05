package main

import (
	"html/template"
	"image/png"
	"log"
	"net/http"
	"os"
	"snsmod/util"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
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

	r.GET("/qrcode/:snstype", func(ctx *gin.Context) {
		snstype := ctx.Param("snstype")
		url := ""

		m := util.MySnsData{}
		m, err := util.GetUserItem()
		if err != nil {
			log.Print("An error has occurred: %s\n", err)
			return
		}

		switch snstype {
		case "facebook":
			url = "https://www.facebook.com/" + m.Facebook
		case "twitter":
			url = "twitter://user?screen_name=" + m.Twitter
		case "instagram":
			url = "https://www.instagram.com/" + m.Instagram
		case "line":
			url = "https://line.me/ti/p/" + m.Line
		default:
			log.Print("An error has occurred: %s\n", err)
			return
		}

		qrCode, err := qr.Encode(url, qr.L, qr.Auto)
		if err != nil {
			log.Print("An error has occurred: %s\n", err)
			return
		}

		qrCode, err = barcode.Scale(qrCode, 512, 512)
		if err != nil {
			log.Print("An error has occurred: %s\n", err)
			return
		}

		file, err := os.Create("./assets/images/qrcode.png")
		if err != nil {
			log.Print("An error has occurred: %s\n", err)
			return
		}
		defer file.Close()

		png.Encode(file, qrCode)

		html := template.Must(template.ParseFiles("templates/navbar.html", "templates/qrcode.html"))
		r.SetHTMLTemplate(html)
		ctx.HTML(http.StatusOK, "navbar.html", gin.H{})
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
