package main

import (
   "net/http"
   "github.com/RangelReale/osin"
   ex "github.com/RangelReale/osin/example"
   _ "github.com/lib/pq"
   "database/sql"
   "fmt"
)



func main() {

   //const configurationUrl string = "http://config-service-san-pos-global-dev.appls.boae.paas.gsnetcloud.corp/master/adn360-front.yml"
   //   config, err := GetRemoteConfig(configurationUrl)
   config, err := GetLocalConfig("goauth.yaml")
   if err != nil {
      panic(err)
   }

   // ex.NewTestStorage implements the "osin.Storage" interface
   server := osin.NewServer(config, ex.NewTestStorage())

   // Access token endpoint
   http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
      resp := server.NewResponse()
      defer resp.Close()

      if ar := server.HandleAccessRequest(resp, r); ar != nil {
         if existUser(ar.Username, ar.Password) {
            ar.Authorized = true
         }
         server.FinishAccessRequest(resp, r, ar)
      }
      osin.OutputJSON(resp, w, r)
   })

   err = http.ListenAndServe(":14000", nil)
   if err != nil {
      panic(err)
   }
}
func existUser (user string, pass string) bool {
   conninfo := "user=postgres password=mysecretpassword dbname=postgres sslmode=disable"
   db, err := sql.Open("postgres", conninfo)
   defer db.Close()
   if err != nil {
      fmt.Println(err)
      return false
   }

   rows, err := db.Query(
      "SELECT * FROM users2 WHERE username=$1 and password=$2", user, pass )
   defer rows.Close()
   if err != nil {
      fmt.Println("[ERROR]", err)
      return false
   }

   return rows.Next()
}