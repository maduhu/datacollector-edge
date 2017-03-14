package common

import (
	"bufio"
	"log"
	"github.com/satori/go.uuid"
	"os"
)

const (
	SDE_ID_FILE = "data/sde.id"
)

type RuntimeInfo struct {
	id string
	logger *log.Logger
}

func (runtimeInfo *RuntimeInfo) init() error {
	runtimeInfo.id = runtimeInfo.getSdeId()
	return nil
}

func (runtimeInfo *RuntimeInfo) getId() string {
	return runtimeInfo.id
}

func (runtimeInfo *RuntimeInfo) getSdeId() string {
	var sdeId string
	if _, err := os.Stat(SDE_ID_FILE); os.IsNotExist(err) {
		f, err := os.Create(SDE_ID_FILE)
		check(err)

		defer f.Close()
		sdeId = uuid.NewV4().String()
		f.WriteString(sdeId)
	} else {
		file, err := os.Open(SDE_ID_FILE)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		sdeId = scanner.Text()
	}

	return sdeId
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NewRuntimeInfo(logger *log.Logger) (*RuntimeInfo, error) {
	runtimeInfo := RuntimeInfo{logger: logger}
	err := runtimeInfo.init()
	if err != nil {
		return nil, err
	}
	return &runtimeInfo, nil
}