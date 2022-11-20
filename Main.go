package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type pwdRequest struct {
	Password string `json:"password"`
}

type pwdResponse struct {
	Password    string `json:"password"`
	Occurrences string `json:"occurrences"`
	Sha1        string `json:"sha1"`
}

func main() {
	router()
}

func router() {
	http.HandleFunc("/check-password", checkPasswordHandler)

	http.ListenAndServe(":8080", nil)
}

func checkPasswordHandler(w http.ResponseWriter, req *http.Request) {

	var request pwdRequest
	json.NewDecoder(req.Body).Decode(&request)

	if req.Method != "POST" {
		http.Error(w, "Only HTTP POST Method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var response pwdResponse = checkPasswordService(request)

	bytesResponse, jsonMarshalErr := json.Marshal(response)
	if jsonMarshalErr != nil {
		http.Error(w, jsonMarshalErr.Error(), http.StatusInternalServerError)
		return
	}

	_, writeErr := w.Write(bytesResponse)
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		return
	}

}

func checkPasswordService(request pwdRequest) pwdResponse {
	var api string = "https://api.pwnedpasswords.com/range/"
	pwd := request.Password
	hash := sha1.New()
	hash.Write([]byte(pwd))
	result := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	response, err := http.Get(api + result[0:5])
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	orig := result
	result = result[5:]
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	passwords := strings.Split(string(contents), "\n")
	pwdMap := make(map[string]string)
	var currentSplit []string
	for i := 0; i < len(passwords); i++ {
		currentSplit = strings.Split(passwords[i], ":")
		pwdMap[currentSplit[0]] = currentSplit[1][0 : len(currentSplit[1])-1]
	}
	checkPwd := pwdMap[result]
	var httpResponse pwdResponse
	if pwdMap[result] != "" {
		httpResponse.Password = pwd
		httpResponse.Occurrences = checkPwd
		httpResponse.Sha1 = orig
	}
	return httpResponse
}
