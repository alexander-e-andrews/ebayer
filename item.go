package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"strconv"

	"github.com/julienschmidt/httprouter"
)

type item struct {
	Images          []string //should be ordered
	Title           string
	Description     string
	ShipL           float64
	ShipW           float64
	ShipH           float64
	ShipWeight      float64
	Listed          bool
	Sold            bool
	Price           float64
	SizeDescription string

	ID int
	sync.Mutex
}

var count = 0

//Making this a pointer so it is easier to work with the mutex
var items []*item

func getList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//get a page of all the different items
	err := itemPageTemplate.ExecuteTemplate(w, "listPage.html", items)
	if err != nil {
		panic(err)
	}
}

func getItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.Atoi(p.ByName("id"))
	fmt.Println(id)
	fmt.Println(len(items))
	if id == len(items) {
		http.Redirect(w, r, "/newItem", 301)
		return
	}
	it := items[id]
	fmt.Println(it)
	err := itemPageTemplate.ExecuteTemplate(w, "itemPage.html", it)
	if err != nil {
		panic(err)
	}
}

// We can run into a race condition here, but we are not worried about it as no more than 1 or 2 people will ever use this
func updateItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.Atoi(p.ByName("id"))
	fmt.Println(id)

	r.ParseMultipartForm(2000)
	fmt.Println(r.FormValue("Title"))
	it := items[id]
	it.Lock()
	it.Title = r.FormValue("Title")
	it.Description = r.FormValue("Description")

	it.ShipL, _ = strconv.ParseFloat(r.FormValue("ShipL"), 64)
	it.ShipW, _ = strconv.ParseFloat(r.FormValue("ShipW"), 64)
	it.ShipH, _ = strconv.ParseFloat(r.FormValue("ShipH"), 64)
	it.ShipWeight, _ = strconv.ParseFloat(r.FormValue("ShipWeight"), 64)
	it.Price, _ = strconv.ParseFloat(r.FormValue("Price"), 64)

	it.Sold, _ = strconv.ParseBool(r.FormValue("Sold"))
	it.Listed, _ = strconv.ParseBool(r.FormValue("Listed"))

	it.SizeDescription= r.FormValue("SizeDescription")

	saveItemsToBackup()
	it.Unlock()
}

func newItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	it := item{Title: "Temp", ID: count, Images: make([]string, 0)}
	items = append(items, &it)
	fmt.Println("wWhat is going wrong")
	fmt.Println(it)
	fmt.Println(&it)
	url := fmt.Sprintf("/item/%d", count)
	count++

	http.Redirect(w, r, url, 301)

}

func uploadFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	//get a ref to the parsed multipart form
	m := r.MultipartForm

	id, _ := strconv.Atoi(m.Value["itemID"][0])

	it := items[id]
	it.Lock()
	//get the *fileheaders
	files := m.File["image"]
	fmt.Println(it)

	fmt.Println(files)
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		it.Images = append(it.Images, files[i].Filename)
		//create destination file making sure the path is writeable.
		dst, err := os.Create("./static/images/" + files[i].Filename)

		defer dst.Close()
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	saveItemsToBackup()
	it.Unlock()

}
