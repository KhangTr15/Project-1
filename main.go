package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Welcome struct {
	Name string
	Time string
}

func main() {

	//Create a function that does this
	pemBytes, err := ioutil.ReadFile("/home/ubuntu/Documents/JOJO.pem")
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatalf("parse key failed:%v", err)
	}
	config := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", "ec2-18-191-194-235.us-east-2.compute.amazonaws.com:22", config)
	if err != nil {
		log.Fatalf("dial failed:%v", err)
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("session failed:%v", err)
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run("ls ")
	if err != nil {
		log.Fatalf("Run failed:%v", err)
	}

	//executing my script
	cmd := exec.Command("bash", "project-1.sh")
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err2 := cmd.Run()
	if err != nil {
		log.Fatal(err2)
	}
	fmt.Printf("Output \n", out.String())

	//running web server
	//Instantiate a Welcome struct object and pass in some random information.
	//We shall get the name of the user as a query parameter from the URL
	welcome := Welcome{"Anonymous", time.Now().Format(time.Stamp)}

	//We tell Go exactly where we can find our html file. We ask Go to parse the html file (Notice
	// the relative path). We wrap it in a call to template.Must() which handles any errors and halts if there are fatal errors

	templates := template.Must(template.ParseFiles("template/welcome.html"))

	//Our HTML comes with CSS that go needs to provide when we run the app. Here we tell go to create
	// a handle that looks in the static directory, go then uses the "/static/" as a url that our
	//html can refer to when looking for our css and other files.

	http.Handle("/static/", //final url can be anything
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")))) //Go looks in the relative "static" directory first using http.FileServer(), then matches it to a
	//url of our choice as shown in http.Handle("/static/"). This url is what we need when referencing our css files
	//once the server begins. Our html code would therefore be <link rel="stylesheet"  href="/static/stylesheet/...">
	//It is important to note the url in http.Handle can be whatever we like, so long as we are consistent.

	//This method takes in the URL path "/" and a function that takes in a response writer, and a http request.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//Takes the name from the URL query e.g ?name=Martin, will set welcome.Name = Martin.
		if name := r.FormValue("name"); name != "" {
			welcome.Name = name
		}
		//If errors show an internal server error message
		//I also pass the welcome struct to the welcome.html file.
		if err := templates.ExecuteTemplate(w, "welcome.html", welcome); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	//Start the web server, set the port to listen to 8080. Without a path it assumes localhost
	//Print any errors from starting the webserver using fmt
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
