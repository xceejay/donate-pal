package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/xceejay/boilerplate/models"
	"github.com/xceejay/boilerplate/services"
)

type AccountController struct {
}

const (
	userkey = "user"
)

/////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////HANDLERS////////////////////////////////////////////////

// HandleLogin is a simple middleware to login
func (accountcontroller AccountController) HandleLogin(c *gin.Context) {

	// session := sessions.Default(c)
	// user := session.Get(userkey)
	// usernameSessionstring := fmt.Sprintf("%s", user)
	// if user == nil {
	// 	c.HTML(http.StatusOK, "login.html", nil)
	// } else {
	// 	c.Redirect(http.StatusPermanentRedirect, "/account/"+usernameSessionstring)
	// }

	if !isLoggedIn(c) {
		c.HTML(http.StatusOK, "login.html", nil)
	} else {

		c.Redirect(http.StatusPermanentRedirect, "/account/admin")
	}
}

func (accountcontroller AccountController) HandleLogout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}

	// session.Delete(userkey)

	session.Set("user", "") // this will mark the session as "written" and hopefully remove the username
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1}) // this sets the cookie with a MaxAge of 0

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	c.Redirect(http.StatusTemporaryRedirect, "/login")
}

// func (accountController AccountController) HandleAccountPage(c *gin.Context) {
// 	session := sessions.Default(c)
// 	username := session.Get(userkey)

// 	usernameSessionstring := fmt.Sprintf("%s", username)

// 	myUrl, err := url.Parse(c.Request.RequestURI)

// 	if err != nil {
// 		log.Errorf("ERROR: %s\nHandleAccountPage: Could Not Parse Request URI\nRedirecting to 404\n", err)
// 		c.Redirect(http.StatusNotFound, "/404")
// 	}
// 	urlUsername := path.Base(myUrl.Path)
// 	fmt.Printf("Request URI: %s\nBase: %s\n", c.Request.RequestURI, urlUsername)

// 	if usernameSessionstring != urlUsername {

// 		c.Redirect(http.StatusPermanentRedirect, "/login")
// 		return
// 	}

// 	if urlUsername == usernameSessionstring {

// 		location := fmt.Sprintf("/%s/%s", "account", username)

// 		fmt.Println("LOCATION:", location)
// 		accountController.ServeAccountPage(c)
// 	}
// }

func (accountController AccountController) HandleAdminAccountPage(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get(userkey)

	usernameSessionstring := fmt.Sprintf("%s", username)

	myUrl, err := url.Parse(c.Request.RequestURI)

	if err != nil {
		log.Errorf("ERROR: %s\nHandleAccountPage: Could Not Parse Request URI\nRedirecting to 404\n", err)
		c.Redirect(http.StatusNotFound, "/404")
	}
	urlUsername := path.Base(myUrl.Path)
	fmt.Printf("Request URI: %s\nBase: %s\n", c.Request.RequestURI, urlUsername)

	if usernameSessionstring != urlUsername {

		c.Redirect(http.StatusPermanentRedirect, "/login")
		return
	}

	if urlUsername == usernameSessionstring {

		location := fmt.Sprintf("/%s/%s", "account", username)

		fmt.Println("LOCATION:", location)
		accountController.ServeAdminAccountPage(c)
	}
}

func (accountController AccountController) HandlePayment(c *gin.Context) {
	c.HTML(http.StatusOK, "payment.html", nil)

}

func (accountController AccountController) HandleAdminDashboardContent(c *gin.Context) {

	if isLoggedIn(c) {

		dashboardContent := c.Param("dashboard_content")
		switch dashboardContent {
		case "transactions":
			accountController.ServeAdminAccountTransactionPage(c)
			return
		case "balance":
			accountController.ServeAdminAccountBalancePage(c)
			return
		case "receipt":
			accountController.ServeAdminAccountReceiptPage(c)
			return
		case "settings":
			accountController.ServeAdminAccountSettingsPage(c)
			return
		case "overview":
			accountController.ServeAdminAccountOverviewPage(c)

			return
		case "":
			accountController.ServeAdminAccountOverviewPage(c)

			return

		}
	}
	accountController.HandleLogin(c)

}

/////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////MISC////////////////////////////////////////////////

