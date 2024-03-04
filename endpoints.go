package main

import (
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", nil)
}

func AddLeague(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, "name is required")
		return
	}
	id := NewID()
	secret := NewID()
	Leagues = append(Leagues, &League{
		Name:    name,
		ID:      id,
		Secret:  secret,
		History: []*HistoryEntry{},
	})
	c.Redirect(http.StatusFound, "/admin/"+secret)
}

func GetLeagueHome(c *gin.Context) {
	id := c.Param("id")
	league := LookupLeague(id)
	if league == nil {
		c.JSON(http.StatusNotFound, "league not found")
		return
	}
	leagueHomeData := struct {
		Name string
		ID   string
	}{league.Name, league.ID}
	c.HTML(http.StatusOK, "l_home.tmpl", leagueHomeData)
}

func GetAdminHome(c *gin.Context) {
	secret := c.Param("secret")
	league := LookupLeagueBySecret(secret)
	if league == nil {
		c.JSON(http.StatusNotFound, "league not found")
		return
	}
	leagueAdminData := struct {
		Name   string
		ID     string
		Secret string
	}{league.Name, league.ID, league.Secret}
	c.HTML(http.StatusOK, "l_admin.tmpl", leagueAdminData)
}

func GetLeagueHistory(c *gin.Context) {
	id := c.Param("id")
	league := LookupLeague(id)
	if league == nil {
		c.JSON(http.StatusNotFound, "league not found")
		return
	}
	histReversed := slices.Clone(league.History)
	slices.Reverse(histReversed)
	c.HTML(http.StatusOK, "api_hist.tmpl", histReversed)
}

func GetLeagueRatings(c *gin.Context) {
	id := c.Param("id")
	league := LookupLeague(id)
	if league == nil {
		c.JSON(http.StatusNotFound, "league not found")
		return
	}
	ratings := CalculateRatings(league.History)
	c.HTML(http.StatusOK, "api_ratings.tmpl", ratings)
}

func AddMatch(c *gin.Context) {
	secret := c.Param("secret")
	league := LookupLeagueBySecret(secret)
	if league == nil {
		c.JSON(http.StatusNotFound, "league not found")
		return
	}
	addedDat := struct {
		Secret  string
		Err     string
		AddedOk bool
	}{league.Secret, "", true}
	winner := c.PostForm("winner")
	loser := c.PostForm("loser")
	wasDraw := c.PostForm("draw") == "on"
	if winner == "" || loser == "" {
		addedDat.Err = "must specify data for all fields"
		addedDat.AddedOk = false
		c.HTML(http.StatusBadRequest, "api_on_match_added.tmpl", addedDat)
		return
	}
	league.History = append(league.History, &HistoryEntry{
		Winner:     strings.ToUpper(strings.ReplaceAll(strings.Trim(winner, " \n\r\t"), " ", "")),
		Loser:      strings.ToUpper(strings.ReplaceAll(strings.Trim(loser, " \n\r\t"), " ", "")),
		Draw:       wasDraw,
		DateString: time.Now().Format("02/01/2006 15:04:05"),
	})
	c.HTML(http.StatusOK, "api_on_match_added.tmpl", addedDat)
}
