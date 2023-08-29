package http

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Listener struct {
	server *http.Server
}

type analyzeWork interface {
	Analyze(string)
}

type middlewaredata struct {
	host        string
	port        string
	infologhost *log.Logger
	workcfg     analyzeWork
}

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func New(listenPort string) (Listener, error) {

	fmt.Println("Use http storage...")

	server := &http.Server{
		Addr:    ":" + listenPort,
		Handler: nil,
	}

	return Listener{server}, nil

}

func (l *Listener) Do(host, port string, workcfg interface{}, infologlocal, infologhost *log.Logger) {

	l.server.Handler = middlewareLog(http.HandlerFunc(all), infologlocal, &middlewaredata{
		host:        host,
		port:        port,
		infologhost: infologhost,
		workcfg:     workcfg.(analyzeWork),
	})
	_ = l.server.ListenAndServe()

}

func middlewareLog(next http.Handler, logdebug *log.Logger, data *middlewaredata) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		if logdebug != nil {
			body, _ := ioutil.ReadAll(r.Body)
			logdebug.Println(string(body))
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(context.Background(), "data", data)))
	}

	return http.HandlerFunc(fn)
}

func all(wr http.ResponseWriter, req *http.Request) {

	data := req.Context().Value("data").(*middlewaredata)
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		log.Fatal(err)
		return
	}
	bodyReader := bytes.NewReader(reqBody)

	requestURL := fmt.Sprintf("http://%s:%s"+req.URL.Path, data.host, data.port)
	reqhost, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		log.Fatal(err)
		return
	}

	copyHeader(reqhost.Header, req.Header)

	res, err := http.DefaultClient.Do(reqhost)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		log.Fatal(err)
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		log.Fatal(err)
		return
	}

	if data.infologhost != nil {
		data.infologhost.Println(string(resBody))
	}

	data.workcfg.Analyze(string(reqBody))

	delHopHeaders(res.Header)
	copyHeader(wr.Header(), res.Header)
	wr.WriteHeader(res.StatusCode)
	wr.Write(resBody)

}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
