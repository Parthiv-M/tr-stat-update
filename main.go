package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/parthiv-m/tr-stat-update/types"
	"github.com/parthiv-m/tr-stat-update/utils"
)

func main() {

	// set logging config
	file, err := utils.OpenLogFile("./debug.log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	// load environment
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load .env file: %e", err)
	}

	// decide dev or prod environment
	var WIKI_BASE_URL string
	WIKI_BASE_URL = os.Getenv("DEV_WIKI_BASE_URL")
	if len(os.Args) > 1 && os.Args[1] == "production" {
		WIKI_BASE_URL = os.Getenv("PROD_WIKI_BASE_URL")
		fmt.Println("Running in production...")
	} else {
		fmt.Println("Running in development...")
	}

	// auth request
	paramMap := map[string]string{
		"grant_type": "password",
		"username":   os.Getenv("TRANSKRIBUS_USERNAME"),
		"password":   os.Getenv("TRANSKRIBUS_PASSWORD"),
		"client_id":  "processing-api-client",
		"scope":      "offline_access",
	}
	headerMap := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	tokenReq := utils.CreateCustomRequest(utils.TRANSKRIBUS_AUTH_API_URL, "POST", paramMap, headerMap, &http.Cookie{})
	tokenRes := tokenReq.MakeRequest()

	if tokenRes.GetResponseStatusCode() != 200 {
		fmt.Println("Error authorising with Transkribus")
		log.Fatal(tokenRes.GetResponseBody())
	}

	var responsedata types.AccessToken
	err = json.Unmarshal(tokenRes.GetResponseBody(), &responsedata)
	accessToken := responsedata.AccessToken

	// stats request
	headerMap = map[string]string{
		"Accept":        "application/json",
		"Authorization": "Bearer " + accessToken,
	}
	statRequest := utils.CreateCustomRequest(utils.TRANSKRIBUS_REST_API_URL+"/credits", "GET", nil, headerMap, &http.Cookie{})
	statRes := statRequest.MakeRequest()

	if statRes.GetResponseStatusCode() != 200 {
		fmt.Println("Error fetching data from Transkribus")
		log.Fatal(statRes.GetResponseBody())
	}

	var creditsData types.Credits
	err = json.Unmarshal(statRes.GetResponseBody(), &creditsData)

	fmt.Println("Credits fetched from Transkribus dashboard")

	// wiki page request
	page_url := WIKI_BASE_URL + "?action=query&prop=revisions&titles=Data:Wikimedia_OCR,_Transkribus_quota.tab&rvslots=*&rvprop=content&formatversion=2&format=json"

	wikiPageRequest := utils.CreateCustomRequest(page_url, "GET", nil, nil, &http.Cookie{})
	wikiPageRes := wikiPageRequest.MakeRequest()

	if wikiPageRes.GetResponseStatusCode() != 200 {
		fmt.Println("Error fetching requested page")
		log.Fatal(wikiPageRes.GetResponseBody())
	}

	var page types.Page
	err = json.Unmarshal(wikiPageRes.GetResponseBody(), &page)

	var wikiPageData types.WikiPageData
	json.Unmarshal([]byte(page.Query.Pages[0].Revisions[0].Slots.Main.Content), &wikiPageData)

	fmt.Println("Wiki page contents fetched")

	// populate the updated wiki page content into a new object
	var newWikiPageData types.WikiPageData
	newWikiPageData.License = wikiPageData.License
	newWikiPageData.Description = wikiPageData.Description
	newWikiPageData.Schema = wikiPageData.Schema
	newWikiPageData.Sources = wikiPageData.Sources
	newWikiPageData.Data = append(wikiPageData.Data, []interface{}{time.Now().UTC().Format("2006-01-02 15:04:05"), creditsData.OverallBalance})

	// bot tokens request
	token_url := WIKI_BASE_URL + "?action=query&meta=tokens&format=json&type=login"

	botTokenRequest := utils.CreateCustomRequest(token_url, "GET", nil, nil, &http.Cookie{})
	botTokenRes := botTokenRequest.MakeRequest()

	if botTokenRes.GetResponseStatusCode() != 200 {
		fmt.Println("Error obtaining login token")
		log.Fatal(botTokenRes.GetResponseBody())
	}

	var loginTokens types.LoginTokens
	err = json.Unmarshal(botTokenRes.GetResponseBody(), &loginTokens)

	loginCookie := strings.Split(strings.Split(botTokenRes.GetResponseHeaders().Get("Set-Cookie"), ";")[0], "=")[1]
	loginCookieName := strings.Split(strings.Split(botTokenRes.GetResponseHeaders().Get("Set-Cookie"), ";")[0], "=")[0]

	fmt.Println("Login tokens fetched")

	// bot login request
	paramMap = map[string]string{
		"action":     "login",
		"lgname":     os.Getenv("MW_BOT_USERNAME"),
		"lgpassword": os.Getenv("MW_BOT_PASSWORD"),
		"lgtoken":    loginTokens.Query.Tokens.Logintoken,
		"format":     "json",
	}

	// login cookie <wikiname>_session
	cookie := &http.Cookie{
		Name:     loginCookieName,
		Value:    loginCookie,
		Path:     "/",
		HttpOnly: true,
	}

	headerMap = map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	botLoginRequest := utils.CreateCustomRequest(WIKI_BASE_URL, "POST", paramMap, headerMap, cookie)
	botLoginRes := botLoginRequest.MakeRequest()

	if botLoginRes.GetResponseStatusCode() != 200 {
		fmt.Println("Error logging in with Bot credentials")
		log.Fatal(botLoginRes.GetResponseBody())
	}

	var loginData types.LoginData
	err = json.Unmarshal(botLoginRes.GetResponseBody(), &loginData)

	if loginData.Login.Result == "Failed" {
		fmt.Println("Bot login failed")
		log.Fatal("Bot login failed with ", paramMap)
	}
	fmt.Println("Bot login successful")

	botLoginCookie := strings.Split(strings.Split(botLoginRes.GetResponseHeaders().Get("Set-Cookie"), ";")[0], "=")[1]
	botLoginCookieName := strings.Split(strings.Split(botLoginRes.GetResponseHeaders().Get("Set-Cookie"), ";")[0], "=")[0]

	// bot cookie <wikiname>_BPsession
	bot_cookie := &http.Cookie{
		Name:     botLoginCookieName,
		Value:    botLoginCookie,
		Path:     "/",
		HttpOnly: true,
	}

	// csrfToken request
	csrf_url := WIKI_BASE_URL + "?action=query&meta=tokens&format=json"

	csrfTokenRequest := utils.CreateCustomRequest(csrf_url, "GET", nil, nil, bot_cookie)
	csrfTokenRes := csrfTokenRequest.MakeRequest()

	if csrfTokenRes.GetResponseStatusCode() != 200 {
		fmt.Println("Error fetching CSRF token")
		log.Fatal(csrfTokenRes.GetResponseBody())
	}

	var csrfToken types.CSRFTokens
	err = json.Unmarshal(csrfTokenRes.GetResponseBody(), &csrfToken)

	if csrfToken.Query.Tokens.CSRFtoken == "+\\" || csrfToken.Query.Tokens.CSRFtoken == "" {
		fmt.Println("Empty CSRF token")
		log.Fatal("CSRF token received: " + csrfToken.Query.Tokens.CSRFtoken)
	}

	fmt.Println("CSRF token obtained")

	jsonString, err := json.Marshal(newWikiPageData)
	editSummary := "Updating Transkribus quota as on " + time.Now().UTC().Format("2006-01-02 15:04:05")

	paramMap = map[string]string{
		"action":        "edit",
		"format":        "json",
		"title":         "Data:Wikimedia_OCR,_Transkribus_quota.tab",
		"bot":           "true",
		"summary":       editSummary,
		"text":          string(jsonString),
		"token":         csrfToken.Query.Tokens.CSRFtoken,
		"formatversion": "2",
	}

	headerMap = map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	editRequest := utils.CreateCustomRequest(WIKI_BASE_URL, "POST", paramMap, headerMap, bot_cookie)
	editRes := editRequest.MakeRequest()

	if editRes.GetResponseStatusCode() != 200 {
		fmt.Println("Error making edit to the page")
		log.Fatal(editRes.GetResponseBody())
	}

	var editResData types.EditResponse
	err = json.Unmarshal(editRes.GetResponseBody(), &editResData)

	if editResData.Edit.Result == "Success" {
		fmt.Println("Data:Wikimedia_OCR,_Transkribus_quota.tab updated successfully")
	} else {
		fmt.Println("Error updating page")
		log.Fatal(editResData)
	}
}
