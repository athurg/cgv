package cgv

import (
	"fmt"
	"net/url"
)

type Movie struct {
	MOV_CD      string
	GRD_SCREENS string
	MOV_NM      string
	MOV_FMT_CD  string
	CNT         string
	MOV_TYP_CD  string
	MOV_ENG_NM  string
	FILE_PATH   string
}

//返回指定影片的CGV官网页面URL
func MovieDetailUrl(movieCode string) string {
	return fmt.Sprintf("%s/movieDetail/gotoMovieDetail.fo?MOV_CD=", CgvAddr, movieCode)
}

func MoviesByThat(thatCode string) ([]Movie, error) {
	p := url.Values{
		"THAT_CD":       {thatCode},
		"SCREEN_GRD_CD": {"all"},
	}

	var info struct{ MovieList []Movie }
	err := request("GET", "/buy/getScheMoviesByThatCd.fo", p, nil, &info)
	if err != nil {
		return nil, err
	}

	for i, movie := range info.MovieList {
		info.MovieList[i].MOV_NM, _ = url.QueryUnescape(movie.MOV_NM)
		info.MovieList[i].FILE_PATH, _ = url.QueryUnescape(movie.FILE_PATH)
	}

	return info.MovieList, nil
}

type MovieSubject struct {
	COMM_CD_NM  string
	GENRE_FG_CD string
}

func MovieSubjects(movieCode string) ([]MovieSubject, error) {
	p := url.Values{"MOV_CD": {movieCode}}
	var info struct{ MOV_SUBJECTS []MovieSubject }
	err := request("GET", "/buy/getMoveSubjects.fo", p, nil, &info)
	if err != nil {
		return nil, err
	}

	return info.MOV_SUBJECTS, nil
}
