package controllers

import (
	"gopkg.in/headzoo/surf.v1"
	//"gopkg.in/headzoo/surf.v2"
	"fmt"
)

func SendLight() {
	bow := surf.NewBrowser()
	//err := bow.Open("http://192.168.188.15")
	err := bow.Open("http://192.168.188.15/user/light.htm")
	checkPanic(err)

	// Outputs: "The Go Programming Language"
	fmt.Println(bow.Title())
	fmt.Println("", bow.Body())

	//LOGIN
	/*
		forms := bow.Forms()
		form := forms[1]
		form.Input("code", "93851851")
		fmt.Println(form.Action())
		err = form.Click("javascript:void(0);")
		//err = form.Click("login_submit")
		//err = form.Submit()
		checkPanic(err)
	*/

	//bow.Back()

	//err = bow.Open("http://192.168.188.15/user/light.htm")
	//checkPanic(err)

	//links := bow.Links()
	//link := links[6]
	//fmt.Println("", link)
	//err = bow.Click("light_link")

	fmt.Println("", bow.Body())

	err = bow.Click("a./user/light.htm")
	//err = bow.Open("user/light.htm")
	checkPanic(err)

}

func checkPanic(err error) {
	if err != nil {
		fmt.Println("", err.Error())
		panic(err)
	}
}
