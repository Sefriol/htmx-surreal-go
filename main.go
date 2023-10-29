package main

import (
	"fmt"
	"html/template"
	"net/http"
	"io"
	"os"
	"path/filepath"
	"log/slog"
	
	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/surrealdb/surrealdb.go"
	
	"github.com/sefriol/htmx-surreal-go/view"
)

type User struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type Relative struct {
    ID      string `json:"id,omitempty"`
    Parent string `json:"in,omitempty"`
    Child  string `json:"out,omitempty"`
}
type SubQueryResult struct {
    Result []struct {
	Query []struct {
		ID     string `json:"id,omitempty"`
		User   User `json:"user"`
		Child  User `json:"child"`
		Parent User `json:"parent"`
	}
    }`json:"result"`
    Status string `json:"status"`
    Time   string `json:"time"`
}

type QueryResult struct {
    Result []struct {
        ID     string `json:"id,omitempty"`
	User   User `json:"user"`
        Child  User `json:"child"`
        Parent User `json:"parent"`
    }`json:"result"`
    Status string `json:"status"`
    Time   string `json:"time"`
}

type Content struct {
    Users []User
    User User
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	tmpl := template.New("index")

	var err error
	db, err := surrealdb.New("ws://localhost:8000/rpc")
	if err != nil {
		panic(err)
	}

	if tmpl, err = tmpl.Parse(view.Index); err != nil {
		fmt.Println(err)
	}

	if tmpl, err = tmpl.Parse(view.Users); err != nil {
		fmt.Println(err)
	}

	if tmpl, err = tmpl.Parse(view.User); err != nil {
		fmt.Println(err)
	}

	if tmpl, err = tmpl.Parse(view.EditUser); err != nil {
		fmt.Println(err)
	}
	if tmpl, err = tmpl.Parse(view.Relations); err != nil {
		fmt.Println(err)
	}
	if tmpl, err = tmpl.Parse(view.Relation); err != nil {
		fmt.Println(err)
	}
	if tmpl, err = tmpl.Parse(view.RelationDialog); err != nil {
		fmt.Println(err)
	}
	e.Use(middleware.Logger())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	}))


	e.Static("/css", "css");
	e.Renderer = &TemplateRenderer{
		templates: tmpl,
	}
	
	if _, err = db.Signin(map[string]interface{}{
		"user": "root",
		"pass": "root",
	}); err != nil {
		panic(err)
	}
	if _, err = db.Use("test", "test"); err != nil {
		panic(err)
	}

	
	// Create relatives graph table
	path := filepath.Join("surql","user.surql")

	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	surql := string(data)
	slog.Info(surql)

	result, err := db.Query(surql, nil)
        if err != nil {
		panic(err)
	}
	log.Warn(result)

	// Get user by ID
	userData, err := db.Select("user")
	if err != nil {
		panic(err)
	}

	// Unmarshal data
	selectedUser := new([]User)
	err = surrealdb.Unmarshal(userData, &selectedUser)
	if err != nil {
		panic(err)
	}

	items := Content{
		Users: *selectedUser,
	}
	e.GET("/", func(c echo.Context) error {
		//log.Warn(surql)
		return c.Render(http.StatusOK, "index", items)
	})

	e.GET("/user/:id/edit", func (c echo.Context) error {
		id := c.Param("id")
		// Insert user
		data, err := db.Select(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		var user User
		// Unmarshal data
		err = surrealdb.Unmarshal(data, &user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.Render(http.StatusOK, "edit-user", user)
	})

	e.GET("/user/:id/relative", func (c echo.Context) error {
		id := c.Param("id")
		c.Logger().Warn(id)
		data, err := db.Select(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		var user User
		// Unmarshal data
		err = surrealdb.Unmarshal(data, &user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}


		userData, err := db.Select("user")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		// Unmarshal data
		users := new([]User)
		err = surrealdb.Unmarshal(userData, &users)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		items := Content{
			Users: *users,
			User: user,
		}

		return c.Render(http.StatusOK, "relation-dialog", items)
	})
	
	e.GET("/user/:id/relatives", func (c echo.Context) error {
		id := c.Param("id")
		c.Logger().Warn(id)
		data, err := db.Query("SELECT id, $user.* as user, in.* AS parent, out.* as child FROM relative WHERE in=$user OR out=$user",
			map[string]string{
			"user": id,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		data_str := fmt.Sprintf("%v", data)
		log.Warn(data_str)
		result := new([]QueryResult)
		// Unmarshal data
		err = surrealdb.Unmarshal(data, &result)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.Render(http.StatusOK, "relations", (*result)[0].Result)
	})

	e.POST("/user/:id/relative", func (c echo.Context) error {
		id := c.Param("id")
		relative := c.FormValue("relative")
		relation := c.FormValue("relation")

		if id == relative {
			return echo.NewHTTPError(http.StatusBadRequest,"User cannot be their own relative!")
		}
		var parent string;
		var child string;

		if relation == "child" {
			parent = id
			child = relative
		} else {
			parent = relative
			child = id
		}

		// Insert relation
		data, err := db.Query("RELATE $parent->relative->$child RETURN (SELECT id, $user.* as user, in.* AS parent, out.* as child FROM relative WHERE in=$user OR out=$user) AS query",
			map[string]string{
			"child": child,
			"parent": parent,
			"user": id,
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		data_str := fmt.Sprintf("%v", data)
		log.Warn(data_str)
		result := make([]SubQueryResult, 1)

                err = surrealdb.Unmarshal(data, &result)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if result[0].Status == "ERR" {
			return echo.NewHTTPError(http.StatusInternalServerError, result[0].Result)
		}

		log.Warn(data_str)
		c.Logger().Warn(data_str)
		return c.Render(http.StatusOK, "relations", result[0].Result[0].Query)
	})

	e.DELETE("/user/:id/relative/:relative", func (c echo.Context) error {
		id := c.Param("id")
		relative := c.Param("relative")

		if id == relative {
			return echo.NewHTTPError(http.StatusBadRequest,"User cannot be their own relative!")
		}

		// Insert relation
		data, err := db.Query("DELETE $id->relative WHERE out=$relative; DELETE $relative->relative WHERE out=$id;",
			map[string]string{
			"id": id,
			"relative": relative,
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		data_str := fmt.Sprintf("%v", data)
		log.Warn(data_str)

		return c.NoContent(http.StatusOK)
	})

	e.GET("/user/:id", func (c echo.Context) error {
		id := c.Param("id")
		// Insert user
		data, err := db.Select(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		var user User
		// Unmarshal data
		err = surrealdb.Unmarshal(data, &user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.Render(http.StatusOK, "user", user)
	})

	e.DELETE("/user/:id", func (c echo.Context) error {
		id := c.Param("id")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
 		_, err = db.Delete(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return c.NoContent(http.StatusOK)
	})

	e.PUT("/user/:id", func (c echo.Context) error {
		id := c.Param("id")
		name := c.FormValue("name")
		surname := c.FormValue("surname")

		changes:= map[string]string{
			"name": name,
			"surname": surname,
		}
		data, err := db.Update(id, changes)
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		
		var user User
		// Unmarshal data
		err = surrealdb.Unmarshal(data, &user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.Render(http.StatusOK, "user", user)
	})

	e.POST("/user", func (c echo.Context) error {
		name := c.FormValue("name")
		surname := c.FormValue("surname")

		// Create user
		user := User{
			Name:    name,
			Surname: surname,
		}

		// Insert user
		data, err := db.Create("user", user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Unmarshal data
		createdUser := make([]User, 1)
		err = surrealdb.Unmarshal(data, &createdUser)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.Render(http.StatusOK, "user", createdUser[0])
	})

	e.Logger.Fatal(e.Start(":1323"))

}
