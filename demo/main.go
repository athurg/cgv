package main

import (
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"cgv"
	"github.com/athurg/wechat_work_sdk"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
)

var wechatClient *wechat.AgentClient
var screenDay = time.Now().Format("2006-01-02")

func main() {
	if _, ok := os.LookupEnv("TENCENTCLOUD_RUNENV"); ok {
		cloudfunction.Start(fetchAndNotify)
		return
	}

	fetchAndNotify()
}

func fetchAndNotify() {
	wechatClient = wechat.NewAgentClientFromEnv()

	cinema, err := cgv.CinemaByName("成都高新")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("影院信息%+v", cinema)


	//获取所有影片信息
	movies, err := cgv.MoviesByThat(cinema.THAT_CD)

	var wg sync.WaitGroup
	for _, movie := range movies {
		wg.Add(1)
		go func(cinemaCode string, m cgv.Movie) {
			defer wg.Done()
			checkMovie(cinemaCode, m)
		}(cinema.THAT_CD, movie)
	}
	wg.Wait()
}

func checkMovie(cinemaCode string, movie cgv.Movie) {
	movie.MOV_NM, _ = url.QueryUnescape(movie.MOV_NM)

	picUrl := "http://www.cgv.com.cn" + movie.FILE_PATH
	detailUrl := "http://www.cgv.com.cn/movieDetail/gotoMovieDetail.fo?MOV_CD=" + movie.MOV_CD

	//获取影片排片日期信息，以确保是否当日正在上映
	dateList, err := cgv.ScheduleDateByCinemaMovie(cinemaCode, movie.MOV_CD)
	if err != nil {
		log.Println("\t", err)
		return
	}

	var isTodayShow bool
	for _, date := range dateList {
		if date.SCN_DY == screenDay {
			isTodayShow = true
			break
		}
	}

	if !isTodayShow {
		return
	}

	message := wechat.NewNewsMessage()
	message.Append(movie.MOV_NM, detailUrl, "", picUrl)

	//获取当日排片表
	schedules, err := cgv.ScheduleInfoToday(cinemaCode, movie.MOV_CD)
	if err != nil {
		message.Append(movie.MOV_NM, detailUrl, err.Error(), "")
		return
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
