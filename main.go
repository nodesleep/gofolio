package main

import (
	"html/template"
	"io"

	"crypto/rand"
	"math/big"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

type Product struct {
	Key string
}

type Image struct {
	Url string
}

type KeyConfig struct {
	KeyLength  int
	GroupCount int
	Charset    string
}

func generateProductKey(config KeyConfig) string {
	key := ""
	for i := 0; i < config.GroupCount; i++ {
		if i > 0 {
			key += "-"
		}
		key += generateRandomString(config.KeyLength, config.Charset)
	}
	return key
}

func generateRandomString(length int, charset string) string {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		result[i] = charset[num.Int64()]
	}
	return string(result)
}


func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {
	config := KeyConfig{
		KeyLength:  5,
		GroupCount: 5,
		Charset:    "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	}

	
	e := echo.New()
	e.Use(middleware.Logger())
	
	htmx := Product{Key: "0"}
	imageUrl := Image{Url: "/static/featured-06-11-24-1.webp"}

	e.Renderer = newTemplate()
	
	e.Static("/static", "static")
	
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", imageUrl)
	})

	e.POST("/about", func(c echo.Context) error {
		return c.Render(200, "about", nil)
	})
	
	e.POST("/greeting", func(c echo.Context) error {
		return c.Render(200, "greeting", imageUrl)
	})

	e.POST("/htmx", func(c echo.Context) error {
		htmx.Key = generateProductKey(config)
		return c.Render(200, "htmx", htmx)
	})
	
	e.Logger.Fatal(e.Start(":42069"))
}
