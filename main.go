package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

type Data struct {
	SourceID string `json:"source_id"`
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Amount int    `json:"amount"`
		Data   []Data `json:"data"`
	} `json:"data"`
}

func main() {
	log := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	logger := logrus.New()

	logger.Panic("test logrus")
	log.Println("test log")

}
