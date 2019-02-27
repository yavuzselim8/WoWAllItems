package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var allAdded = 0
var auctionableAdded = 0
func main(){
	file1, err := os.OpenFile("all.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not open file: %v\n", err)
	}
	file2, err := os.OpenFile("auctionable.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not open file: %v\n", err)
	}
	defer file1.Close()
	defer file2.Close()
	access := getAccesToken()
	t := 101*time.Millisecond
	i := 0
	for range time.Tick( t ) {
		go doRequest(i, access, file1, file2)
		if i % 1000 == 0{
			fmt.Fprintf( os.Stdout, "After %d sec, allAdded = %d, auctionableAdded = %d\n", i * 101 / 1000, allAdded, auctionableAdded )
		}
		if i == 200000{
			break
		}
		i++
	}
}

func doRequest(i int, accessToken string, file1 *os.File, file2 *os.File){
	httpClient := http.Client{}
	reqURL := fmt.Sprintf( "https://eu.api.blizzard.com/wow/item/%d?region=eu&locale=en_US", i)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil{
		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v\n", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %v\n", err)
	}
	defer res.Body.Close()
	if res.StatusCode == 200{
		str := fmt.Sprintf("%d\n", i)
		_, err := file1.WriteString( str )
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not write to the file: %v\n", err)
		}
		file1.Sync()
		allAdded++
		auctionable := gjson.Get(read(res), "isAuctionable").Bool()
		if auctionable == true {
			_, err := file2.WriteString( str )
			if err != nil {
				fmt.Fprintf(os.Stdout, "could not write to the file: %v\n", err)
			}
			file2.Sync()
			auctionableAdded++
		}
	}
}

func read(res *http.Response) string {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stdout, "cannot read body: %v\n", err)
	}
	bs := string(body)
	return bs
}

func getAccesToken() string{
	httpClient := http.Client{}
	reqURL := "https://567f3d8a2b0645cba5005436f5553d95:dGcR2udJoX8oeSqh7zpwc2KjPR6x9za4@eu.battle.net/oauth/token?grant_type=client_credentials"
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v\n", err)
	}
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %v\n", err)

	}
	defer res.Body.Close()
	accessToken := gjson.Get(read(res), "access_token").String()
	return accessToken
}
