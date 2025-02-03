package main

//imports
import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// global variables | PASSWORDS
var adminPass string = "testadmin" //ADMIN PASSWORD
var usrPass string = "testuser"    // USER PASSWORD

// global variables
var sdbgUsed bool = false //current fast debug session ID

// functions
func check(e error, panic_ bool) {
	/*
		Function to panic the server in case of error
		Parameters :
			(error) e : the error (nil if there's no error)
	*/
	if e != nil && panic_ {
		panic(e)
	} else if e != nil && !panic_ {
		fmt.Println(e)
	}
}

func httpCheck(e error, w_ http.ResponseWriter, msg string) {
	/*
		Function to log the errors into the console / panic the server in case of error
		Parameters :
			(error) e : the error (nil if there's no error)
			(http.ResponseWriter) w_ : http writer to send a response to the client
			(string) msg : the error message
	*/
	if e != nil {
		fmt.Fprint(w_, msg, http.StatusBadRequest)
	}
}

func read(filePath string) string {
	/*
		Function to read a file
		Parameters :
			(string) filePath : path to the file
		Returns :
			(string) : the content of the file
	*/
	data, err := os.ReadFile(filePath)
	check(err, false)
	return string(data)
}

func createDir(path string) {
	/*
		Function to create a directroy
		Parameters :
			(string) filePath : path to the file
	*/
	e := os.Mkdir(path, os.ModeDir)
	check(e, false)
}

func write(text string, file_ string, mode int) {
	/*
		Function to write a file
		Parameters :
			(string) text : the file content
			(string) file_ : path to the file
			(int) mode : flag of os to use to open the file
		Returns :
			(string) : the content of the file
	*/
	file, err_ := os.OpenFile(file_, os.O_CREATE|os.O_WRONLY|mode, 0600)
	check(err_, false)
	if _, err := file.WriteString(text); err != nil {
		panic(err)
	}
	e := file.Close()
	check(e, false)
}

func remove(file string) {
	/*
		Function to delete a file or directory
		Parameters :
			(string) file : path to the file

	*/
	err := os.RemoveAll(file)
	check(err, false)
}

func checkFileNotEmpty(path string) bool {
	/*
		Function to check if a file exists and isn't empty
		Parameters :
			(string) path : path to the file
		Returns :
			(bool) : True if the file exists and isn't empty, False otherwise
	*/
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		return false
	} else {
		return true
	}
}

func checkLogin(mode string, password string) bool {
	/*
		Function to check if the provided password is correct
		Parameters :
			(string) mode : used mode (admin or user)
			(string) password : the provided password
		Returns :
			(bool) : True if the pasword is correct, else, returns False
	*/
	var serverPassword string = ""

	//get the password for the provided mode
	if mode == "amdin" {
		serverPassword = adminPass
	} else {
		serverPassword = usrPass
	}
	//checking if it matches with the provided password
	if password == serverPassword {
		return true
	} else if password == adminPass {
		return true
	} else {
		return false
	}
}

func slowDebug(cmd_ string, dir string, out string) {
	/*
		Slow debug function | use it only combined with goroutines
	*/
	write("", "./cache/cmd/"+out, os.O_WRONLY)

	cmd := exec.Command("cmd", "/C", cmd_) //execute the command
	cmd.Dir = "./projects/" + dir

	res, err := cmd.Output()

	check(err, false)

	write(string(res), "./cache/cmd/"+out, os.O_WRONLY)
}

