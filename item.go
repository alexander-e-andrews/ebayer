package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type item struct{
	Images []string //should be ordered
	Title string
	Description string
	ShipL float32
	ShipW float32
	ShipH float32
	ShipLb float32
	ShipOz float32
	Listed bool
	Sold bool
	Price float32
}

var count = 0

var items []item

func getList(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	//get a page of all the different items
}

func getItem(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	id := p.ByName("id")
	fmt.Println(id)
}

func updateItem(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	id := p.ByName("id")
	fmt.Println(id)
}

func newItem(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	count ++
	it := item{}
	items = append(items, it)
	url := fmt.Sprintf("/%d", count)
	http.Redirect(w, r, url, 301)
}