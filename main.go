package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var divaEndpoint string

var b32AddressRegexPattern = regexp.MustCompile(`^[a-z0-9]{52}$`)
var domainRegexPattern = regexp.MustCompile(`^[a-z0-9-_]{3,64}\.i2p$`)

type dnsRecord struct {
	Sequence   int    `json:"seq" binding:"required"`
	Command    string `json:"command" binding:"required"`
	NameServer string `json:"ns" binding:"required"`
	Data       string `json:"d" binding:"required"`
}

func b32AddressValidator(input string) bool {
	return b32AddressRegexPattern.MatchString(input)
}

func domainNameValidator(input string) bool {
	return domainRegexPattern.MatchString(input)
}

func getPing(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, "pong")
}

func getRecords(context *gin.Context) {
	domainName := context.Param("domainName")

	if !domainNameValidator(domainName) {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "domain format invalid"})
		return
	}

	requestURL := fmt.Sprintf("%s/state/decision:I2PDNS:%s", divaEndpoint, domainName)
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not verify with diva endpoint"})
		return
	}
	defer res.Body.Close()

	fmt.Println("The status code we got is:", res.StatusCode)


	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading response body: %s\n", err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not verify with diva endpoint"})
		return
	}

	context.IndentedJSON(res.StatusCode, resBody)
}

func addRecord(context *gin.Context) {
	domainName := context.Param("domainName")
	b32Address := context.Param("b32Address")

	if !domainNameValidator(domainName) {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "domain format invalid"})
		return
	}
	if !b32AddressValidator(b32Address) {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "b32 address format invalid"})
		return
	}

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

	requestURL := fmt.Sprintf("%s/transaction/", divaEndpoint)
	res, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not verify with diva endpoint"})
		return
	}
	defer res.Body.Close()

//	fmt.Printf("%d\n", res.Response.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading response body: %s\n", err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not verify with diva endpoint"})
		return
	}

	context.IndentedJSON(http.StatusCreated, resBody)
}

func main() {
	godotenv.Load()
	divaEndpoint = os.Getenv(`DIVA_ENDPOINT`)

	router := gin.Default()

	router.GET("/ping", getPing)
	router.GET("/:domainName", getRecords)
	router.PUT("/:domainName/:b32Address", addRecord)

	err := router.Run(":9090")

	if err != nil {
		return
	}
}
