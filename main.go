package main

import (
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Docker example https://github.com/olliefr/docker-gs-ping

func main() {
	// Load configuration
	config, err := LoadConfiguration("config.json")
	if err != nil {
		println(err)
	}

	databaseASN, err := GetGeoLite2("GeoLite2-ASN", config.Maxmind.Key)
	if err != nil {
		println(err)
	}
	config.Maxmind.ASN = databaseASN

	databaseCity, err := GetGeoLite2("GeoLite2-City", config.Maxmind.Key)
	if err != nil {
		println(err)
	}
	config.Maxmind.City = databaseCity

	// Template
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	// -----------------------------------------

	e := echo.New()
	e.Renderer = t

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	/*
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			// Be careful to use constant time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
				subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
				return true, nil
			}
			return false, nil
		}))
	*/
	e.Static("/assets", "public/assets")
	e.Static("/csv", "public/csv")
	e.Static("/json", "public/json")

	// Routes
	e.GET("/", config.getIPInfo)
	e.POST("/", config.postIPInfo)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
