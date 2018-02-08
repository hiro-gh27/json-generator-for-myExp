package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type broker struct {
	Type   string `json:"type"`
	Host   string `json:"host"`
	PeerID string `json:"peerID"`
	Num    int    `json:"number"`
}

type excPatterns struct {
	b []broker
}

var allpatterns []excPatterns

func initExcPatterns(input []broker, result []broker) {
	if len(input) <= 0 {
		p := excPatterns{b: result}
		allpatterns = append(allpatterns, p)
		return
	}
	for index := 0; index < len(input); index++ {
		var next []broker
		var rval []broker
		thisloop := input[index]

		next = append(next, input[:index]...)
		next = append(next, input[index+1:]...)
		rval = append(rval, result...)
		rval = append(rval, thisloop)

		initExcPatterns(next, rval)
	}
}

func main() {
	numberFlag := flag.Int("number", 0, "number of client")
	flag.Parse()
	clientNum := *numberFlag

	byteS, err := ioutil.ReadFile("../config/brokers.json")
	if err != nil {
		log.Fatal(err)
	}
	var brokers []broker
	if err := json.Unmarshal(byteS, &brokers); err != nil {
		log.Fatal(err)
	}

	if isExistsDir("../result") {
		fmt.Println("resultフォルダが存在します")
		//os.Exit(0)
	} else {
		fmt.Println("新規作成")
		if err := os.Mkdir("../result", 0777); err != nil {
			log.Fatal(err)
		}
	}

	//テストケース
	var null []broker
	initExcPatterns(brokers, null)
	for _, aPattern := range allpatterns {
		exportPattern := aPattern.b
		fileName := fmt.Sprint("../result/")
		for index := 0; index < len(exportPattern); index++ {
			psType := "sub"
			number := clientNum / 2
			if index == 0 {
				psType = "pub"
				number = clientNum
			}
			exportPattern[index].Host = "tcp://" + exportPattern[index].Host
			exportPattern[index].Type = psType
			exportPattern[index].Num = number
			fileName += fmt.Sprintf("[%s:%s×%d]", exportPattern[index].PeerID, exportPattern[index].Type, exportPattern[index].Num)
		}
		byteData := getByteJSON(exportPattern)
		fileName += ".json"
		export(fileName, byteData)
	}
}

func mkdir(dirName string) {

}

func isExistsDir(filename string) bool {
	f, err := os.Stat(filename)
	return err == nil && f.IsDir()
}

func export(fileName string, data []byte) {
	var file *os.File
	defer file.Close()
	var err error
	if file, err = os.Create(fileName); err != nil {
		log.Fatal(err)
	}
	file.Write(data)
}

func getByteJSON(cs []broker) []byte {
	var err error
	var jb []byte
	rByteBuffer := new(bytes.Buffer)

	if jb, err = json.Marshal(cs); err != nil {
		log.Fatal(err)
	}
	json.Indent(rByteBuffer, jb, "", "    ")
	return rByteBuffer.Bytes()
}
