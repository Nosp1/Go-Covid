package scraper

import (
	"crypto/tls"
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetReportLinks(url string) (string, string) {
	transportrequest := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transportrequest}
	response, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	readBytes, err := ioutil.ReadAll(response.Body)

	newestReport := GetNewReport(string(readBytes))
	splicedHTML := strings.Fields(string(readBytes))

	today := GetUrl(splicedHTML, newestReport)
	yesterday := GetUrl(splicedHTML, newestReport -1)

	return today, yesterday


}

func GetNewReport(html string) int {
	newestReport := 60

	for strings.Contains(html, "sitrep-"+strconv.Itoa(newestReport)+"-covid-19") {
		newestReport++
	}
	newestReport--

	return newestReport
}

func GetUrl(splicedHTML [] string, newestReport int) string {
	sitRepNum := 20

	for i := 0; i < len(splicedHTML); i++ {
		if strings.Contains(splicedHTML[i], "sitrep-"+strconv.Itoa(newestReport)+"-covid-19") {
			fmt.Println("Looking Report", newestReport)
			sitRepNum = i
			break
		}
	}
	url := splicedHTML[sitRepNum]
	begin := 0
	end := 0

	for i := 0; i < len(url); i++ {
		if begin == 0 && []byte(url[i:i+1]) [0] == 34 {
			begin = i + 1
		} else if []byte(url[i : i +1 ]) [0] == 34 {
			end = i
			break
		}
	}
	return "https://www.who.int" + url[begin:end]
}
