package gui

import (
	"io"
	"net/http"
)

func runServer(onShowRequested func()) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		onShowRequested()
		io.WriteString(w, "timetrace says hello\n")
	})

	err := http.ListenAndServe(":15432", nil)
	if err != nil {
		panic(err)
	}
}

func showActiveTimetraceGUI() bool {
	_, err := http.Get("http://localhost:15432")
	return err == nil
}
