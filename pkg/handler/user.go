package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserID struct {
	UserId int `json:"user_id" binding:"required"`
}

type Input struct {
	UserId int `json:"user_id" binding:"required"`
	Amount int `json:"amount" binding:"required"`
}

type Transfer struct {
	UserId int `json:"user_id" binding:"required"`
	ToId   int `json:"to_id" binding:"required"`
	Amount int `json:"amount" binding:"required"`
}

func (h *Handler) Balance(c *gin.Context) {
	var input UserID

	err := c.BindJSON(&input)
	if err != nil {
		logrus.Error("balance: can't get userId: " + err.Error())
		newErrorResponse(http.StatusBadRequest, c, "something went wrong")
		return
	}

	currency := c.Query("currency")

	ans, err := h.services.User.Balance(input.UserId)
	if err != nil {
		logrus.Error("Balance: can't get balance: " + err.Error())
		newErrorResponse(http.StatusBadRequest, c, "user has no balance")
		return
	}

	if currency != "" {
		ans, err = h.services.User.ConvertBalance(ans, currency)
	}

	c.JSON(http.StatusOK, ans)
}

func (h *Handler) TopUp(c *gin.Context) {
	var input Input

	err := c.BindJSON(&input)
	if err != nil {
		logrus.Error("topUp: can't get data: " + err.Error())
		newErrorResponse(http.StatusBadRequest, c, "something went wrong")
		return
	}

	ans, err := h.services.User.TopUp(input.UserId, input.Amount)
	if err != nil {
		logrus.Error("topUp: can't top-up balance: " + err.Error())
		newErrorResponse(http.StatusInternalServerError, c, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, ans)
}

func (h *Handler) Debit(c *gin.Context) {
	var input Input

	err := c.BindJSON(&input)
	if err != nil {
		logrus.Error("debit: can't get data: " + err.Error())
		newErrorResponse(http.StatusBadRequest, c, "something went wrong")
		return
	}

	ans, err := h.services.User.Debit(input.UserId, input.Amount)
	if err != nil {
		logrus.Error("debit: can't debit balance: " + err.Error())

		switch err.Error() {
		case "insufficient funds":
			newErrorResponse(http.StatusBadRequest, c, "insufficient funds")
		default:
			newErrorResponse(http.StatusInternalServerError, c, "something went wrong")
		}

		return
	}

	c.JSON(http.StatusOK, ans)
}

func (h *Handler) Transfer(c *gin.Context) {
	var input Transfer

	err := c.BindJSON(&input)
	if err != nil {
		logrus.Error("transfer: can't get data: " + err.Error())
		newErrorResponse(http.StatusBadRequest, c, "something went wrong")
		return
	}

	ans, err := h.services.User.Transfer(input.UserId, input.ToId, input.Amount)
	if err != nil {
		logrus.Error("transfer: can't transfer balance: " + err.Error())

		switch err.Error() {
		case "insufficient funds":
			newErrorResponse(http.StatusBadRequest, c, "insufficient funds")
		case "the recipient has no balance":
			newErrorResponse(http.StatusBadRequest, c, "the recipient has no balance")
		default:
			newErrorResponse(http.StatusInternalServerError, c, "something went wrong")
		}

		return
	}

	c.JSON(http.StatusOK, ans)
}
