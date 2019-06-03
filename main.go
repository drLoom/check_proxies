package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"io"
	"net/url"
	"net/http"
	"strconv"
	"time"
	"flag"
	"log"
	"crypto/tls"
	. "github.com/logrusorgru/aurora"
)

type Proxy struct {
	ip   string
	port string
}

type CheckedProxy struct {
	Proxy
	httpCode         int
	responseDuration float32
	err              string
}

func main() {
	start := time.Now()

	var source string
	flag.StringVar(&source, "source", "", "Path to CSV file with proxies to check")
	var threads int
	flag.IntVar(&threads, "threads", 5, "Number of threads")
	var output string
	flag.StringVar(&output, "output", "", "Folder, result file will be stored to")
	var testUrl string
	flag.StringVar(&testUrl, "url", "", "URL to test proxy")
	var timeout int
	flag.IntVar(&timeout, "timeout", 5, "Timeout per request, s")
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Log output")

	flag.Parse()

	fmt.Println("source:", source)
	fmt.Println("threads:", threads)
	fmt.Println("output:", output)
	fmt.Println("url:", testUrl)
	fmt.Println("timeout:", timeout)
	fmt.Println("verbose:", verbose)

	proxies, _ := loadProxies(source)
	fmt.Println("Proxies to check", len(proxies))

	toCheck := make(chan Proxy)
	go func() {
		for k := range proxies {
			toCheck <- *proxies[k]
		}
	}()

	checked := make(chan CheckedProxy)
	for i := 0; i < threads; i++ {
		go func() {
			for p := range toCheck {
				checkedP := checkProxy(p, timeout, testUrl, verbose)
				checked <- checkedP
			}
		}()
	}

	checkedFile := output + "/" + "checked_proxies.csv"
	f, err := os.Create(checkedFile)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	w := csv.NewWriter(f)

	w.Write([]string{"IP", "Port", "HTTP Code", "Response Duration", "Error message"})

	for i := 0; i < len(proxies); i++ {
		cP := <-checked
		record := []string{cP.ip, cP.port, strconv.Itoa(cP.httpCode), fmt.Sprintf("%f", cP.responseDuration), cP.err}

		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	w.Flush()

	elapsed := time.Now().Sub(start)
	fmt.Println("Checked file: " + checkedFile)
	fmt.Println(Green("Finished in " + fmt.Sprintf("%f", elapsed.Seconds()) + " sec"))
}

func loadProxies(file string) (map[int]*Proxy, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvf := csv.NewReader(f)

	proxies := map[int]*Proxy{}
	idx := -1
	for {
		idx++

		row, err := csvf.Read()

		if idx == 0 {
			continue
		}

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return proxies, err
		}

		p := &Proxy{}

		p.ip = row[0]
		p.port = row[1]

		proxies[idx] = p

	}
}

func checkProxy(p Proxy, timeout int, tUrl string, verbose bool) (CheckedProxy) {
	start := time.Now()
	testUrl, _ := url.Parse(tUrl)

	proxyURL, _ := url.Parse("http://" + p.ip + ":" + p.port)

	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(time.Duration(timeout) * time.Second),
	}

	request, _ := http.NewRequest("HEAD", testUrl.String(), nil)

	response, err := client.Do(request)
	elapsed := time.Now().Sub(start)

	code := 0
	errMsg := ""

	if err != nil {
		code = 0
		errMsg = err.Error()
	} else {
		code = response.StatusCode
		errMsg = ""
	}

	if verbose {
		fmt.Println(p.ip+":"+p.port+", Code "+strconv.Itoa(code)+" elapsed "+fmt.Sprintf("%f", elapsed.Seconds())+" sec", Red(" "+errMsg))
	}

	return CheckedProxy{p, code, float32(elapsed.Seconds()), errMsg}
}
