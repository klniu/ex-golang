package main

import (
	"log"
	"bufio"
	"os"
	"fmt"
	"strings"
)

/*
CREATE TABLE `region` (
  `district_id` int(10) unsigned NOT NULL,
  `district_name` varchar(128) NOT NULL,
  `city_id` int(10) unsigned NOT NULL,
  `city_name` varchar(128) NOT NULL,
  `province_id` int(10) unsigned NOT NULL,
  `province_name` varchar(128) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/
func main() {
	fr, err := os.Open("regions-of-china.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer fr.Close()

	fw, err := os.Create("cities-of-china.sql")
	if err != nil {
		log.Fatal(err)
	}

	defer fw.Close()

	provinceName := ""
	cityName := ""
	isZX := false
	queryTpl := "('%s', '%s', '%s'),"

	scanner := bufio.NewScanner(fr)
	writer := bufio.NewWriter(fw)
	
	fmt.Fprintln(writer, "INSERT INTO ex_tag_new (`name`, `class`, `description`) VALUES")

	for scanner.Scan() {
		t := scanner.Text()

		if t[0] == '#' {
			continue
		}

		str := ""
		nodes := strings.Split(t, "　")
		switch len(nodes) {
		case 2: // Province
			provinceName = GetProvinceName(nodes[1])
		case 3: // City
			cityName = nodes[2]
			if cityName == "市辖区" {
				cityName = provinceName
			}
			if cityName == "县" {
				fmt.Printf("%s %s\n", cityName, provinceName)
				continue
			}
			
			cityName = GetCityName(cityName)
			if cityName == "" {
				fmt.Printf("%s %s\n", nodes[2], provinceName)
				continue
			}
			if cityName == "省直辖县级行政区划" || cityName == "自治区直辖县级行政区划" {
				isZX = true
				fmt.Printf("%s %s\n", nodes[2], provinceName)
				continue
			} else {
				isZX = false
			}

			str = fmt.Sprintf(queryTpl, cityName, "3", provinceName)
			fmt.Fprintln(writer, str)
		case 4: // District
			// 省直辖县级行政区划, 自治区直辖县级行政区划
			cityName = GetCityName(nodes[3])
			if cityName == "" {
				fmt.Printf("%s %s\n", nodes[3], provinceName)
				continue
			}

			if isZX == true {
				fmt.Printf("%s %s\n", cityName, provinceName)
				str = fmt.Sprintf(queryTpl, cityName, "3", provinceName)
				fmt.Fprintln(writer, str)
			}
		default: // Err
			fmt.Printf("$s\n", t)
		}
	}

	writer.Flush()
}

