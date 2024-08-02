package handlers

import (
	"microservices/app/persistence"
	"net/http"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
)

type product struct {
	log *logrus.Logger
}

func NewProduct(log *logrus.Logger) *product {
	return &product{log: log}
}

func (p *product) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p.log.Println("product handler triggered")
	if req.Method == http.MethodGet {
		p.getProducts(rw, req)
		return
	}

	if req.Method == http.MethodPost {
		p.addProduct(rw, req)
		return
	}

	if req.Method == http.MethodPut {
		p.log.Println(":::::::::::inside put rquest::::::::::::")
		str := req.URL.Path
		regexp := regexp.MustCompile(`/([0-9]+)`)
		group := regexp.FindAllStringSubmatch(str, -1)
		p.log.Println("length of group is ", len(group))
		if len(group) != 1 {
			p.log.Println("Invalid URI more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.log.Println(group[0])
		if len(group[0]) != 2 {
			p.log.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := group[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.log.Println("Invalid URI unable to convert to numer", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.updateProduct(rw, req, id)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)

	//rw.Write(productListJson)
}

func (ph *product) getProducts(rw http.ResponseWriter, req *http.Request) {
	productsList := persistence.GetProducts()
	// marshling can be memory consuming if json object is too big
	//better to use encoder which is also faster than marshling
	//productListJson, err := json.Marshal(productsList)
	rw.Header().Add("Content-type", "application/json")
	err1 := productsList.ToJSON(rw)
	if err1 != nil {
		ph.log.Errorf("not able to convert in json %f", err1)
		http.Error(rw, "not able to convert in json", 500)
	}

}

func (ph *product) addProduct(rw http.ResponseWriter, req *http.Request) {
	ph.log.Println("Handle POST Product")

	prod := &persistence.Product{}

	err := prod.FromJSON(req.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	persistence.AddProduct(prod)
}

func (ph *product) updateProduct(rw http.ResponseWriter, req *http.Request, id int) {
	ph.log.Println("Handle Put Product")

	prod := &persistence.Product{}

	err := prod.FromJSON(req.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = persistence.UpdateProduct(prod, id)
	if err == persistence.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
