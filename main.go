package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/labstack/echo/v4"
)

const ROUTER_URL string = "192.168.15.1"

type RouterInfo struct {
	Link *bool    `json:"link"`
	PPP  *bool    `json:"ppp"`
	Rx   *float64 `json:"rx"`
	Tx   *float64 `json:"tx"`
}

func parseRouterInfo(html string) RouterInfo {

	pppRegex := regexp.MustCompile(`(?m)^\s*var\s*pppStatus\s*=\s*'(.+)'`)
	ppp := pppRegex.FindStringSubmatch(html)[1] == "1"


	linkRegex := regexp.MustCompile(`(?m)^\s*var\s*gponUp\s*=\s*'(.+)'`)
	link := linkRegex.FindStringSubmatch(html)[1] == "1"


	opticalPowerRegex := regexp.MustCompile(`(?m)^\s*var\s*opticalPower\s*='(.+)'`)
	opticalLine := opticalPowerRegex.FindString(html)

	rxRegex := regexp.MustCompile(`RX:\s*(-?\d+\.?\d*)`)
	rxs := rxRegex.FindStringSubmatch(opticalLine)[1]
	rx, _ := strconv.ParseFloat(rxs, 64)

	txRegex := regexp.MustCompile(`TX:\s*(-?\d+\.?\d*)`)
	txs := txRegex.FindStringSubmatch(opticalLine)[1]
	tx, _ := strconv.ParseFloat(txs, 64)

	return RouterInfo{
		Rx: &rx,
		Tx: &tx,
		Link: &link,
		PPP: &ppp,
	}
}

func doRouterInfoRequest() (RouterInfo, error) {
	resp, err := http.Get("http://" + ROUTER_URL)

	if err != nil {
		fmt.Printf("%+v\n", err)
		return RouterInfo{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	// fmt.Printf("%s", body)
	d := parseRouterInfo(string(body))

	return d, nil
}

func getRouterInfo(c echo.Context) error {
	d, _ := doRouterInfoRequest()
	return c.JSON(http.StatusOK, d)
}

func main() {
	e := echo.New()
	e.GET("/api/router-info", getRouterInfo)
	e.Static("/", "static/")

	e.Logger.Fatal(e.Start(":9200"))

}
