package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	captchaBaseURL = "https://bcsapi.xyz"
)

var (
	c Config
)

func getCaptchaKey() int {
	client := http.Client{}
	stringBody := fmt.Sprintf("page_url=https://www.twitch.tv&s_url=https://client-api.arkoselabs.com&site_key=E5554D43-23CC-1982-971D-6A2262A2CA24&access_token=%v", c.CaptchaKey)
	req, err := http.NewRequest("POST", captchaBaseURL+"/api/captcha/funcaptcha", strings.NewReader(stringBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var respJSON FunCaptchaResponse
	b, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, &respJSON)
	if err != nil {
		panic(err)
	}
	if respJSON.Status == "submitted" {
		return respJSON.ID
	}
	return 0
}

func getCaptchaToken(id int) string {
	timesDone := 0
	for {
		time.Sleep(3 * time.Second)
		timesDone++
		if timesDone >= 20 {
			return ""
		}
		resp, err := http.Get(captchaBaseURL + fmt.Sprintf("/api/captcha/%v?access_token=%v", id, c.CaptchaKey))
		if err == nil {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			var r FunCaptchaResponse1
			err = json.Unmarshal(body, &r)
			if err == nil {
				if r.Status != "pending" {
					return r.Solution
				}
			}
		}
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStringBytes generates a random string x letters long
func RandStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func getEmail() string {
	return RandStringBytes(10) + "@mailo.xyz"
}

func getUsername() string {
	resp, err := http.Get("https://api.namefake.com/english-united-states/random")
	if err != nil {
		return RandStringBytes(10) + "_34652"
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var r RandomUnameResponse
	err = json.Unmarshal(body, &r)
	return r.Username
}

func registerAccount(body *TwitchRegisterBody, client *http.Client) string {
	bodyMarshal, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "https://passport.twitch.tv/register", strings.NewReader(string(bodyMarshal)))
	if err != nil {
		fmt.Println(err)
		return "err at step 1"
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "err at step 2"
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return string(bodyBytes)
}

func getUserID(OAuth, clientID string, client *http.Client) string {
	b := []byte(`[{"operationName":"VerifyEmail_CurrentUser","variables":{},"extensions":{"persistedQuery":{"version":1,"sha256Hash":"f9e7dcdf7e99c314c82d8f7f725fab5f99d1df3d7359b53c9ae122deec590198"}}}]`)
	req, err := http.NewRequest("POST", "https://gql.twitch.tv/gql", strings.NewReader(string(b)))
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "OAuth "+OAuth)
	req.Header.Set("Client-Id", clientID)
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	var r TwitchGQLResponse
	l, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(l, &r)
	return r[0].Data.CurrentUser.ID

}

func saveAccount(toSave string) {
	content, _ := ioutil.ReadFile("followbot/accounts.txt")
	c := string(content)
	c += toSave
	ioutil.WriteFile("followbot/accounts.txt", []byte(c), 0064)
}

func setCaptchaBad(capKey int) {
	client := http.Client{}
	bd := fmt.Sprintf(`{"access_token": "%v"}`, c.CaptchaKey)
	req, _ := http.NewRequest("POST", "https://bcsapi.xyz/api/captcha/bad/"+fmt.Sprint(capKey), strings.NewReader(bd))
	client.Do(req)
}

func followThatMan(id, oAuth string, client *http.Client) string {
	body := strings.NewReader(fmt.Sprintf(`[{

    "operationName": "FollowButton_FollowUser",
    "variables": {
        "input": {
            "disableNotifications": false,
            "targetID": "%v"
        }
    },
    "extensions": {
        "persistedQuery": {
            "version": 1,
            "sha256Hash": "3efee1acda90efdff9fef6e6b4a29213be3ee490781c5b54469717b6131ffdfe"
        }
    }
}]`, id))
	req, err := http.NewRequest("POST", "https://gql.twitch.tv/gql", body)
	if err != nil {
		fmt.Println(err)
	}
	req.Host = "gql.twitch.tv"
	req.Header.Set("Connection", "close")
	req.Header.Set("Content-Length", "246")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Authorization", "OAuth "+oAuth)
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	req.Header.Set("Client-Id", "kimne78kx3ncx6brgo4mv6wki5h1ko")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://www.twitch.tv")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.twitch.tv/")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	f, _ := ioutil.ReadAll(resp.Body)
	return string(f)
}

func main() {
	f, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(f, &c)
	if err != nil {
		panic(err)
	}
	// make an infinate for loop that creates 5 captcha keys then makes a goroutine with those keys to generate an account
	for {
		var capKeys []int
		for i := 0; i < 5; i++ {
			key := getCaptchaKey()
			fmt.Println("Got captcha key", key)
			capKeys = append(capKeys, key)
		}
		wg := sync.WaitGroup{}
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(captchaKey int) {
				defer wg.Done()
				tlsConfig := &tls.Config{InsecureSkipVerify: true}
				proxyURL, err := url.Parse(c.Proxy)
				if err != nil {
					panic(err)
				}
				client := &http.Client{Transport: &http.Transport{
					TLSClientConfig: tlsConfig,
					Proxy:           http.ProxyURL(proxyURL),
				},
					Timeout: 60 * time.Second,
				}
				clientID := "kimne78kx3ncx6brgo4mv6wki5h1ko"
				captchaToken := getCaptchaToken(captchaKey)
				if captchaToken == "" {
					return
				}
				email := getEmail()
				username := getUsername()
				account := registerAccount(&TwitchRegisterBody{
					Username: username,
					Password: c.Password,
					Email:    email,
					Birthday: TwitchBirthday{
						Day:   12,
						Month: 2,
						Year:  1998,
					},
					ClientID:                clientID,
					IncludeVerificationCode: true,
					Arkose: TwitchFuncaptcha{
						Token: captchaToken,
					},
				}, client)
				fmt.Println(account)
				var r TwitchRegisterResponse
				er := json.Unmarshal([]byte(account), &r)
				if er != nil {
					fmt.Println("oauth resp err", er)
				}
				Oauth := r.OAuth
				followThatMan(c.TwitchID, Oauth, client)
				saveAccount(fmt.Sprintf("\n=====================\nUsername: %v\nPassword: %v\nEmail: %v\nOAuth: %v\n=====================", username, c.Password, email, Oauth))
			}(capKeys[i])
		}
		wg.Wait()
	}
}
