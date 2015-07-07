package main

import (
	"strings"
)

var ProvinceExtends []string = []string{
	"壮族自治区",
	"回族自治区",
	"维吾尔自治区",
	"自治区",
	"特别行政区",
	"市",
	"省",
}

func GetProvinceName(str string) string {
	for _, item := range(ProvinceExtends) {
		str = strings.TrimSuffix(str, item)
	}

	return str
}

var CityExtends []string = []string{
	// "省直辖县级行政区划",
	// "自治区直辖县级行政区划",
	"朝鲜族自治州",
	"蒙古族藏族自治州",
	"藏族自治州",
	"藏族羌族自治州",
	"哈尼族彝族自治州",
	"彝族自治州",
	"土家族苗族自治州",
	"布依族苗族自治州",
	"依族苗族自治州",
	"苗族侗族自治州",
	"壮族苗族自治州",
	"傣族自治州",
	"白族自治州",
	"傣族景颇族自治州",
	"傈僳族自治州",
	"回族自治州",
	"蒙古自治州",
	"柯尔克孜自治州",
	"哈萨克自治州",
	"黎族苗族自治县",
	"黎族自治县",
	"地区",
	"县",
	"市",
	"盟",
	// "林区",
}

func GetCityName(str string) string {
	for _, item := range(CityExtends) {
		str = strings.TrimSuffix(str, item)
	}

	return str
}