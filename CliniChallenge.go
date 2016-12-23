package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/base64"
	"log"
    "os"
)

const (
	host	 = "https://stg.clinicloud.com"
	basePath = "/api/v2"
	clientId 	  = "X-CHALLENGE-APP"
	clientSecret  = "Y2xpbmljbG91ZGNoYWxsZW5nZXNlY3JldGtleQ=="
	email		  = "test@clinicloud.com"
	password      = "pass123"
)

type TokenStation struct {
    Type string
    Content struct {
    	Access_token string
    	Expires_in int64
    	Host string
    	New_user string
    	Refresh_token string
    	Scope string
    	Token_type string
    	Uuid string
    }
}

type UserInfo struct {
	Type string
    Content struct {
    	Id string
    	Mail string
    	First_name string
    	Last_name string
    	Dob float64
    	Gender int64
    	Country string
    	Password string
    	State string
    	NotifyRecord bool
    	Bio_details struct {
    		UseMetric bool
    		Weight int64
    		Height int64
    	}
    	PriDependents []PriDependent
    }
}

type PriDependent struct {
	Id string
    Mail string
    First_name string
    Last_name string
    Dob float64
    Gender int64
    Country string
    Password string
    State string
    NotifyRecord bool
    Bio_details struct {
    	UseMetric bool
    	Weight int64
    	Height int64
    }
}

type Items struct {
	UserID string
   	UpdateOrder int
}

type ClientCredentials []*Items

func main(){

	var chosen string
    var newEmail string
    var newPsw string
    var authorization string
    fmt.Println("Use Default Account(test@clinicloud.com)? yes/no")
    fmt.Scanln(&chosen)
    if ("no" == chosen) {
        fmt.Print("Email: ")
        fmt.Scanln(&newEmail)
        fmt.Print("Password: ")
        fmt.Scanln(&newPsw)
        authorization = newEmail + ":" + newPsw
    } else {
    	authorization = email + ":" + password
    }
    base64Authorization := base64.StdEncoding.EncodeToString([]byte(authorization))

// new instance client
	client := &http.Client{}

// login ....
	log.Println("login .... ")
	token := new(TokenStation)
	loginStatus, loginContent := login(token, client, base64Authorization)

// Retrieve user information .... 
	log.Println("Retrieving user information .... ")
	userInfo := new(UserInfo)
	userInfoStatus, userInfoContent := getUserInfo(token, userInfo, client)

// Retrieve sessions .... 
	log.Println("Retrieving sessions .... ")
	sessionStatus, sessionContent := getSession(token, client)

// generate html ....
	log.Println("Generating html .... ")
	generateHtml(loginStatus, loginContent, userInfoStatus, userInfoContent, sessionStatus, sessionContent)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func login (token interface{}, client *http.Client, authorization string) (string, string){
	url := host+basePath+"/login"
	log.Println("login url: " + url)
// assemble request
	credentials := map[string]string{
		"client_id": clientId,
		"client_secret": clientSecret}

	bodyJson, _ := json.Marshal(credentials)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("PushToken","")
	req.Header.Add("DeviceType","")
	req.Header.Add("DeviceToken","")
	req.Header.Add("Authorization", "Basic " + authorization)
// sent post request
	res, err := client.Do(req)
	check(err)
// Decode response body into json
    decodeError := json.NewDecoder(res.Body).Decode(token)
    check(decodeError)
// log response
    log.Println("login response Status:", res.Status)
    // log.Println("response Headers:", res.Header)
    log.Println("login response content: ")
    info, _ := json.MarshalIndent(token, "", "  ")
	log.Println(string(info))
    defer res.Body.Close()
    return res.Status, string(info)
}

func getUserInfo (token *TokenStation, userInfo interface{}, client *http.Client) (string, string){
	uuid := token.Content.Uuid
	url := host+basePath + "/user/" + uuid
	log.Println("getUserInfo url: " + url)
// assemble request
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("userID", uuid)
	req.Header.Add("Authorization", "Bearer " + token.Content.Access_token)
// sent get request
	res, err := client.Do(req)
	check(err)

// Decode response body into json
    decodeError := json.NewDecoder(res.Body).Decode(userInfo)
	if decodeError != nil {
		panic(decodeError)
	}
// log response
    log.Println("Retriev UserInfo response Status:", res.Status)
    // log.Println("response Headers:", res.Header)
    info, _ := json.MarshalIndent(userInfo, "", "  ")
    log.Printf("UserInfo content: \n%s", string(info))
    defer res.Body.Close()
    return res.Status, string(info)
}

func getSession (token *TokenStation, client *http.Client) (string, string){
	url := host+basePath+"/sessions/get"
	log.Println("getSession url:  " + url)
// assemble request
	log.Println("getSession assemble request")
	item := Items{UserID: token.Content.Uuid, UpdateOrder: 1 }
	clientCredentials := ClientCredentials{&item}
	bodyJson, _ := json.Marshal(clientCredentials)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + token.Content.Access_token)
// sent post request
	res, err := client.Do(req)
	check(err)
// log response
    log.Println("Retrieve sessions response Status:", res.Status)
    body, _ := ioutil.ReadAll(res.Body)
    var prettyJSON bytes.Buffer
    error := json.Indent(&prettyJSON, body, "", "\t")
    if error != nil {
    	log.Println("Retrieve sessions JSON parse error: ", error)
    	panic(error)
    }

    log.Printf("Retrieved sessions : \n%s", string(prettyJSON.Bytes()))
	defer res.Body.Close()
    return res.Status, string(prettyJSON.Bytes())
}

