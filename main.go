package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"strconv"
	"bufio"
)

type twRate struct {
	//currency string
	inCash, inSpot, inForward10, inForward30, inForward60, inForward90, inForward120, inForward150, inForward180 float64
	outCash, outSpot, outForward10, outForward30, outForward60, outForward90, outForward120, outForward150, outForward180 float64
}

var twRates map[string]twRate

func readFromText(){
	resp, err := http.Get("http://rate.bot.com.tw/xrt/fltxt/0/day")
	//resp, err := http.Get("http://rate.bot.com.tw/xrt/flcsv/0/day")
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	i:=0
	for scanner.Scan() {
		l := scanner.Text()
		if i == 0 {
			i+=1
		} else {
			line := twRateSplit(l)
			rateObj := toTwRateObj(line)
			twRates[line[0]] = rateObj
			//fmt.Printf("%s -  %f\n", rateObj.currency, rateObj.outCash)
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
func toTwRateObj(s []string) twRate{
	var rate twRate
	//rate.currency = s[0]

	rate.inCash, _ = strconv.ParseFloat(s[2], 64)
	rate.inSpot, _ = strconv.ParseFloat(s[3], 64)
	rate.inForward10, _ = strconv.ParseFloat(s[4], 64)
	rate.inForward30, _ = strconv.ParseFloat(s[5], 64)
	rate.inForward60, _ = strconv.ParseFloat(s[6], 64)
	rate.inForward90, _ = strconv.ParseFloat(s[7], 64)
	rate.inForward120, _ = strconv.ParseFloat(s[8], 64)
	rate.inForward150, _ = strconv.ParseFloat(s[9], 64)
	rate.inForward180, _ = strconv.ParseFloat(s[10], 64)
	rate.outCash, _ = strconv.ParseFloat(s[12], 64)
	rate.outSpot, _ = strconv.ParseFloat(s[13], 64)
	rate.outForward10, _ = strconv.ParseFloat(s[14], 64)
	rate.outForward30, _ = strconv.ParseFloat(s[15], 64)
	rate.outForward60, _ = strconv.ParseFloat(s[16], 64)
	rate.outForward90, _ = strconv.ParseFloat(s[17], 64)
	rate.outForward120, _ = strconv.ParseFloat(s[18], 64)
	rate.outForward150, _ = strconv.ParseFloat(s[19], 64)
	rate.outForward180, _ = strconv.ParseFloat(s[20], 64)
	return rate
}
func twRateSplit(s string) []string {
	var slices []string
	str := strings.TrimSpace(s)
	pos := strings.Index(str, " ")
	if pos == -1 {
		slices = append(slices, str)
	} else {
		leftCut := str[0:pos]
		rightCut := strings.TrimSpace(str[pos:])
		slices = append(slices, leftCut)
		slices = append(slices, twRateSplit(rightCut)...)
	}
	return slices
}


func main() {
	ticker := time.NewTicker(1 * 60 * time.Second)

	twRates = make(map[string]twRate, 0)
	readFromText()
	go func() {
		for range ticker.C {
			readFromText()
		}
	}()

	http.HandleFunc("/callback", botCallbackHandler)
	http.HandleFunc("/rate/", rateCallbackHandler)
	http.ListenAndServe(":8888", nil)
}
func botCallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, robot")
}
func rateCallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	cur := strings.TrimPrefix(r.URL.Path, "/rate/")
	_, ok := twRates[cur]
	if ok == true {
		fmt.Fprint(w, cur, " - ", twRates[cur].inCash, twRates[cur].outCash, twRates[cur].inSpot, twRates[cur].outSpot, "\n");
	} else {
		http.NotFound(w, r)
		return
	}
	//fmt.Fprint(w, "USD - ", twRates["USD"].inCash, twRates["USD"].outCash, "\n");
	//fmt.Fprint(w, "JPY - ", twRates["JPY"].inCash, twRates["JPY"].outCash, "\n");
}