// // login is a handler that parses a form and checks for specific data
func (accountcontroller AccountController) PerformLogin(c *gin.Context) {

	user := new(models.User)

	session := sessions.Default(c)
	user.Username = c.PostForm("username")
	user.Password = c.PostForm("password")

	// Validate form input
	if strings.Trim(user.Username, " ") == "" || strings.Trim(user.Password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Check for username and password match from a database
	if !user.AuthencateUser(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// session.Options(sessions.Options{MaxAge:	})
	// Save the username in the session
	session.Set(userkey, user.Username) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusPermanentRedirect, "/account/"+user.Username)

}

func (accountController AccountController) ServeAdminAccountPage(c *gin.Context) {
	accountController.ServeAdminAccountOverviewPage(c)
}

func (accountController AccountController) ServeAdminAccountOverviewPage(c *gin.Context) {
	content := accountController.getAdminDashboardContent("overview", c)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
}

func (accountController AccountController) ServeAdminAccountTransactionPage(c *gin.Context) {
	content := accountController.getAdminDashboardContent("transactions", c)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
}

func (accountController AccountController) ServeAdminAccountBalancePage(c *gin.Context) {
	content := accountController.getAdminDashboardContent("balance", c)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
}

func (accountController AccountController) ServeAdminAccountReceiptPage(c *gin.Context) {
	content := accountController.getAdminDashboardContent("receipt", c)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
}

func (accountController AccountController) ServeAdminAccountSettingsPage(c *gin.Context) {
	content := accountController.getAdminDashboardContent("settings", c)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
}

func (accountController AccountController) getAdminDashboardContent(page string, c *gin.Context) string {
	paths := []string{
		"views/html/account/header.html",
		"views/html/account/footer.html",
		"views/html/account/account_top.html",
		"views/html/account/account_bottom.html",
		"views/html/account/overview.html",
		"views/html/account/transactions.html",
		"views/html/account/balance.html",
		"views/html/account/receipt.html",
		"views/html/account/settings.html",
	}
	vars := make(map[string]interface{})

	var user models.User

	user, err := user.GetAllUserData(c.Query("username"))
	if err != nil {
		panic(err)
	}

	vars["username"] = user.Username
	vars["email"] = user.Email.String
	vars["firstname"] = user.Firstname.String
	vars["lastname"] = user.Lastname.String
	vars["address"] = user.Address.String
	vars["city"] = user.City.String
	vars["country"] = user.Country.String

	templateEngine := new(services.TemplateEngine)

	header := templateEngine.ProcessFile(paths[0], vars)
	footer := templateEngine.ProcessFile(paths[1], vars)
	accounthtmlPageBottom := templateEngine.ProcessFile(paths[3], vars)
	overviewContent := templateEngine.ProcessFile(paths[4], vars)

	switch page {
	case "transactions":
		vars["transactions_active"] = "active"
		vars["fundraiser_name"] = template.HTML("<b>World Health Organization Cancer Project</b>")

		accounthtmlPageTop := templateEngine.ProcessFile(paths[2], vars)
		transactionContent := templateEngine.ProcessFile(paths[5], vars)

		return header + accounthtmlPageTop + transactionContent + accounthtmlPageBottom + footer

	case "balance":

		vars["balance_active"] = "active"
		vars["fundraiser_name"] = template.HTML("<b>World Health Organization Cancer Project</b>")

		accounthtmlPageTop := templateEngine.ProcessFile(paths[2], vars)
		balanceContent := templateEngine.ProcessFile(paths[6], vars)

		return header + accounthtmlPageTop + balanceContent + accounthtmlPageBottom + footer

	case "receipt":
		vars["receipt_active"] = "active"

		receiptModel := new(models.Receipt)

		receipts, err := receiptModel.GetAllReceipts()

		if err != nil {
			fmt.Printf("ERROR GETTING RECIEPT TABLE: %v", err)
		}
		var receiptsTable string

		for index, receipt := range receipts {

			fmt.Println("index:", index)
			receiptsTable += "<tr>"

			receiptsTable += "<td>" + fmt.Sprintf("%v", receipt.Receiptid) + "</td>"

			receiptsTable += "<td>" + receipt.Email.String + "</td>"

			receiptsTable += "<td>" + receipt.Firstname.String + "</td>"
			receiptsTable += "<td>" + receipt.Lastname.String + "</td>"
			receiptsTable += "<td>" + fmt.Sprintf("%f", receipt.Amount) + "</td>"
			receiptsTable += "<td>" + receipt.PaymentMethod + "</td>"
			receiptsTable += "<td>" + receipt.Phone.String + "</td>"
			receiptsTable += "<td>" + receipt.Address.String + "</td>"
			receiptsTable += "<td>" + receipt.DateDonated.Format("2020-12-01") + "</td>"
			receiptsTable += "</tr>"

			// fmt.Println(receiptsTable)

		}
		vars["view_receipts_table"] = template.HTML(receiptsTable)
		vars["fundraiser_name"] = template.HTML("<b>World Health Organization Cancer Project</b>")

		accounthtmlPageTop := templateEngine.ProcessFile(paths[2], vars)
		receiptContent := templateEngine.ProcessFile(paths[7], vars)

		return header + accounthtmlPageTop + receiptContent + accounthtmlPageBottom + footer

	case "settings":
		vars["settings_active"] = "active"
		vars["fundraiser_name"] = template.HTML("<b>World Health Organization Cancer Project</b>")

		accounthtmlPageTop := templateEngine.ProcessFile(paths[2], vars)
		settingsContent := templateEngine.ProcessFile(paths[8], vars)

		return header + accounthtmlPageTop + settingsContent + accounthtmlPageBottom + footer

	default:
		vars["overview_active"] = "active"
		vars["fundraiser_name"] = template.HTML("<b>World Health Organization Cancer Project</b>")
		accounthtmlPageTop := templateEngine.ProcessFile(paths[2], vars)

		return header + accounthtmlPageTop + overviewContent + accounthtmlPageBottom + footer

	}

}

// func (accountController AccountController) ServeAccountPage(c *gin.Context) {

// 	paths := []string{
// 		"views/html/account/header.html",
// 		"views/html/account/footer.html",
// 		"views/html/account/account_top.html",
// 		"views/html/account/account_bottom.html",
// 		"views/html/account/overview.html",
// 		"views/htmlsettings/account/transactions.html",
// 		"views/html/account/balance.html",
// 		"views/html/account/receipt.html",
// 		"views/html/account/settings.html",
// 	}
// 	vars := make(map[string]interface{})

// 	user := accountController.GetAllUserData(c.Query("username"))
// 	vars["name"] = user.name

// 	templateEngine := new(services.TemplateEngine)

// 	header := templateEngine.ProcessFile(paths[0], vars)
// 	footer := templateEngine.ProcessFile(paths[1], vars)
// 	accounthtmlPageTop := templateEngine.ProcessFile(paths[2], vars)
// 	accounthtmlPageBottom := templateEngine.ProcessFile(paths[3], vars)
// 	overviewContent := templateEngine.ProcessFile(paths[4], vars)
// 	transactionContent := templateEngine.ProcessFile(paths[5], vars)
// 	balanceContent := templateEngine.ProcessFile(paths[6], vars)
// 	receiptContent := templateEngine.ProcessFile(paths[7], vars)
// 	settingsContent := templateEngine.ProcessFile(paths[8], vars)

// 	page := "home"
// 	switch page {
// 	case "transactions":
// 		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(header+accounthtmlPageTop+transactionContent+accounthtmlPageBottom+footer))
// 		return
// 	case "balance":
// 		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(header+accounthtmlPageTop+balanceContent+accounthtmlPageBottom+footer))

// 		return
// 	case "receipt":
// 		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(header+accounthtmlPageTop+receiptContent+accounthtmlPageBottom+footer))

// 		return
// 	case "settings":
// 		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(header+accounthtmlPageTop+settingsContent+accounthtmlPageBottom+footer))

// 		return
// 	default:
// 		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(header+accounthtmlPageTop+overviewContent+accounthtmlPageBottom+footer))
// 		return

// 	}

// }

// func status(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
// }

func isLoggedIn(c *gin.Context) bool {

	session := sessions.Default(c)
	user := session.Get(userkey)

	if user == nil {

		return false
	}
	usernameSessionstring := fmt.Sprintf("%s", user)

	return usernameSessionstring == "admin"

}

// func me(c *gin.Context) {
// 	session := sessions.Default(c)
// 	user := session.Get(userkey)
// 	c.JSON(http.StatusOK, gin.H{"user": user})
// }

// func (accountController AccountController) ServeTemplate(c *gin.Context) {

// 	c.HTML(http.StatusOK, "template.html", nil)

// }
