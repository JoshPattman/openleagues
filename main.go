package main

import (
	"database/sql"
	"flag"
	"math"
	"math/rand"
	"sort"

	"github.com/gin-gonic/gin"
)

var AppDB *sql.DB

func main() {

	var databaseFileName string
	flag.StringVar(&databaseFileName, "dbn", "", "Database file name")
	flag.Parse()

	if databaseFileName == "" {
		panic("Database file name is required")
	}

	db, err := DBCreateAndConnect(databaseFileName)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	AppDB = db

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.LoadHTMLGlob("content/*")

	r.GET("/", GetHome)

	r.GET("/league/:id", GetLeagueHome)

	r.GET("/admin/:secret", GetAdminHome)

	r.POST("/api/add-league", AddLeague)
	r.GET("/api/history/:id", GetLeagueHistory)
	r.GET("/api/ratings/:id", GetLeagueRatings)
	r.POST("/api/add-match/:secret", AddMatch)

	r.StaticFile("/style", "./content/style.css")
	r.Run(":8080")
}

type Rating struct {
	Name       string
	Rating     float64
	RatingInt  int
	Wins       int
	Losses     int
	Games      int
	WinPercent int
	Place      int
}

// Will return ratings sorted by rating
func CalculateRatings(history []DBHistory) []*Rating {
	ratings := map[string]*Rating{}
	for _, h := range history {
		if _, ok := ratings[h.Winner]; !ok {
			ratings[h.Winner] = &Rating{
				Name:   h.Winner,
				Rating: 1000,
			}
		}
		if _, ok := ratings[h.Loser]; !ok {
			ratings[h.Loser] = &Rating{
				Name:   h.Loser,
				Rating: 1000,
			}
		}
		if !h.Draw {
			points := 0.1 * ratings[h.Loser].Rating
			ratings[h.Winner].Rating += points
			ratings[h.Loser].Rating -= points
			ratings[h.Winner].Wins++
			ratings[h.Loser].Losses++
		}
		ratings[h.Winner].Games++
		ratings[h.Loser].Games++
	}
	var result []*Rating
	for _, rating := range ratings {
		rating.RatingInt = int(math.Round(rating.Rating))
		if rating.Games > 0 {
			rating.WinPercent = int(math.Round(100 * float64(rating.Wins) / float64(rating.Games)))
		}
		result = append(result, rating)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Rating > result[j].Rating
	})
	for i := range result {
		result[i].Place = i + 1
	}
	return result
}

func NewID() string {
	alpahbet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	text := ""
	for i := 0; i < 10; i++ {
		text += string(alpahbet[rand.Intn(len(alpahbet))])
	}
	return text
}
