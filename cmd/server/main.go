package main

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"sync"
)

type Counter struct {
	value int
	mu    sync.Mutex
}

func (c *Counter) Increase() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *Counter) GetValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	slog.Info("Plant Journal")
	counter := &Counter{}

	bh := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		data := map[string]int{
			"CounterValue": counter.GetValue(),
		}
		tmpl.Execute(w, data)
	}

	ih := func(w http.ResponseWriter, _ *http.Request) {
		tmplStr := "<div id=\"counter\">{{.CounterValue}}</div>"
		tmpl := template.Must(template.New("counter").Parse(tmplStr))
		counter.Increase()
		data := map[string]int{
			"CounterValue": counter.GetValue(),
		}
		tmpl.ExecuteTemplate(w, "counter", data)
	}

	// define handlers
	http.HandleFunc("/", bh)
	http.HandleFunc("/increase", ih)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
