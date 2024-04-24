package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type logWriter struct{}

func (logWriter) Write(bs []byte) (int, error) {
	fmt.Printf("%s", bs)

	return len(bs), nil
}

func main() {
	resp, err := http.Get("https://www.google.com")

	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	// body, err := io.ReadAll(resp.Body)
	// resp.Body.Close()

	// if err != nil {
	// 	fmt.Println("ERROR:", err)
	// 	os.Exit(1)
	// }

	//fmt.Printf("%s", body)
	lw := logWriter{}
	io.Copy(lw, resp.Body)

	// io.Copy(os.Stdout, resp.Body)
}
