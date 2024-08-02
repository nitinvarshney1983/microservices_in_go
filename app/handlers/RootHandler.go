package handlers

import (
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type RootHandler struct {
	log *logrus.Logger
}

func NewRootHandler(l *logrus.Logger) *RootHandler {
	return &RootHandler{l}
}

func (rh *RootHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	rh.log.Info("handling request for root handler")

	reqBodyByteArr, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, "Oops some error is there", http.StatusBadRequest)
		return
	}
	str := string(reqBodyByteArr)
	str = "Hi " + str
	resp.Write([]byte(str))
	//fmt.Fprintf(resp, "Request body is sent with the request is %s", string(reqBodyByteArr))

}
