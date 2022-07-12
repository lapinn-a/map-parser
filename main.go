package main

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	_ "github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type YandexMap struct {
	Maps []Map `json:"maps"`
}

type Map struct {
	GeoObjects GeoObject `json:"geoObjects"`
}

type GeoObject struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Properties Property `json:"properties"`
}

type Property struct {
	Name    string `json:"iconCaption"`
	Contact string `json:"name"`
}

func main() {
	resp, err := http.Get("https://api-maps.yandex.ru/services/constructor/1.0/js/?um=constructor%3Ab81366596a932a75500027b805bfa2a3ef7781a0f2e731b10980559e52457cf0&lang=ru_RU")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	str := strings.Split(string(b), ")}),ym.modules.define(\"params\"")[0]
	str = strings.Split(str, "ym.modules.define(\"map-data\",[],function(e){e(")[1]

	var x YandexMap
	err = json.Unmarshal([]byte(str), &x)

	log.Print(x.Maps[0].GeoObjects.Features)

	f := excelize.NewFile()

	for i, feature := range x.Maps[0].GeoObjects.Features {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i), feature.Properties.Name)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i), feature.Properties.Contact)
	}

	if err := f.SaveAs("result.xlsx"); err != nil {
		log.Fatal(err)
	}
}
