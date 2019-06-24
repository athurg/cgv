package main

import (
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-http/cgv"
	"github.com/go-http/wechat_work"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
)

var wechatClient *wechat.AgentClient

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
	log.Printf("影片信息%+v", movies)

	wechatClient.SendTextToUsers("今日影片提醒", os.Getenv("RECEIVERS"))

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

	picUrl := cgv.WrapPath(movie.FILE_PATH)
	detailUrl := cgv.MovieDetailUrl(movie.MOV_CD)

	//获取影片排片日期信息，以确保是否当日正在上映
	dateList, err := cgv.ScheduleDateByCinemaMovie(cinemaCode, movie.MOV_CD)
	if err != nil {
		log.Println("\t", err)
		return
	}

	screenDay := time.Now().Format("2006-01-02")
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

	//获取当日排片表
	schedules, err := cgv.ScheduleInfoToday(cinemaCode, movie.MOV_CD)
	if err != nil {
		//message.Append(movie.MOV_NM, detailUrl, err.Error(), "")
		return
	}

	var screenTime string
	times := make([]string, 0, len(schedules))
	for _, s := range schedules {
		s.SCN_FR_TM, _ = url.QueryUnescape(s.SCN_FR_TM)
		screenTime, _ = url.QueryUnescape(s.SCN_TM)
		times = append(times, s.SCN_FR_TM)
	}

	detail := "播放时长：" + screenTime + "分钟"
	detail += "\n上映时间：" + strings.Join(times, "、")
	message := wechat.NewNewsMessage()
	message.Append(movie.MOV_NM, detailUrl, detail, picUrl)
	_, err = wechatClient.SendBatchNewsMessageToUsers(message, os.Getenv("RECEIVERS"))
	if err != nil {
		log.Println("发送消息错误", err)
	}
}
