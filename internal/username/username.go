package username

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type SiteData struct {
	ErrorType string      `json:"errorType"`
	ErrorMsg  interface{} `json:"errorMsg"`
	URL       string      `json:"url"`
	URLMain   string      `json:"urlMain"`
	URLProbe  string      `json:"urlProbe"`
	URLError  string      `json:"errorUrl"`
	// TODO: Add headers
	UsedUsername   string `json:"username_claimed"`
	UnusedUsername string `json:"username_unclaimed"`
	RegexCheck     string `json:"regexCheck"`
}

type Result struct {
	Username string
	Exist    bool
	Site     string
	URL      string
	URLProbe string
	Link     string
}

type RequestError interface {
	Error() string
}

var (
	siteData  = map[string]SiteData{}
	waitGroup = &sync.WaitGroup{}
	outStream = new(strings.Builder)

	streamLogger = log.New(outStream, "", 0)
)

func SearchByUsername(username string) {
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Fprintf(
			color.Output,
			"[%s] Cannot open json file:  \"%s\"\n",
			color.HiRedString("!"), err.Error(),
		)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic("Failed to read file:" + "dataFileName")
	} else {
		err := json.Unmarshal(byteValue, &siteData)
		if err != nil {
			panic(err.Error())
		}
	}

	waitGroup.Add(len(siteData))
	guard := make(chan int, 1)
	for site := range siteData {
		guard <- 1
		go func(site string) {
			waitGroup.Done()

			res := Investigo(username, site, siteData[site])
			if res.Exist {
				fmt.Printf("[%s] %s: %s\n", "+", res.Site, res.Link)
			}
			// log.Println(res)
			<-guard
		}(site)
	}
	waitGroup.Wait()
}

func Investigo(username string, site string, data SiteData) Result {
	var u, urlProbe string
	var result Result
	u = strings.Replace(data.URL, "{}", username, 1)

	u = strings.Replace(data.URL, "{}", username, 1)

	// URL used to check if user exists.
	// Mostly same as variable `u`
	if data.URLProbe != "" {
		urlProbe = strings.Replace(data.URLProbe, "{}", username, 1)
	} else {
		urlProbe = u
	}

	r, err := Request(urlProbe)
	if err != nil {
		if r != nil {
			r.Body.Close()
		}
		return Result{
			Username: username,
			URL:      u,
			Exist:    false,
			Site:     site,
		}
	}

	switch data.ErrorType {
	case "status_code":
		if r.StatusCode == http.StatusOK {
			result = Result{
				Username: username,
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Exist:    true,
				Link:     u,
				Site:     site,
			}
		} else {
			result = Result{
				Username: username,
				URL:      data.URL,
				Site:     site,
				Exist:    false,
			}
		}
	case "message":
		switch errType := data.ErrorMsg.(type) {
		case string:
			if !strings.Contains(ReadResponseBody(r), data.ErrorMsg.(string)) {
				result = Result{
					Username: username,
					URL:      data.URL,
					URLProbe: data.URLProbe,
					Exist:    true,
					Link:     u,
					Site:     site,
				}
			} else {
				result = Result{
					Username: username,
					URL:      data.URL,
					Site:     site,
					Exist:    false,
				}
			}
		case []interface{}:
			_flag := false
			for _, msgString := range (data.ErrorMsg).([]interface{}) {
				if strings.Contains(ReadResponseBody(r), msgString.(string)) {
					_flag = true
					break
				}
			}
			result = Result{
				Username: username,
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Exist:    _flag,
				Link:     u,
				Site:     site,
			}
		default:
			panic(fmt.Sprintf("%s: Unsupported errorMsg type %T", site, errType))
		}
	case "response_url":
		// In the original Sherlock implementation,
		// the error type `response_url` works as `status_code`.
		if (r.StatusCode <= 300 || r.StatusCode < 200) && r.Request.URL.String() == u {
			result = Result{
				Username: username,
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Exist:    true,
				Link:     u,
				Site:     site,
			}
		} else {
			result = Result{
				Username: username,
				URL:      data.URL,
				Site:     site,
				Exist:    false,
			}
		}
	default:
		result = Result{
			Username: username,
			Exist:    false,
			Site:     site,
		}
	}

	r.Body.Close()

	return result
}

func Request(target string) (*http.Response, RequestError) {
	request, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")

	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	return client.Do(request)
}

func ReadResponseBody(response *http.Response) string {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return string(bodyBytes)
}
