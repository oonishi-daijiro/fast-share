package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// fnet http://example.com -d hoge
type directoryInfo struct {
	Pathes      []string `json:"pathes"`
	PackageName string   `json:"packageName"`
}

type argument struct {
	host       string
	reqDirName string
}

func main() {
	clArg, err := initArg()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	packageInfo, requestErr := requestPackageInfo(clArg.reqDirName, clArg.host)
	if requestErr != nil {
		fmt.Println("[ERROR]:" + requestErr.Error())
		return
	}
	createDirectoryStruct(packageInfo.Pathes)
	dlErr := downloadDirContent(packageInfo.Pathes, clArg.host)
	if dlErr != nil {
		fmt.Println("[ERROR]:" + dlErr.Error())
		return
	}
	fmt.Println("[STATS]:Done")
}

func initArg() (*argument, error) {
	if len(os.Args) != 4 {
		err := errors.New("Few or many argument")
		return nil, err
	}
	userOptions := new(argument)
	hostAndParms := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	reqDirName := hostAndParms.String("d", "", "The name of directory that you want to request")
	host := os.Args[1]
	if tf, err := isCorrectAdress(host); !tf {
		return nil, err
	}
	hostAndParms.Parse(os.Args[2:])
	userOptions.host = host
	userOptions.reqDirName = *reqDirName
	return userOptions, nil
}

func isCorrectAdress(adress string) (bool, error) {
	if !strings.HasPrefix(adress, "http://") {
		err := errors.New("It is not include \"http://\"")
		return false, err
	}
	return true, nil
}

func downloadDirContent(pathes []string, host string) error {
	fmt.Println("[STATS]:Download")
	for _, eachPath := range pathes {
		respFileContent, respErr := reqDirData(eachPath, host)
		if respErr != nil {
			return respErr
		}
		outputFile, createErr := os.OpenFile(eachPath, os.O_CREATE|os.O_RDWR, 0666)
		if createErr != nil {
			return createErr
		}
		defer outputFile.Close()
		_, cpErr := io.Copy(outputFile, respFileContent.Body)
		// fInfo, _ := outputFile.Stat()
		// *totalSize += fInfo.Size()
		if cpErr != nil {
			return cpErr
		}
	}
	return nil
}

func requestPackageInfo(dirName string, host string) (directoryInfo, error) { // first
	fmt.Println("[STATS]:Requesting directory information")
	requestParm, postParmErr := createPostParms("state", "requestDirInfo", "dirName", dirName)
	if postParmErr != nil {
		return directoryInfo{}, postParmErr
	}
	response, postErr := http.PostForm(host, requestParm)
	if postErr != nil {
		return directoryInfo{}, postErr
	}
	defer response.Body.Close()
	responseContent, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return directoryInfo{}, readErr
	}
	if response.StatusCode == 404 {
		notFound := errors.New("The directory was not found you specified")
		return directoryInfo{}, notFound
	} else if response.StatusCode == 403 {
		accessDenied := errors.New("Access denied")
		return directoryInfo{}, accessDenied
	}
	var information directoryInfo
	jsonParseErr := json.Unmarshal(responseContent, &information)
	if jsonParseErr != nil {
		return directoryInfo{}, jsonParseErr
	}
	return information, nil
}

func reqDirData(path string, host string) (*http.Response, error) { //second
	requestParm, createErr := createPostParms("state", "requestDirData", "dirName", path)
	if createErr != nil {
		return &http.Response{}, createErr
	}
	response, respErr := http.PostForm(host, requestParm)
	if respErr != nil {
		return &http.Response{}, respErr
	}
	return response, nil
}

func createPostParms(postParms ...string) (url.Values, error) {
	requestBody := url.Values{}
	if len(postParms)%2 != 0 {
		lengthErr := errors.New("The argument is not ehough")
		return nil, lengthErr
	}
	for i := 1; i <= len(postParms)/2; i++ {
		index := 2*i - 2
		requestBody.Add(postParms[index], postParms[index+1])
	}
	return requestBody, nil
}

func createDirectoryStruct(dirStruct []string) error {
	fmt.Println("[STATS]:Creating each directory")
	for _, eachPath := range dirStruct {
		if err := os.MkdirAll(filepath.Dir(eachPath), 0777); err != nil {
			return err
		}
	}
	return nil
}
