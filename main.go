package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
)

type Film struct {
	Title    string
	Director string
}

type omdbResponse struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	Poster string `json:"Poster"`
	Plot   string `json:"Plot"`
}

type FilmCard struct {
	Title  string
	Year   string
	Poster string
	Plot   string
}

type Comment struct {
	Username string
	Comment  string
}

var comments []Comment

func main() {
	engine := html.New("./", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(cors.New())

	app.Static("/assets/css/bootstrap.min.css", "assets/css/bootstrap.min.css")
	app.Post("/add-film/", addFilm)
	app.Post("/add-comment/", addComment)
	app.Get("/get-film-info/", getFilmInfo)
	app.Get("/ws/comments/", websocket.New(handleConnections))
	app.Get("/", rootRoute)

	log.Fatal(app.Listen(":8000"))
	fmt.Println("server is running on port 8000")
}

func rootRoute(c *fiber.Ctx) error {
	films := []Film{
		{Title: "The Godfather", Director: "Francis Ford Coppola"},
		{Title: "Blade Runner", Director: "Ridley Scott"},
		{Title: "The Thing", Director: "John Carpenter"},
	}
	return c.Render("index", fiber.Map{
		"Films":    films,
		"Comments": comments,
	})
}

func addFilm(c *fiber.Ctx) error {
	time.Sleep(1 * time.Second)
	Title := c.FormValue("title")
	Director := c.FormValue("director")
	return c.Render("film-list-element", Film{Title: Title, Director: Director})
}

func getFilmInfo(c *fiber.Ctx) error {
	title := c.Query("title")
	if title == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing title parameter")
	}

	film, err := fetchFilm(title)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Render("film-info-element", FilmCard{Poster: film.Poster, Title: film.Title, Year: film.Year, Plot: film.Plot})
}

func fetchFilm(title string) (*omdbResponse, error) {
	/*resp, err := http.Get(fmt.Sprintf("https://www.omdbapi.com/?t=%s", title))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var film omdbResponse
	err = json.Unmarshal(body, &film)
	if err != nil {
		return nil, err
	}
	return &film, nil*/

	// Seed the random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random integer between 0 and 99
	num := rng.Intn(100) + 1900

	// Convert the integer to a string
	numStr := strconv.Itoa(num)

	return &omdbResponse{
		Title:  "placeholder - " + numStr,
		Year:   numStr,
		Poster: "https://i.ytimg.com/vi/lgwT6tDniko/maxresdefault.jpg",
		Plot:   "place holder wow what a fantastic movie",
	}, nil
}

func addComment(c *fiber.Ctx) error {
	username := c.FormValue("username")
	comment := c.FormValue("comment")
	comments = append(comments, Comment{Username: username, Comment: comment})
	broadcast(Comment{Username: username, Comment: comment})
	return c.Render("comments-element", Comment{Username: username, Comment: comment})
}

var conns = make(map[*websocket.Conn]bool)

func handleConnections(c *websocket.Conn) {
	conns[c] = true
	defer delete(conns, c)
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		if messageType == websocket.TextMessage {
			var comment Comment
			if err := json.Unmarshal(message, &comment); err == nil {
				comments = append(comments, comment)
				broadcast(comment)
			}
		}
	}
}

func broadcast(comment Comment) {
	message, _ := json.Marshal(comment)
	for conn := range conns {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
