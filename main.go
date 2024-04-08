package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func primeNumberChecker(num int) bool {
	if num <= 1 {
		return false
	}

	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func randomRangeNum(rangeNum int) {
	/*
		Tolong buat satu array / list dari 1 sampai 100. Print semua angka ini dalam urutan terbalik, tetapi ada beberapa peraturan :

		1. Jangan print angka bilangan prima.

		2. Ganti angka yang dapat dibagi dengan angka 3 dengan text "Foo".

		3. Ganti angka yang dapat dibagi dengan angka 5 dengan text "Bar".

		4. Ganti angka yang dapat dibagi dengan angka 3 dan 5 dengan text "FooBar".

		5. Print angka menyamping tidak ke bawah.
	*/

	result := ""

	for i := rangeNum; i >= 1; i-- {
		if test := primeNumberChecker(i); test {
			continue
		}

		if i != rangeNum {
			result += ", "
		}

		if i%3 == 0 && i%5 == 0 {
			result += "FooBar"
		} else if i%3 == 0 {
			result += "Foo"
		} else if i%5 == 0 {
			result += "Bar"
		} else {
			result += strconv.Itoa(i)
		}

	}

	fmt.Print(result)
}

type Forecast struct {
	List []struct {
		DtTxt string `json:"dt_txt"`
		Main  struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	} `json:"list"`
}

func formatDate(dtTxt string) (string, error) {
	// Parse the date string into a time.Time object
	t, err := time.Parse("2006-01-02 15:04:05", dtTxt)
	if err != nil {
		return "", err
	}

	// Format the time into the desired format
	formattedDate := t.Format("Mon, 02 Jan 2006")
	return formattedDate, nil
}

func weather(apiKey string, city string, countryCode string) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?q=%s,%s&appid=%s&units=metric", city, countryCode, apiKey)

	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		fmt.Println("Error fetching data:", getErr)
		return
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	var forecast Forecast
	err = json.NewDecoder(res.Body).Decode(&forecast)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// After decoding JSON response into forecast variable
	if len(forecast.List) == 0 {
		fmt.Println("No forecast data available")
		return
	}

	prevDate := ""
	printedDays := 0
	for i := range forecast.List {

		dateTime := forecast.List[i].DtTxt
		temperature := forecast.List[i].Main.Temp

		formattedDate, err := formatDate(dateTime)
		if err != nil {
			fmt.Println("Error formatting date:", err)
			continue
		}

		if formattedDate != prevDate {
			printedDays++
			prevDate = formattedDate
			if printedDays > 5 {
				break
			}

			fmt.Printf("%s: %.1f°C\n", formattedDate, temperature)

		}

	}

}

func main() {
	// 1. Program kecil
	randomRangeNum(100)

	// 2. Menampilkan ramalan cuaca kota Jakarta untuk 5 hari kedepan
	const apiKey = "7cb2a32a387a4f96f0b7c4b9d65e1556"
	const city = "Jakarta"
	const countryCode = "ID"
	weather(apiKey, city, countryCode)
}
