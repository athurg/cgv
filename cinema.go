package cgv

import (
	"fmt"
)

//影院信息
type Cinema struct {
	EXCEPT_YN     string
	THAT_ADDR     string //影院地址
	CITY_CD       string //城市代码
	OPEN_YN       string //是否营业，Y或N
	SARFT_THAT_CD string //广电总局影院代码
	THAT_CD       string //影院内部代码
	ASSEMBLE_YN   string //建造中
	THAT_NM       string //影院名
}

//城市信息
type City struct {
	CITY_CD     string   //城市代码
	CITY_NM     string   //城市名
	CITY_EN_NM  string   //城市英文名
	CINEMACOUNT string   //影城数量
	CINEMALIST  []Cinema //城市的影院信息
}

//获取所有的影院列表
func CinemaList() ([]Cinema, error) {
	var info struct {
		CityList struct {
			CityThatInfoList []City
		}
	}

	err := request("GET", "/index/queryCityList.fo", nil, nil, &info)
	if err != nil {
		return nil, err
	}

	cinemas := make([]Cinema, 0)
	for _, cities := range info.CityList.CityThatInfoList {
		for _, cinema := range cities.CINEMALIST {
			cinemas = append(cinemas, cinema)
		}
	}

	return cinemas, nil
}

//获取指定名称的影院信息
func CinemaByName(name string) (*Cinema, error) {
	cinemas, err := CinemaList()
	if err != nil {
		return nil, err
	}

	for _, info := range cinemas {
		if info.THAT_NM == name {
			return &info, nil
		}
	}

	return nil, fmt.Errorf("未找到")
}
