package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type geoResponse struct {
	IP          string
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	RegionCode  string `json:"region_code"`
	RegionName  string `json:"region_name"`
	City        string
	ZipCode     string `json:"zip_code"`
	TimeZone    string `json:"time_zone"`
	Latitude    float32
	Longitude   float32
	MetroCode   int `json:"metro_code"`
}

var mx sync.Mutex
var wg sync.WaitGroup

func main() {
	file, _ := os.Open("ips.csv")
	outfile, _ := os.Create("out.csv")
	outfile.WriteString("install_id,ip,city,country,region\n")
	body := csv.NewReader(file)

	records, _ := body.ReadAll()

	for i := range records {

		err := locate(records[i], outfile)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func locate(record []string, outfile *os.File) error {
	res, err := http.Get("http://www.freegeoip.net/json/" + record[1])
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(res.Body)

	var gr geoResponse

	err = json.Unmarshal(body, &gr)
	if err != nil {
		return err
	}

	for ii := range record {

		outfile.WriteString(record[ii] + ",")

	}

	outfile.WriteString(gr.City + "," + gr.CountryName + "," + gr.RegionName + "\n")

	return nil
}
