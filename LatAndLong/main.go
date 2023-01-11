package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var clientId string = os.Getenv("CLIENT_ID")
	var clientSecret string = os.Getenv("CLIENT_SECRET")

	// fmt.Print("주소: ")

	// in := bufio.NewReader(os.Stdin)
	// address, err := in.ReadString('\n')
	// if err != nil {
	// 	log.Panic(err)
	// }

	// fmt.Println(address)

	f, err := excelize.OpenFile("address.xlsx")
	if err != nil {
		log.Panic(err)
		return
	}

	rows := f.GetRows("mapo")

	// excel file 생성
	newFile := excelize.NewFile()
	sheet := newFile.NewSheet("mapo")
	newFile.SetActiveSheet(sheet)
	// sheet 셀 설정
	newFile.SetCellValue("mapo", "A1", "Address")
	newFile.SetCellValue("mapo", "B1", "Latitude")
	newFile.SetCellValue("mapo", "C1", "Longtitude")

	for rowIdx, row := range rows {
		for _, address := range row {
			value := url.Values{}
			value.Set("query", address)

			naverMapUrl := &url.URL{
				Scheme:   "https",
				Host:     "naveropenapi.apigw.ntruss.com",
				Path:     "map-geocode/v2/geocode",
				RawQuery: value.Encode(),
			}

			request, err := http.NewRequest("GET", naverMapUrl.String(), nil)

			if err != nil {
				log.Panic(err)
			}

			request.Header.Add("X-NCP-APIGW-API-KEY-ID", clientId)
			request.Header.Add("X-NCP-APIGW-API-KEY", clientSecret)

			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				log.Panic(err)
			}

			defer response.Body.Close()

			if response.StatusCode != 200 {
				log.Fatal(response.StatusCode)
			}

			bytes, err := ioutil.ReadAll(response.Body)

			if err != nil {
				log.Panic(err)
			}

			var resParse map[string]interface{}
			e := json.Unmarshal(bytes, &resParse)

			if e != nil {
				log.Panic(e)
			}
			// fmt.Println(resParse)

			addresses := resParse["addresses"]
			addressResult := addresses.([]interface{})[0]

			lat := addressResult.(map[string]interface{})["y"]
			long := addressResult.(map[string]interface{})["x"]

			var cellIdx string = strconv.Itoa(rowIdx + 2)
			newFile.SetCellValue("mapo", "A"+cellIdx, address)
			newFile.SetCellValue("mapo", "B"+cellIdx, lat)
			newFile.SetCellValue("mapo", "C"+cellIdx, long)

			er := newFile.SaveAs("result.xlsx")
			if er != nil {
				log.Panic(er)
			}
		}
		fmt.Println()
	}
}
