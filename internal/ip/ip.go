package ip

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/fatih/color"
)

const apiIP = "https://ipapi.co/"

type ipAPIStruct struct {
	IP                 string  `json:"ip"`
	Version            string  `json:"version"`
	City               string  `json:"city"`
	Region             string  `json:"region"`
	RegionCode         string  `json:"region_code"`
	Country            string  `json:"country"`
	CountryName        string  `json:"country_name"`
	CountryCode        string  `json:"country_code"`
	CountryCodeIso3    string  `json:"country_code_iso3"`
	CountryCapital     string  `json:"country_capital"`
	CountryTld         string  `json:"country_tld"`
	ContinentCode      string  `json:"continent_code"`
	InEu               bool    `json:"in_eu"`
	Postal             string  `json:"postal"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Timezone           string  `json:"timezone"`
	UtcOffset          string  `json:"utc_offset"`
	CountryCallingCode string  `json:"country_calling_code"`
	Currency           string  `json:"currency"`
	CurrencyName       string  `json:"currency_name"`
	Languages          string  `json:"languages"`
	CountryArea        float64 `json:"country_area"`
	CountryPopulation  int     `json:"country_population"`
	Asn                string  `json:"asn"`
	Org                string  `json:"org"`
}

func SearchByIP(ip string) {
	ipRegex(ip)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiIP+ip+"/json/", nil)
	req.Header.Set("User-Agent", "mosint")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	var res ipAPIStruct
	json.Unmarshal(body, &res)

	printIPAPI(res)
}

func ipRegex(ip string) {
	re := regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if !re.MatchString(ip) {
		fmt.Println("\nInvalid IP!")
		os.Exit(0)
	}
}

func printIPAPI(apiapi_result ipAPIStruct) {
	fmt.Println("IPAPI Results:")
	fmt.Println("|- IP:", color.YellowString(apiapi_result.IP))
	fmt.Println("|- City:", color.WhiteString(apiapi_result.City))
	fmt.Println("|- Region:", color.WhiteString(apiapi_result.Region))
	fmt.Println("|- Country:", color.WhiteString(apiapi_result.CountryName))
	fmt.Println("|- Country Code:", color.WhiteString(apiapi_result.CountryCode))
	fmt.Println("|- Timezone:", color.WhiteString(apiapi_result.Timezone))
	fmt.Println("|- Organization:", color.WhiteString(apiapi_result.Org))
	fmt.Println("|- ASN:", color.WhiteString(apiapi_result.Asn))
}
