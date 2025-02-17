package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type CacheObj struct {
	Response     *http.Response
	ResponseBody []byte
}

type Proxy struct {
	Domain        string
	ForwardDomain string
	Cache         map[string]*CacheObj
	Mutex         sync.RWMutex
}

func NewProxyServer(domain, forwardDomain string) *Proxy {
	return &Proxy{
		Domain:        domain,
		ForwardDomain: forwardDomain,
		Cache:         map[string]*CacheObj{},
	}
}

func (p *Proxy) HandleFunc(w http.ResponseWriter, r *http.Request) {
	requestURL := p.ForwardDomain + r.URL.Path

	fmt.Println("REQUEST URL:", requestURL)

	p.Mutex.RLock()
	if cache, exist := p.Cache[requestURL]; exist {
		RespondWithHeaders(w, *cache.Response, cache.ResponseBody, "HIT")
		return
	}
	p.Mutex.RUnlock()

	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error Forwarding Request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Error Forwarding Body Request", http.StatusInternalServerError)
		return
	}
	fmt.Println(res.StatusCode)
	go func() {
		p.Mutex.RLock()
		p.Cache[requestURL] = &CacheObj{
			Response:     res,
			ResponseBody: body,
		}
		p.Mutex.RUnlock()
	}()
	RespondWithHeaders(w, *res, body, "MISS")
}

func (p *Proxy) ClearCache() {
	p.Mutex.RLock()
	p.Cache = map[string]*CacheObj{}
	p.Mutex.RUnlock()
}

func RespondWithHeaders(w http.ResponseWriter, response http.Response, body []byte, cacheHeader string) {
	fmt.Printf("Cache : %s\n", cacheHeader)
	w.Header().Set("X-Cache", cacheHeader)
	w.WriteHeader(response.StatusCode)
	for k, v := range response.Header {
		w.Header()[k] = v
	}
	w.Write(body)
}
