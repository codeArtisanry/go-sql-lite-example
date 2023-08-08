package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// LoadENV
func ConnectENV() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env file not loaded properly")
	}
}

func main() {
	// ConnectENV()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
	}))

	db := Database("main.db")

	sts := `
  DROP TABLE IF EXISTS cars;
CREATE TABLE cars(id INTEGER PRIMARY KEY, name TEXT, price INT);
INSERT INTO cars(name, price) VALUES('Audi',52642);
INSERT INTO cars(name, price) VALUES('Mercedes',57127);
INSERT INTO cars(name, price) VALUES('Skoda',9000);
INSERT INTO cars(name, price) VALUES('Volvo',29000);
INSERT INTO cars(name, price) VALUES('Bentley',350000);
INSERT INTO cars(name, price) VALUES('Citroen',21000);
INSERT INTO cars(name, price) VALUES('Hummer',41400);
INSERT INTO cars(name, price) VALUES('Volkswagen',21600);
`
	_, err := db.Exec(sts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("table cars created")
	app.Use(logger.New())

	// get one row data from the cars table by id
	app.Get("/cars", func(c *fiber.Ctx) error {
		// id := c.Params("id")
		query := fmt.Sprintf("SELECT * FROM cars WHERE price = %d", 9000)
		rows, err := db.Query(query)
		if err != nil {

			return c.SendString("Error")
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var name string
			var price int
			err = rows.Scan(&id, &name, &price)
			if err != nil {
				return c.SendString("Error")
			}
			fmt.Println(id, name, price)
			jsonRespose := map[string]interface{}{
				"id":    id,
				"name":  name,
				"price": price,
			}
			return c.JSON(jsonRespose)
		}

		return nil
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	// app.Group("/", middleware.Validate())
	// app.Use(middleware.Validate())

	defer db.Close()
	log.Fatal(app.Listen(":3001"))
}

// Database Connection
func Database(database string) *sql.DB {
	// Database connection
	db, err := sql.Open("sqlite3", database)
	// For memory database
	// db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal("Database Not Connected Due To: ", err)
	}

	return db
}

// Migration
func Migration(database string) {
	db := Database(database)
	fmt.Println("Waiting For Migrations...")
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	// Apply Migration
	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		log.Fatal("Migration Not Apply Due To: ", err)
	}
	fmt.Printf("Applied %d migrations!\n", n)
}
