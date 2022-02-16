package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/emanuelhristea/lime/server/models"
	"github.com/gin-gonic/gin"
)

// GetSubscriptionList is a ...
func GetSubscriptionList(c *gin.Context) {
	customerId := c.Param("customerId")
	preload, exists := c.GetQuery("load")
	subscriptionList := &[]models.Subscription{}
	if exists {
		subscriptionList = models.SubscriptionsList(customerId, strings.Split(preload, ",")...)
	} else {
		subscriptionList = models.SubscriptionsList(customerId)
	}
	respondJSON(c, http.StatusOK, subscriptionList)
}

// GetSubscription is a ...
func GetSubscription(c *gin.Context) {
	id := c.Param("id")
	subscriptionId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	modelSubscription := models.Subscription{}
	_subscription, err := modelSubscription.FindSubscriptionByID(subscriptionId)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(c, http.StatusOK, _subscription)
}

func extractSubscriptionFromForm(c *gin.Context) (*models.Subscription, bool) {
	stripe := c.PostForm("stripe_id")
	if stripe == "" {
		respondJSON(c, http.StatusBadRequest, "Subscription identifier is invalid")
		return nil, true
	}

	tariffId, err := strconv.ParseUint(c.PostForm("tariff_id"), 10, 64)
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Invalid plan selected")
		return nil, true
	}

	status := false
	if c.PostForm("status") != "" {
		status = true
	}

	id := c.Param("id")
	customerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return nil, true
	}

	modelSubscription := &models.Subscription{
		StripeID:   stripe,
		TariffID:   tariffId,
		CustomerID: customerId,
		Status:     status,
	}
	return modelSubscription, false
}

func CreateSubscription(c *gin.Context) {
	modelSubscription, shouldReturn := extractSubscriptionFromForm(c)
	if shouldReturn {
		return
	}

	_subscription, err := modelSubscription.SaveSubscription()
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Cannot save subscription. Duplicate payment id or plan")
		return
	}

	respondJSON(c, http.StatusOK, _subscription)
}

func UpdateSubscription(c *gin.Context) {
	sId := c.Param("sid")
	subscriptionId, err := strconv.ParseUint(sId, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	modelSubscription, shouldReturn := extractSubscriptionFromForm(c)
	if shouldReturn {
		return
	}

	_subscription, err := modelSubscription.UpdateSubscription(subscriptionId)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, _subscription)
}

func DeleteSubscription(c *gin.Context) {
	sId := c.Param("sid")
	subscriptionId, err := strconv.ParseUint(sId, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	rows, err := models.DeleteSubscription(subscriptionId)
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Cannot delete customer that has active subscriptions")
		return
	}

	respondJSON(c, http.StatusOK, fmt.Sprintf("%d", rows))
}