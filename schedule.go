package cgv

import (
	"net/url"
	"time"
)

type ScheduleDate struct {
	SCN_DY      string
	WEEK        string
	GRD_SCREENS string
}

func ScheduleDateByCinemaMovie(thatCode, movieCode string) ([]ScheduleDate, error) {
	p := url.Values{
		"THAT_CD":       {thatCode},
		"MOV_CD":        {movieCode},
		"SCREEN_GRD_CD": {"all"},
	}

	var info struct {
		DateList []ScheduleDate
	}
	err := request("GET", "/buy/getScheduleDate.fo", p, nil, &info)
	if err != nil {
		return nil, err
	}

	return info.DateList, nil
}

type Schedule struct {
	THAT_CD        string //影城代码 1030
	MOV_CD         string //10003235
	MOV_NM         string //影片名
	MOV_CPT_LAG_CD string //01
	MOV_FMT        string
	MOV_FMT_CD     string
	MOV_TYP_CD     string //影片类型代码
	MOV_TYP        string //影片类型

	SIMG_URL        string
	SCN_TM          string //110
	FILM_LAG        string //影片语种
	PRICEMULTI_YN   string
	SPECL_SCREEN_CD string
	SEAT_CNT        int //座位数量

	SCN_DY      string //放映日期2019-06-09
	SCN_FR_TM   string //放映开始时间17:45
	SCN_TO_TM   string //放映结束时间19:35
	SCN_SCH_SEQ int    //放映序号

	SCREEN_CD string //影厅代码
	SCREEN_NM string //影厅名

	MBR_CRD_PRC string //会员价格
	STD_PRC     int    //原价
	BKT_AMT     int    //正价

	SEAT_GRD_CD string //01
	BFLAG       string //1
}

//获取所有的影厅、放映时间
func ScheduleInfoByDay(thatCode, movieCode, scheduleDay string) ([]Schedule, error) {
	d := map[string]string{
		"THAT_CD": thatCode,
		"MOV_CD":  movieCode,
		"SCN_DY":  scheduleDay,
	}

	var info struct {
		ScrNmSet []string //影厅名列表
		ScheList []Schedule
	}

	err := request("POST", "/buy/getScheduleInfo.fo", nil, d, &info)
	if err != nil {
		return nil, err
	}

	return info.ScheList, nil
}

//获取当日指定影院指定影片的放映计划
func ScheduleInfoToday(thatCode, movieCode string) ([]Schedule, error) {
	return ScheduleInfoByDay(thatCode, movieCode, time.Now().Format("2006-01-02"))
}
