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
	fr, err := os.Open("area-of-china.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer fr.Close()

	fw, err := os.Create("area-of-china.sql")
	if err != nil {
		log.Fatal(err)
	}

	defer fw.Close()

	i := 0
	provinceID := "0"
	provinceName := ""
	cityID := "0"
	cityName := ""
	queryTpl := "('%s', '%s', '%s', '%s', '%s', '%s'),"

	scanner := bufio.NewScanner(fr)
	writer := bufio.NewWriter(fw)
	
	fmt.Fprintln(writer, "INSERT INTO region (`district_id`, `district_name`, `city_id`, `city_name`, `province_id`, `province_name`) VALUES")

	for scanner.Scan() {
		i++
		t := scanner.Text()

		if t[0] == '#' {
			continue
		}

		str := ""
		nodes := strings.Split(t, "　")
		switch len(nodes) {
		case 2: // Province
			provinceID = nodes[0]
			provinceName = nodes[1]
			str = fmt.Sprintf(queryTpl, "0", "", "0", "", provinceID, provinceName)
		case 3: // City
			cityID = nodes[0]
			cityName = nodes[2]
			
			/*if nodes[2] == "县" {
				continue
			}
			if nodes[2] == "市辖区" {
				cityID = nodes[0]
				cityName = provinceName
			}*/

			str = fmt.Sprintf(queryTpl, "0", "", cityID, cityName, provinceID, provinceName)
		case 4: // District
			str = fmt.Sprintf(queryTpl, nodes[0], nodes[3], cityID, cityName, provinceID, provinceName)
		default: // Err
			fmt.Printf("$s\n", t)
		}

		fmt.Fprintln(writer, str)
	}

	writer.Flush()
}
