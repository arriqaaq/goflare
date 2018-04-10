package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	CLOUD_FLARE_URL = "https://cloudflare-dns.com/dns-query?ct=%s&name=%s&type=%s"
	CONTENT_TYPE    = "application/dns-json"
)

/***
	Documentation for cloudflare:

	https://developers.cloudflare.com/1.1.1.1/dns-over-https/json-format/
***/

type CloudFlareResponse struct {
	Status   int  `json:"Status"`
	TC       bool `json:"TC"`
	RD       bool `json:"RD"`
	RA       bool `json:"RA"`
	AD       bool `json:"AD"`
	CD       bool `json:"CD"`
	Question []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
	} `json:"Question"`
	Answer []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
		TTL  int    `json:"TTL"`
		Data string `json:"data"`
	} `json:"Answer"`
}

type CloudFlare struct {
	client *http.Client
}

func (c *CloudFlare) Query(name string, qtype interface{}) (string, error) {
	queryUrl := fmt.Sprintf(CLOUD_FLARE_URL, CONTENT_TYPE, name, qtype)
	response, err := c.client.Get(queryUrl)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	contents, cErr := ioutil.ReadAll(response.Body)
	if cErr != nil {
		return "", cErr
	}
	return string(contents), nil
}
func (c *CloudFlare) Resolve(name string, qtype interface{}) (string, error) {
	var responseObj CloudFlareResponse
	queryUrl := fmt.Sprintf(CLOUD_FLARE_URL, CONTENT_TYPE, name, qtype)
	response, err := c.client.Get(queryUrl)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	contents, cErr := ioutil.ReadAll(response.Body)
	if cErr != nil {
		return "", cErr
	}
	json.Unmarshal([]byte(contents), &responseObj)

	if len(responseObj.Answer) > 0 {
		return responseObj.Answer[0].Data, nil
	}

	return "", errors.New("Unable to resolve IP")

}

func main() {
	var name, qtype, action string

	httpClient := http.Client{
		Timeout: 1 * time.Second,
	}

	gc := CloudFlare{
		client: &httpClient,
	}
	flag.StringVar(&name, "name", "", "--name")
	flag.StringVar(&action, "action", "", "--action")
	flag.StringVar(&qtype, "qtype", "AAAA", "--qtype")
	flag.Parse()

	switch action {
	case "query":
		result, err := gc.Query(name, qtype)
		if err != nil {
			log.Fatalln("error resolving dns query: ", err)
		} else {
			log.Println(result)
		}
	case "resolve":
		result, err := gc.Resolve(name, qtype)
		if err != nil {
			log.Fatalln("error resolving dns query: ", err)
		} else {
			log.Println(result)
		}
	default:
		log.Println("invalid option")
	}

}
