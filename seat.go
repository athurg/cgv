package cgv

type Seat struct {
	ISTIP         string //是否包含提示:1或0
	CONTENT       string //提示内容
	SEAT_GRD_PRCS []struct {
		STD_PRD     string
		SEAT_GRD_CD string
		BKT_AMT     string
	}

	TOTAL_COLS int //座位总列数
	TOTAL_ROWS int //座位总排数

	SeatData []struct {
		SEAT_LOC_NO      string //座位号00300501
		LEVELCODE        string //座位等级代码
		IBKT_STM_SEAT_CD string //座位标准码5101500101#02#01

		YCOORD    string //第几排（1开始不包含空列）
		XCOORD    string //第几号（1开始，不包含空行）
		ROWNUM    string //行号2
		COLUMNNUM string //列号1

		SEAT_GRD_CD string //": "01",
		STATUS      string //状态：可售Available、不可售Unavailable、线上专座Yellow、已售Locked
		GROUPCODE   string //编组码: "ONLINE003005",
	}
}

//获取影片座位信息，参数: 广电总局影厅编码、影城排片序号
func SeatInfoByDay(sarftThatCode, scheduleSeq string) (*Seat, error) {
	d := map[string]string{
		"SCN_SCH_SEQ":   scheduleSeq,
		"SARFT_THAT_CD": sarftThatCode,
	}

	var info Seat

	err := request("POST", "/buy/scheSeat.fo", nil, d, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
