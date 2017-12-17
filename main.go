package main

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "net/http"
  "fmt"
  "encoding/json"
  "strconv"
)

type Gopay struct {
  gorm.Model
  ExternalId string
  Type string
  Amount float32
  Passphrase string
}

type Result struct {
  Account Gopay
  Status string
}

func main() {
  http.HandleFunc("/", Serve)
  http.HandleFunc("/topup", ServeTopup)
  http.ListenAndServe(":8080", nil)
}

func connectDB() *gorm.DB {
  db, err := gorm.Open("postgres", "host=127.0.0.1 user=gopay_development dbname=gopay_development sslmode=disable password=123456")
  if err != nil {
    panic("failed to connect database")
  }
  db.AutoMigrate(&Gopay{})
  return db
}


func Serve(w http.ResponseWriter, r *http.Request) {
  db := connectDB()
  switch r.Method{
    case "POST":
      fmt.Println("Requested POST")
      r.ParseForm()
      fmt.Println(r)
      fmt.Println(r.Form)
      fmt.Println(r.Form.Get("id"))
      gopay := Gopay{
        ExternalId: r.Form.Get("id"),
        Type: r.Form.Get("type"),
        Amount: 0,
        Passphrase: r.Form.Get("passphrase"),
      }
      fmt.Println(gopay.ExternalId + gopay.Type)
      var result Result 
      if (gopay.ExternalId == "" || gopay.Type == "" || gopay.Passphrase == "") {
        result = Result{gopay, "FAILED"}
      } else {
        db.Create(&gopay)
        result = Result{gopay, "OK"}
      }
      fmt.Println(result.Status)
      bs, err := json.Marshal(result)

      if err != nil {
        fmt.Println(err)
      }
      w.Header().Set("Content-Type", "application/json")
      fmt.Fprintln(w, string(bs))

    //TODO : also make case for PATCH
    case "PUT":
      r.ParseForm()

      var gopay Gopay
      db.First(&gopay, "external_id = ? and passphrase = ? and type = ?", r.Form.Get("id"), r.Form.Get("passphrase"), r.Form.Get("type"))

      var result Result
      use,_ := strconv.ParseFloat(r.Form.Get("amount"), 32)

      if (gopay.ExternalId == "") {
        result = Result{gopay, "UNAUTHORIZED"}
      } else if (float32(use) > float32(gopay.Amount)) {
        result = Result{gopay, "INSUFFICIENT"}
      } else {
        db.Model(&gopay).Update("Amount", float32(gopay.Amount) - float32(use))
        result = Result{gopay, "OK"}
      }

      bs, err := json.Marshal(result)

      if err != nil {
        fmt.Println(err)
      }
      w.Header().Set("Content-Type", "application/json")
      fmt.Fprintln(w, string(bs))
  }
}

func ServeTopup(w http.ResponseWriter, r *http.Request) { 
  db := connectDB()
  switch r.Method{
    case "PUT":
      r.ParseForm()

      var gopay Gopay
      db.First(&gopay, "external_id = ? and passphrase = ?", r.Form.Get("id"), r.Form.Get("passphrase"))

      var result Result
      use,_ := strconv.ParseFloat(r.Form.Get("amount"), 32)

      if (gopay.ExternalId == "") {
        result = Result{gopay, "UNAUTHORIZED"}
      } else if (use <= 0) {
        result = Result{gopay, "INVALID"}
      } else {
        db.Model(&gopay).Update("Amount", float32(gopay.Amount) + float32(use))
        result = Result{gopay, "OK"}
      }

      bs, err := json.Marshal(result)

      if err != nil {
        fmt.Println(err)
      }
      w.Header().Set("Content-Type", "application/json")
      fmt.Fprintln(w, string(bs))
  }
}
