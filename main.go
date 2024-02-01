package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type App struct {
	Client       *http.Client
	DivaEndpoint string
	ApiToken     string
}

var b32AddressRegexPattern = regexp.MustCompile(`^[a-z0-9]{52}$`)
var domainRegexPattern = regexp.MustCompile(`^[a-z0-9-_]{3,64}\.i2p$`)

type dnsRecord struct {
	Sequence   int    `json:"seq" binding:"required"`
	Command    string `json:"command" binding:"required"`
	NameServer string `json:"ns" binding:"required"`
	Data       string `json:"d" binding:"required"`
}

type ApiToken struct {
	Header string `json:"header"`
	Token  string `json:"token"`
}

func (app *App) getRecords(context *gin.Context) {
	domainName := context.Param("domainName")

	if !domainRegexPattern.MatchString(domainName) {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "domain format invalid"})
		return
	}

	requestURL := fmt.Sprintf("%s/state/I2PDNS:%s", app.DivaEndpoint, domainName)
	res, err := app.Client.Get(requestURL)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not verify with diva endpoint"})
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error reading response body"})
		return
	}

	context.IndentedJSON(res.StatusCode, body)
}

func (app *App) addRecord(context *gin.Context) {
	domainName := context.Param("domainName")
	b32Address := context.Param("b32Address")

	if !domainRegexPattern.MatchString(domainName) || !b32AddressRegexPattern.MatchString(b32Address) {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input format"})
		return
	}

	app.refreshApiToken()

	payload := dnsRecord{
		Sequence:   1,
		Command:    "data",
		NameServer: fmt.Sprintf("I2PDNS:%s", domainName),
		Data:       fmt.Sprintf("%s=%s", domainName, b32Address),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not convert payload to json"})
		return
	}

	requestURL := fmt.Sprintf("%s/tx/", app.DivaEndpoint)
	req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error creating http request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("diva-token-api", app.ApiToken)

	res, err := app.Client.Do(req)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error executing http request"})
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error reading response body"})
		return
	}

	context.IndentedJSON(http.StatusCreated, body)
}

func (app *App) refreshApiToken() {
	requestURL := fmt.Sprintf("%s/testnet/token", app.DivaEndpoint)

	res, err := app.Client.Get(requestURL)
	if err != nil {
		fmt.Printf("error getting token: %s\n", err)
		return
	}
	defer res.Body.Close()

	var _apiToken ApiToken
	if err := json.NewDecoder(res.Body).Decode(&_apiToken); err != nil {
		fmt.Printf("Error decoding JSON: %s\n", err)
		return
	}

	app.ApiToken = _apiToken.Token
}

func (app *App) initializeRoutes(router *gin.Engine) {
	router.GET("/ping", app.getPing)
	router.GET("/:domainName", app.getRecords)
	router.PUT("/:domainName/:b32Address", app.addRecord)
}

func (app *App) getPing(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, "pong")
}

func main() {
	app := App{
		Client: &http.Client{},
	}

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error: No .env file found")
	}

	app.DivaEndpoint = os.Getenv("DIVA_ENDPOINT")
	if app.DivaEndpoint == "" {
		fmt.Println("Error: DIVA_ENDPOINT is not set")
		os.Exit(1)
	}

	router := gin.Default()
	app.initializeRoutes(router)

	if err := router.Run(":9090"); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
		os.Exit(1)
	}
}
