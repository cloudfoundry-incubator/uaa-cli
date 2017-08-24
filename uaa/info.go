package uaa

import (
	"encoding/json"
	"errors"
)

type UaaInfo struct {
	App uaaApp `json:"app"`
	Links uaaLinks `json:"links"`
	Prompts uaaPrompts `json:"prompts"`
	ZoneName string `json:"zone_name"`
	EntityId string `json:"entityID"`
	CommitId string `json:"commit_id"`
	Timestamp string `json:"timestamp"`
}

type uaaApp struct {
	Version string `json:"version"`
}

type uaaLinks struct {
	ForgotPassword string `json:"passwd"`
	Uaa string `json:"uaa"`
	Registration string `json:"register"`
	Login string `json:"login"`
}

type uaaPrompts struct {
	Username []string `json:"username"`
	Password []string `json:"password"`
}

func Info(context UaaContext) (UaaInfo, error) {
	infoBytes, err := UnauthenticatedGetter{}.Get(context, "info", "")
	if err != nil {
		return UaaInfo{}, err
	}

	info := UaaInfo{}
	err = json.Unmarshal(infoBytes,&info)
	if err != nil {
		return UaaInfo{}, parseError("/info", infoBytes)
	}
	return info, err
}

func requestError(url string) error {
	return errors.New("An unknown error occurred while calling " + url)
}

func parseError(url string, body []byte) error {
	errorMsg := "An unknown error occurred while parsing response from " + url + ". Response was " + string(body)
	return errors.New(errorMsg)
}

func unknownError() error {
	return errors.New("An unknown error occurred")
}