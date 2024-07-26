package main

import (
	"net/http"
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
	if err := DBCreateLeage(name, id, secret); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/admin/"+secret)
}

func GetLeagueHome(c *gin.Context) {
	id := c.Param("id")
	league, err := DBLookupLeageByID(id)
	if err == ErrLeagueNotFound {
		c.JSON(http.StatusNotFound, "league not found")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
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
	league, err := DBLookupLeagueBySecret(secret)
	if err == ErrLeagueNotFound {
		c.JSON(http.StatusNotFound, "league not found")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	leagueAdminData := struct {
		Name   string
		ID     string
		Secret string
	}{league.Name, league.ID, league.SecretID}
	c.HTML(http.StatusOK, "l_admin.tmpl", leagueAdminData)
}

func GetLeagueHistory(c *gin.Context) {
	id := c.Param("id")
	_, hist, err := DBLookupLeagueByIDWithHistory(id)
	if err == ErrLeagueNotFound {
		c.JSON(http.StatusNotFound, "league not found")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.HTML(http.StatusOK, "api_hist.tmpl", hist)
}

func GetLeagueRatings(c *gin.Context) {
	id := c.Param("id")
	_, hist, err := DBLookupLeagueByIDWithHistory(id)
	if err == ErrLeagueNotFound {
		c.JSON(http.StatusNotFound, "league not found")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ratings := CalculateRatings(hist)
	c.HTML(http.StatusOK, "api_ratings.tmpl", ratings)
}

func AddMatch(c *gin.Context) {
	secret := c.Param("secret")
	winner := c.PostForm("winner")
	loser := c.PostForm("loser")
	wasDraw := c.PostForm("draw") == "on"

	addedDat := struct {
		Secret  string
		Err     string
		AddedOk bool
	}{secret, "", true}

	if winner == "" || loser == "" {
		addedDat.Err = "must specify data for all fields"
		addedDat.AddedOk = false
		c.HTML(http.StatusBadRequest, "api_on_match_added.tmpl", addedDat)
		return
	}

	hist := DBHistory{
		LeagueID:   secret,
		Winner:     strings.ToUpper(strings.ReplaceAll(strings.Trim(winner, " \n\r\t"), " ", "")),
		Loser:      strings.ToUpper(strings.ReplaceAll(strings.Trim(loser, " \n\r\t"), " ", "")),
		Draw:       wasDraw,
		DateString: time.Now().Format("02/01/2006 15:04:05"),
	}
	if err := DBCreateHistory(secret, []DBHistory{hist}); err != nil {
		addedDat.Err = err.Error()
		addedDat.AddedOk = false
		c.HTML(http.StatusInternalServerError, "api_on_match_added.tmpl", addedDat)
		return
	}

	c.HTML(http.StatusOK, "api_on_match_added.tmpl", addedDat)
}
