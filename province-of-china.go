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

	fw, err := os.Create("province-of-china.sql")
	if err != nil {
		log.Fatal(err)
	}

	defer fw.Close()

	provinceName := ""
	province := ""
	cities := ""
	queryTpl := "('%s', '%s', '%s', '%s'),"

	scanner := bufio.NewScanner(fr)
	writer := bufio.NewWriter(fw)
	
	fmt.Fprintln(writer, "DELETE FROM ex_tag WHERE `class` = 20;\nINSERT INTO ex_tag (`id`, `name`, `class`, `comment`) VALUES")

	for scanner.Scan() {
		t := scanner.Text()

		if t[0] == '#' {
			continue
		}

		str := ""
		nodes := strings.Split(t, "　")
		switch len(nodes) {
		case 2: // Province
			// Save Province
			if province != "" {
				province = fmt.Sprintf(province, cities)
				fmt.Fprintln(writer, province)
				cities = ""
			}

			// Save Province
			provinceName = nodes[1]
			str = fmt.Sprintf(queryTpl, "NULL", provinceName, "20", "%s")
			province = str
		case 3: // City
			if nodes[2] == "市辖区" || nodes[2] == "县" {
				cities += " "+provinceName
			} else {
				cities += " "+nodes[2]
			}
			
		case 4: // District
		default: // Err
			fmt.Printf("$s\n", t)
		}
	}

	writer.Flush()
}
