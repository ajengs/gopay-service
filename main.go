package main

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "net/http"
  "fmt"
  "encoding/json"
  "strconv"
)

type OrderLog struct {
  gorm.Model
  ExternalId string `json:"external_id"`
  Status string     `json:"status"`
  Notes string      `json:"notes"`
  GopayID  uint
}

type Gopay struct {
  gorm.Model
  ExternalId string `json:"external_id"`
  Type string       `json:"type"`
  Amount float32    `json:"amount"`
  Passphrase string `json:"passphrase"`
}

type Result struct {
  Account Gopay     `json:"account"`
  Status string     `json:"status"`
}

func main() {
  http.HandleFunc("/", Serve)
  http.HandleFunc("/topup", ServeTopup)
  http.ListenAndServe(":8080", nil)
}

func Serve(w http.ResponseWriter, r *http.Request) {
  switch r.Method{
    case "POST":
      result := register(w,r)
      sendResponse(w, r, result)

    case "PUT":
      result := updateAmount(w, r, "substract")
      sendResponse(w, r, result)
    case "PATCH":
      result := updateAmount(w, r, "substract")
      sendResponse(w, r, result)
  }
}

func ServeTopup(w http.ResponseWriter, r *http.Request) { 
  switch r.Method{
    case "PUT":
      result := updateAmount(w, r, "add")
      sendResponse(w, r, result)
    case "PATCH":
      result := updateAmount(w, r, "add")
      sendResponse(w, r, result)
  }
}

func sendResponse(w http.ResponseWriter, r *http.Request, result Result) {
  bs, err := json.Marshal(result)

  if err != nil {
    fmt.Println(err)
  }
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintln(w, string(bs))
}

func register(w http.ResponseWriter, r *http.Request) Result {
  db := connectDB()
  r.ParseForm()
  gopay := Gopay{
    ExternalId: r.Form.Get("id"),
    Type: r.Form.Get("type"),
    Amount: 0,
    Passphrase: r.Form.Get("passphrase"),
  }

  var result Result
  if (gopay.ExternalId == "" || gopay.Type == "" || gopay.Passphrase == "") {
    result = Result{gopay, "FAILED"}
  } else {
    db.Create(&gopay)
    result = Result{gopay, "OK"}
  }
  return result
}

func updateAmount(w http.ResponseWriter, r *http.Request, act string) Result {
  db := connectDB()
  r.ParseForm()
  var gopay Gopay
  db.First(&gopay, "external_id = ? and passphrase = ? and type = ?", r.Form.Get("id"), r.Form.Get("passphrase"), r.Form.Get("type"))

  var result Result
  use,_ := strconv.ParseFloat(r.Form.Get("amount"), 32)

  if (gopay.ExternalId == "") {
    result = Result{gopay, "UNAUTHORIZED"}
  } else if (act == "substract" && float32(use) > float32(gopay.Amount)) {
    result = Result{gopay, "INSUFFICIENT"}
  } else if (act == "add" && use <= 0) {
    result = Result{gopay, "INVALID"}
  } else {
    if (act == "add"){
      db.Model(&gopay).Update("Amount", float32(gopay.Amount) + float32(use))
    } else { 
      db.Model(&gopay).Update("Amount", float32(gopay.Amount) - float32(use)) 
    }
    logOrder(w, r, act, gopay, use)
    result = Result{gopay, "OK"}
  }
  fmt.Println(result)
  return result
}

func logOrder(w http.ResponseWriter, r *http.Request, act string, g Gopay, use float64) {
  db := connectDB()
  
  order := OrderLog{
    ExternalId: r.Form.Get("order_id"),
    Status: r.Form.Get("order_status"),
    Notes: act + " " + strconv.Itoa(int(use)),
    GopayID: g.ID,
  }
  db.Create(&order) 
}

func connectDB() *gorm.DB {
  db, err := gorm.Open("postgres", "host=127.0.0.1 user=gopay_development dbname=gopay_development sslmode=disable password=123456")
  if err != nil {
    panic("failed to connect database")
  }
  db.AutoMigrate(&Gopay{})
  db.AutoMigrate(&OrderLog{})
  return db
}

