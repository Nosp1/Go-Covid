package analytics

import (
	"bytes"
	"crypto/tls"
	"github.com/Nosp1/TA/Covid19/scraper"
	"github.com/ledongthuc/pdf"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Stats struct {
	TotalDeath string
	NewInfected string
	TotalInfected string
	NewDeath string
}

var fileNumber int

func GetStats() []Stats {
	newUrl, oldUrl := scraper.GetReportLinks("https://www.who.int/emergencies/diseases/novel-coronavirus-2019/situation-reports")

	today := getNorwayStatsUrl(newUrl)
	yesterday := getNorwayStatsUrl(oldUrl)

	tempStats := []Stats{today,yesterday}
	return tempStats
}

func readPlainTextPDF(pdfLocation string) (contents string, err error) {
	file, tempReader, err := pdf.Open(pdfLocation)
	defer file.Close()
	if err != nil {
		return
	}
	var buffer bytes.Buffer
	tempBuff, err := tempReader.GetPlainText()
	if err != nil {
		return
	}
	buffer.ReadFrom(tempBuff)
	contents = buffer.String()
	return
}

func getNorwayStatsUrl(url string) Stats {
	if err := download(strconv.Itoa(fileNumber)+".pdf", url); err != nil {
		panic(err)
	}
	convPdf, err :=  readPlainTextPDF(strconv.Itoa(fileNumber)+".pdf")
	if err != nil {
		panic(err)
	}
	splicepdf := strings.Fields(convPdf)
	getNorway := 0

	for i := 0; i < len(splicepdf); i++ {
		if splicepdf[i] == "Norway" {
			getNorway = i
			break

		}
	}
	fileNumber++
	tempStats := Stats{TotalInfected: splicepdf[getNorway+1], NewInfected: splicepdf[getNorway+2], TotalDeath:splicepdf[getNorway+3], NewDeath:splicepdf[getNorway+4] }

	return tempStats
}

func download(path string, url string) error {
	transport := &http.Transport{
		TLSClientConfig:&tls.Config{InsecureSkipVerify:true},
	}

	client := &http.Client{Transport:transport}

	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	genFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer genFile.Close()

	_, err = io.Copy(genFile, response.Body)
	return err
}