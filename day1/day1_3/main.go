package main

import (
	"GeekProject/day1/day1-3/server"
)

func main() {
	webServer := server.NewSdkHttpServer("mysql")
	//post
	webServer.Route("/user/signUp", server.SignUp)
	//log.Fatal(http.ListenAndServe("localhost:8080", nil))
	err := webServer.Start(":8080")
	if err != nil {
		return
	}
}
