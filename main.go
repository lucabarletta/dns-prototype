package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"regexp"
	"time"
)

var ipRegexPattern = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
var domainRegexPattern = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}$`)

type dnsRecord struct {
	Created time.Time `json:"created"`
	Domain  string    `json:"domain" binding:"required,domainNameValidator"`
	SubName string    `json:"subName" binding:"required"`
	Name    string    `json:"name" binding:"required"`
	Type    string    `json:"type" binding:"required"`
	Records []string  `json:"records" binding:"required,recordArrayValidator"`
	Ttl     int       `json:"ttl" binding:"required"`
}

func recordArrayValidator(input validator.FieldLevel) bool {
	recordList := input.Field().Interface().([]string)
	if len(recordList) < 1 {
		return false
	}
	for _, ip := range recordList {
		if !ipRegexPattern.MatchString(ip) {
			return false
		}
	}
	return true
}

func domainNameValidator(input validator.FieldLevel) bool {
	domainName := input.Field().Interface().(string)
	if !domainRegexPattern.MatchString(domainName) {
		return false
	}
	return true
}

var rc = map[string]dnsRecord{
	// just for testing & demonstration
	"domain": {Created: time.Now(), Domain: "domain", SubName: "subName", Name: "name", Type: "type", Records: []string{"192.168.1.1", "192.168.1.2"}, Ttl: 3600},
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

	if err := context.ShouldBindJSON(&dto); err != nil {
		println(err.Error())
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given type could not be parsed"})
		return
	}

	dto.Created = time.Now()
	rc[domainName] = dto
	context.IndentedJSON(http.StatusCreated, dto)
}

func registerValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("recordArrayValidator", recordArrayValidator)
		if err != nil {
			return
		}
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("domainNameValidator", domainNameValidator)
		if err != nil {
			return
		}
	}
}

func main() {
	router := gin.Default()
	registerValidator()
	router.GET("/ping", getPing)
	router.GET("/:domainName", getRecords)
	router.PUT("/:domainName", addRecord)
	err := router.Run(":9090")
	if err != nil {
		return
	}
}
