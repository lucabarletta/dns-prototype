package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type dnsRecord struct {
	Created time.Time `json:"created"`
	Domain  string    `json:"domain" binding:"required"`
	SubName string    `json:"subName" binding:"required"`
	Name    string    `json:"name" binding:"required"`
	Type    string    `json:"type" binding:"required"`
	Record  string    `json:"record" binding:"required"`
	Ttl     int       `json:"ttl" binding:"required"`
}

var rc = map[string]dnsRecord{
	// just for testing & demonstration
	"domain": {time.Now(), "domain", "subName", "name", "type", "192.168.1.1", 3600},
}

func getPing(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, "pong")
}

func getRecords(context *gin.Context) {
	domainName := context.Param("domainName")
	record, err := getRecordsByDomainName(domainName)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	context.IndentedJSON(http.StatusOK, record)
}

func getRecordsByDomainName(domainName string) (*dnsRecord, error) {
	if rec, found := rc[domainName]; found {
		return &rec, nil
	} else {
		return nil, errors.New("not found")
	}
}

func addRecord(context *gin.Context) {
	var dto dnsRecord
	domainName := context.Param("domainName")

	if err := context.BindJSON(&dto); err != nil {
		println(err.Error())
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given type could not be parsed"})
		return
	}

	dto.Created = time.Now()
	rc[domainName] = dto
	context.IndentedJSON(http.StatusCreated, dto)
}

func main() {
	router := gin.Default()
	router.GET("/ping", getPing)
	router.GET("/:domainName", getRecords)
	router.POST("/:domainName", addRecord)
	router.Run(":9090")
}
