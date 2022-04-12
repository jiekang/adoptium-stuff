package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const API_BASE string = "https://api.adoptium.net/v3/stats/downloads/total"

var VERSIONS = []int{8, 11, 17}

type kv struct {
	Key   string
	Value int
}

func main() {

	vMap := make(map[int]map[string]int)
	for _, i := range VERSIONS {
		url := API_BASE + "/" + strconv.Itoa(i)
		r, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(err)
		}
		kMap := make(map[string]int)
		for key, _ := range data {
			url2 := url + "/" + key
			r, err := http.Get(url2)
			if err != nil {
				panic(err)
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}

			var data map[string]interface{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				panic(err)
			}
			for key, value := range data {
				d := int(value.(float64))
				// OpenJDK8U-debugimage_aarch64_linux_hotspot_8u302b08.tar.gz
				index := strings.Index(key, "_hotspot")
				if index != -1 {
					mk := strings.Split(key[:index], "-")[1]
					kMap[mk] = kMap[mk] + d
				}
			}
		}
		vMap[i] = kMap
	}

	for _, i := range VERSIONS {
		unsortedMap := vMap[i]

		var sortedMap []kv
		for k, v := range unsortedMap {
			sortedMap = append(sortedMap, kv{k, v})
		}

		sort.Slice(sortedMap, func(i, j int) bool {
			return sortedMap[i].Value > sortedMap[j].Value
		})

		fmt.Printf("Version: %v\n", i)
		for _, kv := range sortedMap {
			fmt.Printf("%s : %d\n", kv.Key, kv.Value)
		}
		fmt.Println()
	}
}
