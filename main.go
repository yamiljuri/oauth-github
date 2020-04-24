package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

const clientID = "f9b65d732ae557be7722"
const clientSecret = "9b6f4ab3b30555468a17a05939cd8a9e9cfd26e9"

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
type UserResponse struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

func getAuthGithub(code string) (*OAuthAccessResponse, error) {
	var oAuthAccessResponse OAuthAccessResponse
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&oAuthAccessResponse); err != nil {
		return nil, err
	}
	return &oAuthAccessResponse, nil
}

func getUserAuthenticated(token string) (*UserResponse, error) {
	var userResponse UserResponse
	httpClient := http.Client{}
	requestUser, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	requestUser.Header.Set("Authorization", "token "+token)
	if err != nil {
		return nil, err
	}
	requestApiUser, err := httpClient.Do(requestUser)
	if err != nil {
		return nil, err
	}
	defer requestApiUser.Body.Close()
	if err := json.NewDecoder(requestApiUser.Body).Decode(&userResponse); err != nil {
		return nil, err
	}
	return &userResponse, nil
}

func main() {
	urlGithub := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=http://localhost:8080/oauth/redirect", clientID)
	exec.Command("open", urlGithub).Start()
	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		code := r.FormValue("code")

		oAuthAccessResponse, err := getAuthGithub(code)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		userLogged, err := getUserAuthenticated(oAuthAccessResponse.AccessToken)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		responseString := fmt.Sprintf("Id:%d\n\rNombre: %s\n\rUsuario: %s\n\r", userLogged.Id, userLogged.Name, userLogged.Login)
		w.WriteHeader(http.StatusFound)
		w.Write([]byte(responseString))
	})
	http.ListenAndServe(":8080", nil)
}
