package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/log"

	// "strings"
	// "strconv"
	// "math"
)

const (
	AccessKey    = "**"
	SecretKey    = "**"
	Server       = "http://pili.qiniuapi.com"
	DownloadPath = "/statd/downflow?begin=20180301000000&end=20180401000000&g=5min&select=flow&$domain="
	UpPath       = "/statd/upflow?begin=20180301000000&end=20180401000000&g=5min&select=flow&$domain="
)

func GetDownloadFlow(domain string) (peak DownloadFlow, err error) {
	url := fmt.Sprintf("%s%s", Server, DownloadPath+domain)
	resData, resErr := getRequest(url)
	if resErr != nil {
		fmt.Printf("connect http failure, %s\n", resErr.Error())
		err = resErr
		return
	}
	var res []DownloadFlow
	umErr := json.Unmarshal(resData, &res)
	if umErr != nil {
		log.Errorf("Unmarshal json failure, %s\n", umErr.Error())
		err = umErr
		return
	}

	//排序
	sort.Sort(FlowWrapper{res, func(p, q *DownloadFlow) bool {
		return q.Value.Flow < p.Value.Flow // flow 递减排序
	}})

	// fmt.Printf("count : %d\n", len(res))
	index := (((24*60)*30)/5)/20 + 1
	peak = res[index]
	return
}

func GetUpFlow(domain string) (peak DownloadFlow, err error) {
	url := fmt.Sprintf("%s%s", Server, UpPath+domain)
	resData, resErr := getRequest(url)
	if resErr != nil {
		fmt.Printf("connect http failure, %s\n", resErr.Error())
		err = resErr
		return
	}
	var res []DownloadFlow
	umErr := json.Unmarshal(resData, &res)
	if umErr != nil {
		log.Errorf("Unmarshal json failure, %s\n", umErr.Error())
		err = umErr
		return
	}

	//排序
	sort.Sort(FlowWrapper{res, func(p, q *DownloadFlow) bool {
		return q.Value.Flow < p.Value.Flow // flow 递减排序
	}})

	fmt.Printf("count : %d\n", len(res))
	//替换当月 天数
	index := (((24*60)*30)/5)/20 + 1
	peak = res[index]
	return
}

func GetLiveDomain() {
	url := fmt.Sprintf("%s%s", Server, "/v1/hubs")
	resData, resErr := getRequest(url)
	if resErr != nil {
		fmt.Printf("connect http failure, %s\n", resErr.Error())
		// err = resErr
		return
	}
	var res Domains
	umErr := json.Unmarshal(resData, &res)
	if umErr != nil {
		log.Errorf("Unmarshal json failure, %s\n", umErr.Error())
		// err = umErr
		fmt.Println(umErr.Error())
		return
	}
	f, err := os.Create("live_domain.txt")
	if err != nil {
		fmt.Printf("create map file error: %v\n", err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for i := 0; i < len(res.Items); i++ {
		fmt.Printf("%s\n", res.Items[i].Name)
		for j := 0; j < len(res.Items[i].Domains); j++ {
			// fmt.Printf("%s\n", res.Items[i].Domains[j].Domain)
			if res.Items[i].Domains[j].Type != "publishRtmp" {
				fmt.Fprintln(w, res.Items[i].Domains[j].Domain)
			}
		}
	}
	w.Flush()
}

func GetPublishDomain() {
	url := fmt.Sprintf("%s%s", Server, "/v1/hubs")
	resData, resErr := getRequest(url)
	if resErr != nil {
		fmt.Printf("connect http failure, %s\n", resErr.Error())
		// err = resErr
		return
	}
	var res Domains
	umErr := json.Unmarshal(resData, &res)
	if umErr != nil {
		log.Errorf("Unmarshal json failure, %s\n", umErr.Error())
		// err = umErr
		fmt.Println(umErr.Error())
		return
	}
	f, err := os.Create("publish_domain.txt")
	if err != nil {
		fmt.Printf("create map file error: %v\n", err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for i := 0; i < len(res.Items); i++ {
		fmt.Printf("%s\n", res.Items[i].Name)
		for j := 0; j < len(res.Items[i].Domains); j++ {
			// fmt.Printf("%s\n", res.Items[i].Domains[j].Domain)
			if res.Items[i].Domains[j].Type == "publishRtmp" {
				fmt.Fprintln(w, res.Items[i].Domains[j].Domain)
			}
		}
	}
	w.Flush()
}

func getRequest(url string) (data []byte, err error) {
	mac := qbox.NewMac(AccessKey, SecretKey)
	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		log.Errorf("get reqeust failure, ", reqErr)
		err = reqErr
		return
	}
	accessToken, _ := mac.SignRequest(req)
	req.Header.Set("Authorization", "QBox "+accessToken)
	client := &http.Client{}
	resp, respErr := client.Do(req)
	if respErr != nil {
		log.Errorf("connect failure,", respErr)
		err = respErr
		return
	}
	defer resp.Body.Close()
	body, rErr := ioutil.ReadAll(resp.Body)
	if rErr != nil {
		log.Errorf("read response body failure,", err)
		err = rErr
		return
	}
	if resp.StatusCode != 200 {
		log.Error("response failure, %s\n", string(body))
		errInfo := fmt.Sprintf("http status %d, %s\n", resp.StatusCode, string(body))
		err = errors.New(errInfo)
		return
	}
	data = body
	// fmt.Println(string(body))
	return
}