func generateHtml (loginStatus string, loginContent string, userInfoStatus string, 
    userInfoContent string, sessionStatus string, sessionContent string) string{
    htmlString1 := `<!DOCTYPE html>
<html>
<head>
<style>
div.container {
    width: 100%;
    border: 1px solid gray;
}
header, footer {
    padding: 1em;
    color: white;
    background-color: black;
    clear: left;
    text-align: center;
}
nav {
    text-align: center;
    float: left;
    max-width: 160px;
    margin: 0;
    padding: 1em;
}
article {
    margin-left: 170px;
    border-left: 1px solid gray;
    padding: 1em;
    overflow: hidden;
}
</style>
</head>
<body>
<div class="container">

<header><h1>CliniCloud Challenge Outcome</h1></header>
  
<nav>Login Infomation</nav>
<article>
<pre><code>`

htmlString2 := `</code></pre>
</article>
<hr>

<nav>User Infomation</nav>
<article>
<pre><code>`

htmlString3 := `</code></pre>
</article>
<hr>

<nav>Sessions</nav>
<article>
<pre><code>`

htmlString4 := `</code></pre>
</article>

<footer>Ethan (Xingyu Ji)</footer>
</div>
</body>
</html>`

    var httpBuffer bytes.Buffer
    httpBuffer.WriteString(htmlString1)
    httpBuffer.WriteString("Login response Status:" +loginStatus+"<br>")
    httpBuffer.WriteString(loginContent)
    httpBuffer.WriteString(htmlString2)
    httpBuffer.WriteString("Retrieve userInfo Status:" +userInfoStatus+ "<br>")
    httpBuffer.WriteString(userInfoContent)
    httpBuffer.WriteString(htmlString3)
    httpBuffer.WriteString("Retrieve sessions Status:" +sessionStatus + "<br>")
    httpBuffer.WriteString(sessionContent)
    httpBuffer.WriteString(htmlString4)
    httpContent := []byte(httpBuffer.String())
    f, err := os.Create("./Challenge_Outcome.html")
    check(err)
    err = ioutil.WriteFile("./Challenge_Outcome.html", httpContent, 0644)
    check(err)
    defer f.Close()
    log.Println("Render Html successful")
    return "success"
}