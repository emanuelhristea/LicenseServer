package controllers

import (
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/emanuelhristea/lime/license"
	"github.com/emanuelhristea/lime/server/middleware"
	"github.com/emanuelhristea/lime/server/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// MainHandler is a ...
func MainHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(middleware.IdentityKey)

	if user == nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Authorization",
		})
	} else {
		action := c.Param("action")
		switch action {
		case "/new":
			c.HTML(http.StatusOK, "new_customer.html", gin.H{
				"title": "Create new customer",
				"Roles": license.GetRoleValues(),
			})

		case "/update":
			customer, err := customerFromParam(c)
			if err == nil {
				c.HTML(http.StatusOK, "new_customer.html", gin.H{
					"title":    "Update customer",
					"Customer": customer,
					"Roles":    license.GetRoleValues(),
				})
			}
		default:
			log.Print("serve main")
			customersList := models.CustomersList()
			c.HTML(http.StatusOK, "customers.html", gin.H{
				"title":     "Customers",
				"customers": customersList,
			})
		}
	}
}

func customerFromParam(c *gin.Context) (*models.Customer, error) {
	cID := c.Param("id")
	customerID, err := strconv.ParseUint(cID, 10, 64)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/")
	}
	customer := &models.Customer{}
	customer, err = customer.FindCustomerByID(customerID)
	if err != nil {
		c.Redirect(http.StatusNotFound, "/admin/")
	}
	return customer, err
}

func subscriptionFromParam(c *gin.Context) (*models.Subscription, error) {
	sID := c.Param("sid")
	subscriptionID, err := strconv.ParseUint(sID, 10, 64)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/")
	}
	subscription := &models.Subscription{}
	subscription, err = subscription.FindSubscriptionByID(subscriptionID)
	if err != nil {
		c.Redirect(http.StatusNotFound, "/admin/")
	}
	return subscription, err
}

func CustomerRowHandler(c *gin.Context) {
	customer, err := customerFromParam(c)
	if err == nil {
		c.HTML(http.StatusOK, "_customer_row.html", customer)
	}
}

func CustomerSubscriptionAction(c *gin.Context) {
	customerID := c.Param("id")
	action := c.Param("action")

	switch action {
	case "/":
	case "/new":
		name := CustomerNameFromID(customerID)
		tariffsList := models.TariffsList()
		c.HTML(http.StatusOK, "new_subscription.html", gin.H{
			"title":      "Subscription and Licenses for " + name,
			"customerID": customerID,
			"Tariffs":    tariffsList,
		})
	case "/update":
		name := CustomerNameFromID(customerID)
		tariffsList := models.TariffsList()
		subscription, err := subscriptionFromParam(c)
		if err == nil {
			c.HTML(http.StatusOK, "new_subscription.html", gin.H{
				"title":        "Subscription and Licenses for " + name,
				"Subscription": subscription,
				"customerID":   customerID,
				"Tariffs":      tariffsList,
			})
		}

	default:
		c.Redirect(http.StatusFound, "/admin/customer/"+customerID+"/subscriptions/")
	}

}

func CustomerSubscriptionLicenseAction(c *gin.Context) {
	cID := c.Param("id")
	sID := c.Param("sid")
	action := c.Param("action")

	switch action {
	case "/":
	case "/new":
		name := SubscriptionNameFromID(sID)
		if name == "" {
			c.Redirect(http.StatusFound, "/admin/customer/"+cID+"/subscriptions/")
			return
		}

		c.HTML(http.StatusOK, "new_license.html", gin.H{
			"title":          "Add license for " + name,
			"subscriptionID": sID,
			"customerID":     cID,
		})

	}
}

// CustomerSubscriptionList is a ...
func CustomerSubscriptionList(c *gin.Context) {
	customerID := c.Param("id")
	name := CustomerNameFromID(customerID)
	subscriptionsList := models.SubscriptionsList(customerID, "Licenses", "Licenses.Subscription", "Customer", "Tariff")

	action := c.Param("action")

	switch action {
	case "/":
		c.HTML(http.StatusOK, "subscriptions.html", gin.H{
			"title":         "Subscription and Licenses for " + name,
			"customerID":    customerID,
			"Subscriptions": subscriptionsList,
		})
	case "/new":
	case "/error":
		message, exists := c.GetQuery("message")
		if exists {
			c.HTML(http.StatusOK, "subscriptions.html", gin.H{
				"error":         message,
				"title":         "Subscription and Licenses for " + name,
				"customerID":    customerID,
				"Subscriptions": subscriptionsList,
			})
		}
	default:
		c.Redirect(http.StatusFound, "/admin/customer/"+customerID+"/subscriptions/")
	}
}

func CustomerNameFromID(customerID string) string {
	customer := &models.Customer{}
	_id, err := strconv.ParseUint(customerID, 10, 64)
	if err != nil {
		return ""
	}

	_customer, err := customer.FindCustomerByID(_id)
	if err != nil {
		return ""
	}

	return _customer.Name
}

func SubscriptionNameFromID(subscriptionID string) string {
	subscription := &models.Subscription{}
	_id, err := strconv.ParseUint(subscriptionID, 10, 64)
	if err != nil {
		return ""
	}

	_subscription, err := subscription.FindSubscriptionByID(_id)
	if err != nil {
		return ""
	}

	return _subscription.StripeID
}

// DownloadLicense is a ...
func DownloadLicense(c *gin.Context) {
	licenseID := c.Param("id")
	license := models.License{}

	uid, err := strconv.ParseUint(licenseID, 10, 64)
	if err != nil {
		panic(err)
	}

	license.FindLicenseByID(uid)

	body := base64.StdEncoding.EncodeToString(license.License)
	reader := strings.NewReader(body)
	contentLength := int64(len(body))
	contentType := "application/octet-stream"
	extraHeaders := map[string]string{"Content-Disposition": `attachment; filename="` + license.Hash + `"`}
	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

// TariffsList is a ...
func TariffsList(c *gin.Context) {
	action := c.Param("action")

	switch action {
	case "/":
		tariffsList := models.TariffsList()
		c.HTML(http.StatusOK, "tariffs.html", gin.H{
			"title":   "Pricing",
			"Tariffs": tariffsList,
		})
	case "/new":
		c.HTML(http.StatusOK, "new_tariff.html", gin.H{
			"title": "Create new pricing",
		})
	}
}

// TariffsAction is a ..
func TariffAction(c *gin.Context) {
	id := c.Param("id")
	action := c.Param("action")
	switch action {
	case "/":
	case "/delete":
		tariffId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			panic(err)
		}
		models.DeleteTariff(tariffId)
		tariffsList := models.TariffsList()
		c.HTML(http.StatusOK, "tariffs.html", gin.H{
			"title":   "Pricing",
			"Tariffs": tariffsList,
		})
	case "/update":
		tariffId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.Redirect(http.StatusFound, "/admin/tariffs")
		}
		tariff := &models.Tariff{}
		tariff, err = tariff.FindTariffByID(tariffId)
		if err != nil {
			c.Redirect(http.StatusNotFound, "/admin/tariffs")
		}

		c.HTML(http.StatusOK, "new_tariff.html", gin.H{
			"title":  "Update pricing plan",
			"Tariff": tariff,
		})
	}
}
