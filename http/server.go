package http

import (
    "github.com/Skyhark-Projects/golang-bic-from-iban/log"
    "encoding/json"
    "net/http"
    "strings"
    "errors"
    "io"
)

var routes = []Route{}

type Route struct {
    Path []string
    Callback func(req Request) (interface{}, error)
    Method string
}

type Request struct {
    Url string
    Variables map[string]string
    Post map[string]interface{}
    ClientIP string
}

func (r *Request) Trusted() bool {
    return r.ClientIP == "127.0.0.1" || r.ClientIP[:3] == "10."
}

type Server struct {

}

func Start(port string) {
    server := &Server{}

    log.Info("Started http server", "port", port)
    if err := http.ListenAndServe(":" + port, server); err != nil {
        log.Error("Error starting http server", "error", err)
    }
}

func AddRoute(path string, handler func(req Request) (interface{}, error)) {
    routes = append(routes, Route{
        Path: strings.Split(path[1:], "/"),
        Callback: handler,
        Method: "get",
    })
}

type errResult struct {
    Error string `json:"error"`
}

func (s *Server) handle(w http.ResponseWriter, req *http.Request) (interface{}, error) {
    uri := strings.Split(req.RequestURI[1:], "/")
    if len(uri) > 1 && uri[0] == "api" {
        uri = uri[1:]
    }

    for _, route := range routes {
        if len(route.Path) != len(uri) {
            continue
        }

        variables := map[string]string{}
        same := true
        for index, r := range route.Path {
            if len(r) > 1 && r[:1] == ":" {
                variables[r[1:]] = uri[index]
            } else if r != uri[index] {
                same = false
                break
            }
        }

        if !same {
            continue
        }

        //get ip
        clientIP := ""
        if req.Header.Get("cf-connecting-ip") != "" {
            clientIP = req.Header.Get("cf-connecting-ip")
        } else if req.Header.Get("x-real-ip") != "" {
            clientIP = req.Header.Get("x-real-ip")
        } else if req.Header.Get("x-forwarded-for") != "" {
            clientIP = req.Header.Get("x-forwarded-for")
        } else {
            clientIP = strings.Split(req.RemoteAddr, ":")[0]
        }

        //-----------

        post := map[string]interface{}{}
        if req.Method == "POST" {
            data := json.NewDecoder(req.Body)
            err := data.Decode(&post)
            if err != nil && err != io.EOF {
                return nil, err
            }
        }

        sreq := Request{
            Url: "/" + strings.Join(uri, "/"),
            Variables: variables,
            Post: post,
            ClientIP: clientIP,
        }

        if len(route.Path) > 0 && len(route.Path[0]) > 0 && route.Path[0][:1] == "$" && !sreq.Trusted() {
            continue
        }

        return route.Callback(sreq)
    }
    
    return nil, errors.New("The requested api does not exists")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Cache-Control", "must-revalidate")

    res, err := s.handle(w, req)
    if err != nil {
        bs, err2 := json.Marshal(&errResult{ Error: err.Error() })
        if err2 != nil {
            log.Error("Could not marshal json error", "err", err2)
            http.Error(w, err.Error(), 400)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write(bs)
        }
    } else {
        bs, err := json.Marshal(res)
        if err != nil {
            log.Error("Could not marshal json result", "err", err)
            http.Error(w, err.Error(), 400)
        } else {
            w.Write(bs)
        }
    }
}

func init() {
    AddRoute("/", func(req Request) (interface{}, error) {
        res := []string{}
        for _, route := range routes {
            if len(route.Path) > 0 && len(route.Path[0]) > 1 && route.Path[0][:1] == "$" && !req.Trusted() {
                continue
            }

            res = append(res, "/" + strings.Join(route.Path, "/"))
        }

        return res, nil
    })
}