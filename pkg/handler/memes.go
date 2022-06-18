package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (h *Handler) GetAllMemes(c *gin.Context) {
	memes, err := h.services.GetAllMemes()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, memes)
}

func (h *Handler) LoadVezdekod(c *gin.Context) {
	h.services.LoadVezdekod()
	c.JSON(http.StatusOK, map[string]string{"status": "200"})
}

type addGroupInput struct {
	Password string `json:"password"`
	Address  string `json:"address"`
}

func (h *Handler) LoadGroup(c *gin.Context) {
	var input addGroupInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Password != viper.GetString("internal-password") {
		newErrorResponse(c, http.StatusUnauthorized, "wrong password")
		return
	}
	h.services.LoadGroup(input.Address)

	c.JSON(http.StatusOK, map[string]string{"status": "200"})
}

type inputFeed struct {
	Page int `json:"page"`
}

func (h *Handler) Feed(c *gin.Context) {
	var input inputFeed

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Page < 0 {
		newErrorResponse(c, http.StatusBadRequest, "page < 0")
	}

	meme, err := h.services.Feed(input.Page)

	if err != nil {
		logrus.Fatalf("error in /feed: %s", err)
	}

	c.JSON(http.StatusOK, meme)
}

type inputLike struct {
	MemeId string `json:"meme_id"`
}

func (h *Handler) Like(c *gin.Context) {
	var input inputLike

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, _ := strconv.Atoi(input.MemeId)

	err := h.services.Like(id, c.Request.Host)

	if err != nil {
		if err.Error() == "cannot" {
			newErrorResponse(c, http.StatusConflict, err.Error())
			return
		} else {
			logrus.Fatalf("error in /like: %s", err)
		}
	}

	c.JSON(http.StatusOK, map[string]string{"status": "200"})
}

type inputPromote struct {
	MemeId   string `json:"meme_id"`
	Password string `json:"password"`
}

func (h *Handler) Promote(c *gin.Context) {

	var input inputPromote

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Password != viper.GetString("internal-password") {
		newErrorResponse(c, http.StatusUnauthorized, "wrong password")
		return
	}

	id, _ := strconv.Atoi(input.MemeId)
	err := h.services.Promote(id)

	if err != nil {
		logrus.Fatalf("error in /promote: %s", err)
	}

	c.JSON(http.StatusOK, map[string]string{"status": "200"})
}
