package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"math/rand"
	"net/http"
)

func main() {
	http.HandleFunc("/test", testWX)
	http.HandleFunc("/add", newImageHandler)
	http.HandleFunc("/random", randomHandler)
	_ = http.ListenAndServe(":8082", nil)
}

// test if wx backend http usable.
func testWX(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("WX backend test OK!"))
}

// handler add flower
func newImageHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name            string
		Baike_info      string
		Image_in_base64 string
	}
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)
	fmt.Println(string(body))
	if !existFlower(input.Name) {
		println("not existFlower")
		insertNewImage(input.Name, input.Baike_info, input.Image_in_base64)
	}
}

// search if exist in db
func existFlower(name string) bool {
	println("testing existFlower")
	println(name)
	db, err := sql.Open("postgres", "user=wx dbname=wx password=password sslmode=disable port=5438")
	if err != nil {
		println(err.Error())
	}
	sqlReq, err := db.Prepare("SELECT * FROM image WHERE name=$1")
	rows, err := sqlReq.Exec(name)
	if err != nil {
		println(err.Error())
	}
	exist := false
	num, err := rows.RowsAffected()
	println(num)
	if num > 0 {
		exist = true
	}
	db.Close()
	println("tested existFlower")
	println(exist)
	return exist
}

// insert new flower into db
func insertNewImage(name string, baikeInfo string, imageBase64 string) {
	db, _ := sql.Open("postgres", "user=wx dbname=wx password=password sslmode=disable port=5438")
	sqlReq, _ := db.Prepare("INSERT INTO image(name,baike_info,image_in_base64) VALUES($1, $2, $3)")
	_, err := sqlReq.Exec(name, baikeInfo, imageBase64)
	if err != nil {
		println(err.Error())
	}
	db.Close()
}

// get a random flower from db
func randomHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("postgres", "user=wx dbname=wx password=password sslmode=disable port=5438")
	flowers, err := db.Query("SELECT * FROM image")
	if err != nil {
		println(err.Error())
	}
	type flower struct {
		name  string
		baike string
		base  string
	}
	var flowerList []flower
	for flowers.Next() {
		var a_flower flower
		var userlessId int
		err = flowers.Scan(&userlessId, &a_flower.name, &a_flower.baike, &a_flower.base)
		fmt.Println(a_flower.name)
		flowerList = append(flowerList, a_flower)
	}
	var randFlowerJSON struct {
		Name         string `json:"name"`
		Baike_info   string `json:"baike_info"`
		Image_base64 string `json:"image_in_base64"`
	}
	randFlower := flowerList[rand.Intn(len(flowerList))]
	randFlowerJSON.Name = randFlower.name
	randFlowerJSON.Baike_info = randFlower.baike
	randFlowerJSON.Image_base64 = randFlower.base
	response, err := json.Marshal(randFlowerJSON)
	_, _ = w.Write(response)
}
