package main

import (
   "net/http"
   "github.com/RangelReale/osin"
   ex "github.com/RangelReale/osin/example"
   _ "github.com/lib/pq"
   "database/sql"
)


func main() {

   // ex.NewTestStorage implements the "osin.Storage" interface
   server := osin.NewServer(NewServerConfig(), ex.NewTestStorage())

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

   http.ListenAndServe(":14000", nil)
}
func existUser (user string, pass string) bool {
   conninfo := "user=postgres password=mysecretpassword dbname=postgres sslmode=disable"
   db, _ := sql.Open("postgres", conninfo)
   defer db.Close()
   rows, _ := db.Query(
      "SELECT * FROM users2 WHERE username=$1 and password=$2", user, pass )
   defer rows.Close()
   return rows.Next()
}

func NewServerConfig() *osin.ServerConfig {
   return &osin.ServerConfig{
      AuthorizationExpiration:   250,
      AccessExpiration:          3600,
      TokenType:                 "Bearer",
      AllowedAuthorizeTypes:     osin.AllowedAuthorizeType{osin.CODE},
      AllowedAccessTypes:        osin.AllowedAccessType{osin.PASSWORD},
      ErrorStatusCode:           200,
      AllowClientSecretInParams: true,
      AllowGetAccessRequest:     false,
      RetainTokenAfterRefresh:   false,
   }
}