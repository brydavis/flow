package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

type Results struct {
	Query,
	Data string
}

type Response map[string]string

func ListenAndServe(port int, conn *sql.DB) error {
	http.HandleFunc("/run/", func(res http.ResponseWriter, req *http.Request) {
		RunHandler(res, req, conn)
	})
	http.HandleFunc("/save/", SaveHandler)
	http.HandleFunc("/open/", OpenHandler)
	http.HandleFunc("/static/", StaticHandler)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		RootHandler(res, req, conn)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		return err
	}
	return nil
}

func StaticHandler(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, req.URL.Path[1:])

}

func RootHandler(res http.ResponseWriter, req *http.Request, conn *sql.DB) {
	// b, _ := ioutil.ReadFile("views/index.html")
	b, _ := ioutil.ReadFile("static/demo/gosublime.html")

	t := template.New("")
	t, _ = t.Parse(string(b))

	t.Execute(res, Results{Query: "", Data: ""})

}

func SaveHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		filepath := req.FormValue("code_filepath")
		results := req.FormValue("res_filepath")

		if err := ioutil.WriteFile("./saved/"+filepath, []byte(req.FormValue("code")), 0700); err != nil {
			fmt.Println(err)
		}

		if err := ioutil.WriteFile("./saved/"+results, []byte(req.FormValue("results")), 0700); err != nil {
			fmt.Println(err)
		}

		// exec.Command("goimports", "-w", filepath).Run()
		// exec.Command("go", "fmt", filepath).Run()
		file, _ := ioutil.ReadFile(filepath)
		res.Write(file)
	}
}

func RunHandler(res http.ResponseWriter, req *http.Request, conn *sql.DB) {
	// if req.Method == "POST" {
	// b, _ := ioutil.ReadFile("static/demo/gosublime.html")

	// t := template.New("")
	// t, _ = t.Parse(string(b))

	req.ParseForm()
	script := req.FormValue("data")

	queryset := strings.Split(script, ";")

	var metastore [][]map[string]interface{}

	for _, query := range queryset {
		rows, err := conn.Query(query)
		if err != nil {
			fmt.Println(err.Error(), "\n")
		}
		defer rows.Close()

		columns, _ := rows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		var store []map[string]interface{}

		for rows.Next() {
			for i, _ := range columns {
				valuePtrs[i] = &values[i]
			}

			rows.Scan(valuePtrs...)

			row := make(map[string]interface{})
			for i, col := range columns {
				var v interface{}
				val := values[i]
				b, ok := val.([]byte)

				if ok {
					v = string(b)
				} else {
					v = val
				}
				row[col] = v
			}

			store = append(store, row)
		}

		if store != nil {
			metastore = append(metastore, store)
		}

	}

	data, _ := json.Marshal(metastore)
	if js, err := json.Marshal(Response{"code": script, "output": string(data)}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.Write(js)
	}

}

func OpenHandler(res http.ResponseWriter, req *http.Request) {
	filepath := req.FormValue("filepath")
	// exec.Command("goimports", "-w", filepath).Run()
	// exec.Command("go", "fmt", filepath).Run()
	if file, err := ioutil.ReadFile("./saved/" + filepath); err != nil {
		fmt.Println(err)
	} else {
		res.Write(file)
	}
}
