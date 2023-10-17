package beurse

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
)

func GCFPostHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang"
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// func InsertUser(MONGOCONNSTRINGENV, dbname, collectionname string, userdata User) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	hash, _ := HashPassword(userdata.Password)
// 	userdata.Password = hash
// 	atdb.InsertOneDoc(mconn, collectionname, userdata)
// 	return "Ini username : " + userdata.Username + " ini password : " + userdata.Password
// }

func InsertUser(r *http.Request) string {
	var Response Credential
	var userdata User
		err := json.NewDecoder(r.Body).Decode(&userdata)
		if err != nil { 
			Response.Message = "error parsing application/json: " + err.Error() 
			return GCFReturnStruct(Response) 
		}
	hash, _ := HashPassword(userdata.Password)
	userdata.Password = hash
	atdb.InsertOneDoc(SetConnection("MONGOSTRING", "db_urse"), " user", userdata)
	return "Ini username : " + userdata.Username + "ini password : " + userdata.Password
}