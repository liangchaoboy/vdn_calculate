package main

import (
	"fmt"
	"time"
	"server"
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"strconv"
)

func main() {

	//----------------- 获取 uid 下的所有域名 -----------
	// 获取域名
	//server.GetLiveDomain()
	//server.GetPublishDomain()

	 //----------------- 计算所有域名下的 上行 带宽峰值 ------------
	 //读取域名计算
	 fi, err := os.Open("publish_domain.txt")
	 if err != nil {
	 	fmt.Printf("Error: %s\n", err)
	 	return
	 }
	 defer fi.Close()

	 //写文件
	 f, err := os.Create("2018_0301-2018-0401_publish_domain.xls")
	 if err != nil {
	 	panic(err)
	 }
	 defer f.Close()

	 f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	 w := csv.NewWriter(f)
	 w.Write([]string{"域名", "取值点日期", "取值点时间", "带宽(bps)"})

	 br := bufio.NewReader(fi)
	 for {
	 	a, _, c := br.ReadLine()
	 	if c == io.EOF {
	 		break
	 	}
	 	domain := string(a)
	 	res, err := server.GetUpFlow(domain)
	 	if err != nil {
	 		fmt.Printf("%s get info failure, %s\n", domain, err)
	 	}
	 	bandwidth := float64(res.Value.Flow*8) / 300.0
	 	d, t := timeChange(res.Time)
	 	w.Write([]string{domain, d, t, strconv.FormatFloat(bandwidth, 'f', 2, 64)})
	 	fmt.Printf("%s %s : %f\n", d, t, bandwidth)
	 }
	 w.Flush()
}


func timeChange(toBeCharge string) (d, t string) {
	//获取本地location
	// toBeCharge := "2017-04-02T12:05:00+08:00"                    //待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	timeLayout := "2006-01-02T15:04:05+08:00"                       //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	d = fmt.Sprintf("%04d-%02d-%02d", theTime.Year(), theTime.Month(), theTime.Day())
	t = fmt.Sprintf("%02d:%02d:%02d", theTime.Hour(), theTime.Minute(), theTime.Second())
	return
}