// requests handler
func handlerFunc(w http.ResponseWriter, r *http.Request) {
	/*
		Function to handle http/https requests
		Parameters :
			(http.ResponseWriter) w : a system to write the responses
			(*http.Request) r : a pointer that points to the request
	*/

	//Access-Control headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "GET" && r.URL.Path == "/test" {
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Server OK")

	} else if r.Method == "POST" && r.URL.Path == "/login" { //login
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")

		//reading r body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		httpCheck(err, w, "Unable to read request body")
		//defining JSON data structure
		type Login struct {
			Mode     string `json:"mode"`
			Password string `json:"password"`
		}
		//deconding the JSON
		var login Login
		e := json.Unmarshal(body, &login)
		httpCheck(e, w, "Unable to read the JSON data")

		var m string = login.Mode
		var p string = login.Password
		//action
		if checkLogin(m, p) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Password OK")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401, wrong password")
		}

	} else if r.Method == "POST" && r.URL.Path == "/listProjects" { //get the list of all the projects on the server
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")

		//reading r body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		httpCheck(err, w, "Unable to read request body")
		//defining JSON data structure
		type Login_ struct {
			Mode     string `json:"mode"`
			Password string `json:"password"`
		}
		//deconding the JSON
		var login_ Login_
		e := json.Unmarshal(body, &login_)
		httpCheck(e, w, "Unable to read the JSON data")
		//action
		if checkLogin(login_.Mode, login_.Password) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, read("./projects/__list__.txt"))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401, wrong password")
		}
	} else if r.Method == "POST" && r.URL.Path == "/listProjectBranches" { //get the list of files and folders (architectrue) of a project
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")

		//reading r body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		httpCheck(err, w, "Unable to read request body")
		//defining JSON data structure
		type Login struct {
			Mode     string `json:"mode"`
			Password string `json:"password"`
			Project  string `json:"project"`
		}
		//deconding the JSON
		var login Login
		e := json.Unmarshal(body, &login)
		httpCheck(e, w, "Unable to read the JSON data")
		//action
		var bfile string = "./projects/" + login.Project + "/__architecture__.txt"
		if checkLogin("", login.Password) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, read(bfile))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401, wrong password")
		}

	} else if r.Method == "POST" && r.URL.Path == "/read" { //read the content of an existing file
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")

		//reading r body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		httpCheck(err, w, "Unable to read request body")
		//defining JSON data structure
		type Login struct {
			Mode     string `json:"mode"`
			Password string `json:"password"`
			Project  string `json:"project"`
			Path     string `json:"path"`
		}
		//deconding the JSON
		var login Login
		e := json.Unmarshal(body, &login)
		httpCheck(e, w, "Unable to read the JSON data")
		//action
		var file string = "./projects/" + login.Project + "/" + login.Path
		var fl string = file + "locker.txt"
		if checkLogin(login.Mode, login.Password) && !checkFileNotEmpty(fl) {
			write(r.Header.Get("Origin"), fl, os.O_WRONLY)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, read(file))
		} else if checkFileNotEmpty(fl) {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, "409, file already in use")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401, wrong password")
		}

	} else if r.Method == "POST" && r.URL.Path == "/mkdir" { //create a new directory (folder)
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")

		//reading r body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		httpCheck(err, w, "Unable to read request body")
		//defining JSON data structure
		type Login struct {
			Mode     string `json:"mode"`
			Password string `json:"password"`
			Project  string `json:"project"`
			Path     string `json:"path"`
		}
		//deconding the JSON
		var login Login
		e := json.Unmarshal(body, &login)
		httpCheck(e, w, "Unable to read the JSON data")
		//action
		var dir string = "./projects/" + login.Project + "/" + login.Path
		if checkLogin(login.Mode, login.Password) {
			w.WriteHeader(http.StatusOK)
			createDir(dir)
			write("1", "./cache/mc2.txt", os.O_WRONLY)
			fmt.Fprintf(w, "Done")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401, wrong password")
		}
	} else if r.Method == "POST" && r.URL.Path == "/write" { //write into an existing file (must be locked for edition before use)
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")

		//reading r body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		httpCheck(err, w, "Unable to read request body")
		//defining JSON data structure
		type Login struct {
			Mode     string `json:"mode"`
			Password string `json:"password"`
			Project  string `json:"project"`
			Path     string `json:"path"`
			Content  string `json:"content"`
		}
		//deconding the JSON
		var login Login
		e := json.Unmarshal(body, &login)
		httpCheck(e, w, "Unable to read the JSON data")
		//action
		var file string = "./projects/" + login.Project + "/" + login.Path
		if checkLogin(login.Mode, login.Password) && read(file+"locker.txt") == r.Header.Get("Origin") && read(file+"locker.txt") != "" {
			w.WriteHeader(http.StatusOK)
			write(login.Content, file, os.O_WRONLY)
			fmt.Fprintf(w, "Done")
		} else if checkFileNotEmpty(file + "locker.txt") {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, "409 Conflict, file already in use")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401, wrong password")
		}

	} else if r.Method == "POST" && r.URL.Path == "/writef" { //create a new file
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")

		//reading r body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		httpCheck(err, w, "Unable to read request body")
		//defining JSON data structure
		type Login struct {
			Mode     string `json:"mode"`
			Password string `json:"password"`
			Project  string `json:"project"`
			Path     string `json:"path"`
			Content  string `json:"content"`
		}
		//deconding the JSON
		var login Login
		e := json.Unmarshal(body, &login)
		httpCheck(e, w, "Unable to read the JSON data")
		//action
		var file string = "./projects/" + login.Project + "/" + login.Path
		if checkLogin(login.Mode, login.Password) {
			w.WriteHeader(http.StatusOK)
			write(login.Content, file, os.O_WRONLY)
			write("1", "./cache/mc2.txt", os.O_WRONLY)
			fmt.Fprintf(w, "Done")
		} else if checkFileNotEmpty(file + "locker.txt") {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, "409 Conflict, file already in use")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401, wrong password")
		}

	} else if r.Method == "POST" && r.URL.Path == "/delete" { //delete a file or directory
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")

		//reading r body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		httpCheck(err, w, "Unable to read request body")
		//defining JSON data structure
		type Login struct {
			Mode     string `json:"mode"`
			Password string `json:"password"`
			Project  string `json:"project"`
			Path     string `json:"path"`
		}
		//deconding the JSON
		var login Login
		e := json.Unmarshal(body, &login)
		httpCheck(e, w, "Unable to read the JSON data")
		//action
		var file string = "./projects/" + login.Project + "/" + login.Path
		if checkLogin(login.Mode, login.Password) && !checkFileNotEmpty(file+"locker.txt") && login.Mode == "admin" {
			w.WriteHeader(http.StatusOK)
			remove(file)
			write("1", "./cache/mc2.txt", os.O_WRONLY)
			fmt.Fprintf(w, "Done")
		} else if checkFileNotEmpty(file+"locker.txt") && read(file+"locker.txt") == r.Header.Get("Origin") {
			w.WriteHeader(http.StatusOK)
			remove(file)
			write("1", "./cache/mc2.txt", os.O_WRONLY)
			fmt.Fprintf(w, "Done")
		} else if checkFileNotEmpty(file + "locker.txt") {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, "409 Conflict, file is used by another user")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401, wrong password")
		}

	} else if r.Method == "POST" && r.URL.Path == "/close" { //unlock a file wich was locked for edition
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")

		//reading r body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		httpCheck(err, w, "Unable to read request body")
		//defining JSON data structure
		type Login struct {
			Mode     string `json:"mode"`
			Password string `json:"password"`
			Proj     string `json:"project"`
			Path     string `json:"path"`
		}
		//deconding the JSON
		var login Login
		e := json.Unmarshal(body, &login)
		httpCheck(e, w, "Unable to read the JSON data")
		//action
		var fl string = "./projects/" + login.Proj + "/" + login.Path + "locker.txt"

		if checkLogin(login.Mode, login.Password) && checkFileNotEmpty(fl) {
			w.WriteHeader(http.StatusOK)
			remove(fl)
			fmt.Fprintf(w, "Done")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401, wrong password")
		}

	} else if r.Method == "POST" && r.URL.Path == "/fdebug" { //fast debug
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")

		//reading r body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		httpCheck(err, w, "Unable to read request body")
		//defining JSON data structure
		type Debug struct {
			Mode     string `json:"mode"`
			Password string `json:"password"`
			Project  string `json:"project"`
			Command  string `json:"cmd"`
		}
		//deconding the JSON
		var dbg Debug
		e := json.Unmarshal(body, &dbg)
		httpCheck(e, w, "Unable to read the JSON data")
		//action

		var dir string = "./projects/" + dbg.Project
		var cmd_ string = dbg.Command

		if checkLogin(dbg.Mode, dbg.Password) {

			cmd := exec.Command("cmd", "/C", cmd_) //execute the command
			cmd.Dir = dir

			res, err := cmd.Output()

			check(err, false)

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, string(res)+"\n")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401, wrong password")
		}

	} else if r.Method == "POST" && r.URL.Path == "/debug" { //dual or slow debug
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")

		//reading r body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		httpCheck(err, w, "Unable to read request body")
		//defining JSON data structure
		type Debug struct {
			Mode        string `json:"mode"`
			Password    string `json:"password"`
			Project     string `json:"project"`
			Path        string `json:"path"`
			SlowCommand string `json:"slow_cmd"`
			Command     string `json:"cmd"`
		}
		//deconding the JSON
		var dbg Debug
		e := json.Unmarshal(body, &dbg)
		httpCheck(e, w, "Unable to read the JSON data")
		//action
		var proj string = dbg.Project + "/" + dbg.Path

		var cmd string = dbg.Command
		var scmd string = dbg.SlowCommand

		if checkLogin(dbg.Mode, dbg.Password) && !sdbgUsed {
			sdbgUsed = true
			go slowDebug(cmd, proj, "OUT.txt")
			go slowDebug(scmd, proj, "sOUT.txt")

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Started.")
		} else if checkLogin(dbg.Mode, dbg.Password) && !sdbgUsed {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Slow/Dual debug is already being used.")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401, wrong password")
		}

	} else if r.Method == "POST" && r.URL.Path == "/debugRecall" { //retrieve the dual/slow debug results
		//Content-Type header
		w.Header().Set("Content-Type", "text/plain")
		//reading r body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		httpCheck(err, w, "Unable to read request body")
		//defining JSON data structure
		type Debug struct {
			Mode     string `json:"mode"`
			Password string `json:"password"`
			Session  string `json:"session"`
		}
		//deconding the JSON
		var dbg Debug
		e := json.Unmarshal(body, &dbg)
		httpCheck(e, w, "Unable to read the JSON data")

		if checkLogin(dbg.Mode, dbg.Password) {
			//Content-Type header
			w.Header().Set("Content-Type", "application/json")
			//checking the outputs
			var outpath string = "./cache/cmd/OUT.txt"
			var s_outpath string = "./cache/cmd/sOUT.txt"

			var out string = read(outpath)
			var s_out string = read(s_outpath)
			//making the return JSON
			var json = "{\"output\":\"" + out + "\"}" + "," + "{\"slow_output\":\"" + s_out + "\"}"
			//clear the cache & free the slow debugger
			if s_out != "" {
				remove(out)
				remove(s_out)
				sdbgUsed = false
			}
			//send the response to the client
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, json)

		} else {
			//Content-Type header
			w.Header().Set("Content-Type", "text/plain")

			//response
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401, wrong password")
		}

	} else if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "GET,POST,OPTIONS")
	}
}

// main code
func main() {
	// Folder where TLS certificates are stored
	certDir := "https"

	// Paths to cert.pem and key.pem (definition)
	certFile := filepath.Join(certDir, "cert.pem")
	keyFile := filepath.Join(certDir, "key.pem")

	// Check if cert.pem and key.pem exist and are not empty
	certExists := checkFileNotEmpty(certFile)
	keyExists := checkFileNotEmpty(keyFile)

	// Create a router for http requests
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerFunc)

	// If TLS certificates exist and are not empty, configure the HTTPS server
	if certExists && keyExists {
		fmt.Println("Starting HTTPS server on port 8443...")
		tlsConfig := &tls.Config{}
		server := &http.Server{
			Addr:      ":8443",
			Handler:   mux,
			TLSConfig: tlsConfig,
		}
		log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	} else {
		// Else, start the http server
		fmt.Println("Starting HTTP server on port 8080...")
		log.Fatal(http.ListenAndServe(":8080", mux))
	}
}
