package main

import (
	"log"
	"net/url"
	"os"
	"time"

	"cgv"
	"github.com/athurg/wechat_work_sdk"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
)

func main() {
	if _, ok := os.LookupEnv("TENCENTCLOUD_RUNENV"); ok {
		cloudfunction.Start(fetchAndNotify)
		return
	}

	fetchAndNotify()
}

func fetchAndNotify() {
	wechatClient := wechat.NewAgentClientFromEnv()

	cinema, err := cgv.CinemaByName("成都高新")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("影院信息%+v", cinema)

	screenDay := time.Now().Format("2006-01-02")

	//获取所有影片信息
	movies, err := cgv.MoviesByThat(cinema.THAT_CD)
	for _, movie := range movies {
		movie.MOV_NM, _ = url.QueryUnescape(movie.MOV_NM)

		picUrl := "http://www.cgv.com.cn" + movie.FILE_PATH
		detailUrl := "http://www.cgv.com.cn/movieDetail/gotoMovieDetail.fo?MOV_CD=" + movie.MOV_CD

		//获取影片排片日期信息，以确保是否当日正在上映
		dateList, err := cgv.ScheduleDateByCinemaMovie(cinema.THAT_CD, movie.MOV_CD)
		if err != nil {
			log.Println("\t", err)
			continue
		}

		var isTodayShow bool
		for _, date := range dateList {
			if date.SCN_DY == screenDay {
				isTodayShow = true
				break
			}
		}

		if !isTodayShow {
			continue
		}

		message := wechat.NewNewsMessage()
		message.Append(movie.MOV_NM, detailUrl, "", picUrl)

		//获取当日排片表
		schedules, err := cgv.ScheduleInfoToday(cinema.THAT_CD, movie.MOV_CD)
		if err != nil {
			message.Append(movie.MOV_NM, detailUrl, err.Error(), "")
			break
		}

		for _, s := range schedules {
			s.SCREEN_NM, _ = url.QueryUnescape(s.SCREEN_NM)
			s.SCN_FR_TM, _ = url.QueryUnescape(s.SCN_FR_TM)
			s.SCN_TO_TM, _ = url.QueryUnescape(s.SCN_TO_TM)
			message.Append(s.SCN_FR_TM+"～"+s.SCN_TO_TM+"，"+s.SCREEN_NM, "", "", "")
		}

		_, err = wechatClient.SendBatchNewsMessageToUsers(message, "@all")
		if err != nil {
			log.Println("发送消息错误", err)
		}
	}
}
