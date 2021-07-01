package main

// Config is the struct for the config.json file
type Config struct {
	TwitchID   string `json:"twitchID"`
	CaptchaKey string `json:"captchaKey"`
	Proxy      string `json:"proxy"`
	Password   string `json:"password"`
}

// TwitchRegisterBody is the body submitted when registering an account
type TwitchRegisterBody struct {
	Username                string           `json:"username"`
	Password                string           `json:"password"`
	Email                   string           `json:"email"`
	Birthday                TwitchBirthday   `json:"birthday"`
	ClientID                string           `json:"client_id"`
	IncludeVerificationCode bool             `json:"include_verification_code"`
	Arkose                  TwitchFuncaptcha `json:"arkose"`
}

// TwitchBirthday filler
type TwitchBirthday struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

// TwitchFuncaptcha filler
type TwitchFuncaptcha struct {
	Token string `json:"token"`
}

// RandomUnameResponse the response from api.namefake.com
type RandomUnameResponse struct {
	Username string `json:"username"`
}

// FunCaptchaResponse is the response from the API giving back a valid token
type FunCaptchaResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// FunCaptchaResponse1 filler
type FunCaptchaResponse1 struct {
	ID       int    `json:"id"`
	Status   string `json:"status"`
	Solution string `json:"solution"`
}

// TwitchRegisterResponse filler
type TwitchRegisterResponse struct {
	OAuth string `json:"access_token"`
}

// TwitchGQLResponse GQL is basically the standard for twitches api
type TwitchGQLResponse []struct {
	Data struct {
		CurrentUser struct {
			ID              string `json:"id"`
			HasPrime        bool   `json:"hasPrime"`
			DisplayName     string `json:"displayName"`
			Email           string `json:"email"`
			IsEmailVerified bool   `json:"isEmailVerified"`
			Typename        string `json:"__typename"`
		} `json:"currentUser"`
		RequestInfo struct {
			CountryCode string `json:"countryCode"`
			Typename    string `json:"__typename"`
		} `json:"requestInfo"`
	} `json:"data"`
	Extensions struct {
		DurationMilliseconds int    `json:"durationMilliseconds"`
		OperationName        string `json:"operationName"`
		RequestID            string `json:"requestID"`
	} `json:"extensions"`
}

// PersistedQuery filler
type PersistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

// Extensions filler
type Extensions struct {
	PersistedQuery PersistedQuery `json:"persistedQuery"`
}
