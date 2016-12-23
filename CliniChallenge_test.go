package main
 
import (
    "testing"
    "net/http"
	"encoding/base64"
	"io/ioutil"
	"log"
)

func TestLoginValidStatus(t *testing.T){
  log.SetOutput(ioutil.Discard)
  println("Valid Test Status: Login")
  expectedStatus := "200 OK"
  client := &http.Client{}
  token := new(TokenStation)
  authorization := email + ":" + password
  base64Authorization := base64.StdEncoding.EncodeToString([]byte(authorization))
  actual, _ := login(token, client, base64Authorization)
  if actual != expectedStatus {
    t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedStatus, actual)
  }
}

func TestLoginValidContent(t *testing.T){
  log.SetOutput(ioutil.Discard)
  println("Valid Test Content: Login")
  expectedReturnType := "Access_token"
  client := &http.Client{}
  token := new(TokenStation)
  authorization := email + ":" + password
  base64Authorization := base64.StdEncoding.EncodeToString([]byte(authorization))
  login(token, client, base64Authorization)
  actual := token.Type
  if actual != expectedReturnType {
    t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedReturnType, actual)
  }
}

func TestLoginInvalidContent(t *testing.T){
  println("Invalid Test Content: Login")
  expectedReturnType := "error"
  client := &http.Client{}
  token := new(TokenStation)
  authorization := "wrong@account.com:notexistedpsw"
  base64Authorization := base64.StdEncoding.EncodeToString([]byte(authorization))
  login(token, client, base64Authorization)
  actual := token.Type
  if actual != expectedReturnType {
    t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedReturnType, actual)
  }
}

func TestLoginInvalidStatus(t *testing.T){
  log.SetOutput(ioutil.Discard)
  println("Invalid Test Status: Login")
  expectedStatus := "401 Unauthorized"
  client := &http.Client{}
  token := new(TokenStation)
  authorization := "wrong@account.com:notexistedpsw"
  base64Authorization := base64.StdEncoding.EncodeToString([]byte(authorization))
  actual, _ := login(token, client, base64Authorization)
  if actual != expectedStatus {
    t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedStatus, actual)
  }
}

func TestRetrieveUserInfoValid(t *testing.T){
  println("Valid Test: RetrieveUserInfo")
  expectedReturnType := "user_info"
  client := &http.Client{}
  token := new(TokenStation)
  authorization := email + ":" + password
  base64Authorization := base64.StdEncoding.EncodeToString([]byte(authorization))
  login(token, client, base64Authorization)

  userInfo := new(UserInfo)
  getUserInfo(token, userInfo, client)

  actual := userInfo.Type
  if actual != expectedReturnType {
    t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedReturnType, actual)
  }
}

func TestRetrieveUserInfoInvalid(t *testing.T){
  println("Invalid Test: RetrieveUserInfo")
  expectedReturnType := ""
  client := &http.Client{}
  token := new(TokenStation)
  authorization := "wrong@account.com:notexistedpsw"
  base64Authorization := base64.StdEncoding.EncodeToString([]byte(authorization))
  login(token, client, base64Authorization)

  userInfo := new(UserInfo)
  getUserInfo(token, userInfo, client)

  actual := userInfo.Type
  if actual != expectedReturnType {
    t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedReturnType, actual)
  }
}

func TestRetrieveSessionValid(t *testing.T){
  println("Valid Test: RetrieveSession")
  expectedReturnStatus := "200 OK"
  client := &http.Client{}
  token := new(TokenStation)
  authorization := email + ":" + password
  base64Authorization := base64.StdEncoding.EncodeToString([]byte(authorization))
  login(token, client, base64Authorization)

  actual, _ := getSession(token, client)
  
  if actual != expectedReturnStatus {
    t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedReturnStatus, actual)
  }
}

func TestRetrieveSessionInvalid(t *testing.T){
  println("Invalid Test: RetrieveSession")
  expectedReturnStatus := "401 Unauthorized"
  client := &http.Client{}
  token := new(TokenStation)
  authorization := "wrong@account.com:notexistedpsw"
  base64Authorization := base64.StdEncoding.EncodeToString([]byte(authorization))
  login(token, client, base64Authorization)
  actual, _ := getSession(token, client)
  if actual != expectedReturnStatus {
    t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedReturnStatus, actual)
  }
}

func TestGenerateHtml(t *testing.T){
  println("Valid Test: GenerateHtml")
  expected := "success"
  client := &http.Client{}
  token := new(TokenStation)
  authorization := email + ":" + password
  base64Authorization := base64.StdEncoding.EncodeToString([]byte(authorization))
  loginStatus, loginContent := login(token, client, base64Authorization)
  userInfo := new(UserInfo)
  userInfoStatus, userInfoContent := getUserInfo(token, userInfo, client)
  sessionStatus, sessionContent := getSession(token, client)
  actual := generateHtml(loginStatus, loginContent, userInfoStatus,
  	userInfoContent, sessionStatus, sessionContent)
  if actual != expected {
    t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
  }
}
