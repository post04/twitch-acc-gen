package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// Conn is one of the vars in the global vars
var (
	Conn *websocket.Conn
	Mux  sync.Mutex
)

// Account is the struct that holds account data in the accounts.txt file
type Account struct {
	Username string
	Oauth    string
	Password string
	Email    string
}

func followThatMan(id, oAuth string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

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
		return ""
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b)
}

func unfollow(id, oauth string) {
	payload := strings.NewReader(fmt.Sprintf(`[{
	
		"operationName": "FollowButton_UnfollowUser",
		"variables": {
			"input": {
				"targetID":"%v"
			}
		},
		"extensions": {
			"persistedQuery": {
				"version": 1,
				"sha256Hash": "d7fbdb4e9780dcdc0cc1618ec783309471cd05a59584fc3c56ea1c52bb632d41"
			}
		}
	}]`, id))
	client := http.Client{}
	req, err := http.NewRequest("POST", "https://gql.twitch.tv/gql", payload)
	if err != nil {
		fmt.Println(err)
	}
	req.Host = "gql.twitch.tv"
	req.Header.Set("Connection", "close")
	req.Header.Set("Content-Length", "246")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Authorization", "OAuth "+oauth)
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
}

func getAllAccounts(accounts string) []*Account {
	var toReturn = []*Account{}
	active := false
	var email, password, oauth, username string
	for _, line := range strings.Split(accounts, "\n") {
		if line == "=====================" && active {
			toReturn = append(toReturn, &Account{
				Email:    email,
				Password: password,
				Oauth:    oauth,
				Username: username,
			})
			active = false
		} else if line == "=====================" && !active {
			active = true
		} else {
			parts := strings.Split(line, ": ")
			switch parts[0] {
			case "OAuth":
				oauth = parts[1]
			case "Username":
				username = parts[1]
			case "Password":
				password = parts[1]
			case "Email":
				email = parts[1]
			default:
				break
			}
		}
	}
	return toReturn
}

func printTemplate(accounts int) {
	fmt.Printf(`	Loaded %v accounts!
	
	Options:
		1.) Follow
		2.) UnFollow
		3.) MessageBot
	
	Please pick an option from 1 to 7: `, accounts)
}

func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func convert(thing string) int {
	converted, err := strconv.Atoi(thing)
	if err != nil {
		return 0
	}
	return converted
}

func newSocket(oauth, username string) {
	if c, _, err := websocket.DefaultDialer.Dial("wss://irc-ws.chat.twitch.tv:443/irc", nil); err == nil {
		Conn = c
	}

	login(oauth, username)
}

func writeRawString(content string) {
	Mux.Lock()
	Conn.WriteMessage(1, []byte(content))
	Mux.Unlock()
}

func login(oauth, username string) {
	writeRawString(fmt.Sprintf("PASS %v", oauth))
	writeRawString(fmt.Sprintf("NICK %v", username))
}

func join(channel string) {
	writeRawString(fmt.Sprintf("JOIN #%v", channel))
}

func sendmessage(message, streamer string) {
	writeRawString(fmt.Sprintf("PRIVMSG #%v :%v", streamer, message))
}

func main() {
	accounts, _ := ioutil.ReadFile("accounts.txt")
	accountsString := string(accounts)
	allAccounts := getAllAccounts(accountsString)

	var option string
	for {
		printTemplate(len(allAccounts))
		fmt.Scanln(&option)
		switch option {
		case "1":
			// 1 is for follow botting
			clearConsole()
			var id string
			var amount string
			fmt.Printf("	Userid: ")
			fmt.Scanln(&id)
			fmt.Printf("	Amount: ")
			fmt.Scanln(&amount)
			amountint := convert(amount)
			if amountint < 1 {
				fmt.Println("	Amount isn't an int or is below 1!")
				break
			}
			if amountint > len(allAccounts) {
				amountint = len(allAccounts)
			}
			wg := sync.WaitGroup{}
			for _, account := range allAccounts[:amountint] {
				wg.Add(1)
				go func(account *Account) {
					defer wg.Done()
					followThatMan(id, account.Oauth)
				}(account)
			}
			wg.Wait()
			fmt.Println("	Done!")
		case "2":
			// 2 is for unfollow botting
			clearConsole()
			var id string
			var amount string
			fmt.Printf("	Userid: ")
			fmt.Scanln(&id)
			fmt.Printf("	Amount: ")
			fmt.Scanln(&amount)
			amountint := convert(amount)
			if amountint < 1 {
				fmt.Println("	Amount isn't an int or is below 1!")
				break
			}
			if amountint > len(allAccounts) {
				amountint = len(allAccounts)
			}
			wg := sync.WaitGroup{}
			for _, account := range allAccounts[:amountint] {
				wg.Add(1)
				go func(account *Account) {
					defer wg.Done()
					unfollow(id, account.Oauth)
				}(account)
			}
			wg.Wait()
			fmt.Println("	Done!")

		case "3":
			// 3 is for message spam botting
			clearConsole()
			var streamerName string
			var amount string
			fmt.Print("	Streamer name: ")
			fmt.Scanln(&streamerName)
			fmt.Printf("	Amount: ")
			fmt.Scanln(&amount)
			amountint := convert(amount)
			if amountint < 1 {
				fmt.Println("	Amount isn't an int or is below 1!")
				break
			}
			if amountint > len(allAccounts) {
				amountint = len(allAccounts)
			}
			fmt.Print("	Message to send: ")
			in := bufio.NewReader(os.Stdin)
			messageToSend, _ := in.ReadString('\n')
			wg := sync.WaitGroup{}
			for _, account := range allAccounts[:amountint] {
				wg.Add(1)
				go func(account *Account) {
					defer wg.Done()
					newSocket("oauth:"+account.Oauth, account.Username)
					join(streamerName)
					sendmessage(messageToSend, streamerName)
				}(account)
			}
			wg.Wait()
			fmt.Println("	Done!")
		default:
			clearConsole()
			fmt.Println("	Incorrect option!")
			break
		}
	}
}
