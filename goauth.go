package main

import (
   "net/http"
   "github.com/RangelReale/osin"
//   ex "github.com/RangelReale/osin/example"
   _ "github.com/lib/pq"
   "database/sql"
   "fmt"
   "log"
   "encoding/json"
   "github.com/ory/osin-storage/storage/postgres"
)

const conninfo = "user=postgres password=mysecretpassword dbname=postgres sslmode=disable"

var store osin.Storage

func main() {

   //const configurationUrl string = "http://config-service-san-pos-global-dev.appls.boae.paas.gsnetcloud.corp/master/adn360-front.yml"
   //   config, err := GetRemoteConfig(configurationUrl)
   config, err := GetLocalConfig("goauth.yaml")
   if err != nil {
      panic(err)
   }

   db, err := sql.Open("postgres", conninfo)
   if err != nil {
      panic(err)
   }

   store = postgres.New(db)
//   initPgStorage(store.(*postgres.Storage))

   // ex.NewTestStorage implements the "osin.Storage" interface
   server := osin.NewServer(config, store)

   // Access token endpoint
   http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
      resp := server.NewResponse()
      defer resp.Close()

      if ar := server.HandleAccessRequest(resp, r); ar != nil {
         if existUser(ar.Username, ar.Password) {
            ar.Authorized = true
         }
         ar.UserData = ar.Username
         server.FinishAccessRequest(resp, r, ar)
      }
      osin.OutputJSON(resp, w, r)
   })

   // Introspect
   http.HandleFunc("/introspect", func(w http.ResponseWriter, r *http.Request) {
      //TODO
      r.ParseForm()
      log.Println(r.Form)
      token := r.Form.Get("token")
      tokenInfo, err := getTokenInfo(token)
      if err != nil {
         w.WriteHeader(http.StatusInternalServerError)
         return
      }
      log.Println(tokenInfo)
      data, _ := json.Marshal(tokenInfo)
      fmt.Fprintf(w, string(data))
   })

   err = http.ListenAndServe(":14000", nil)
   if err != nil {
      panic(err)
   }
}
func initPgStorage(store *postgres.Storage) {
   store.CreateSchemas()
   client := osin.DefaultClient {Id: "1234", Secret: "aabbccdd", RedirectUri: "not_used"}
   store.CreateClient(&client)

}
type TokenInfo struct {
   Active bool `json:"active"`
   Client_id string `json:"client_id"`
   Username string `json:"username"`
   Scope string `json:"scope"`
   Exp int `json:"exp"`
}

func getTokenInfo(token string) (*TokenInfo, error) {
   accessData, err := store.LoadAccess(token)
   if err != nil {
      return nil, err
   }
   tokenInfo := &TokenInfo{
      Active: true,
      Client_id: "app",
      Username: accessData.UserData.(string),
      Exp: 1,
      Scope: "user.read",
   }
   return tokenInfo, nil
}

func existUser (user string, pass string) bool {
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