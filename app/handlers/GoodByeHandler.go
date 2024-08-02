package handlers

import (
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type GoodBye struct {
	log *logrus.Logger
}

func NewGoodBuy(log *logrus.Logger) *GoodBye {
	return &GoodBye{log: log}
}

func (gb *GoodBye) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	gb.log.Info("handling good bye request")
	_, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, "some error occurred while reading request", 500)
	}
	resp.Write([]byte("Good bye to microservice world"))
}
