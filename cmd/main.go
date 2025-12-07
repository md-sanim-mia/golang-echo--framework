package main

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getUser(c echo.Context) error {
	id := c.Param("id")

	res, err := http.Get("https://jsonplaceholder.typicode.com/todos/" + id)

	if err != nil {

		return c.String(http.StatusInternalServerError, "json placeholder api fatch problem")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {

		return c.String(http.StatusInternalServerError, "file to convert json data")
	}

	return c.String(http.StatusOK, string(body))
}


func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/api/v1", func(c echo.Context) error {

		return c.String(http.StatusOK, "yeha ho golang echo famwork is run now ")
	})

	e.GET("/api/v1/users/:id", getUser)

	e.Logger.Fatal(e.Start(":5000"))
}
