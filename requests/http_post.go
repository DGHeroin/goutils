package requests

import (
    "io"
    "net/http"
)

var (
    GetClient = func() *http.Client {
        return http.DefaultClient
    }
)

func PostJson(url string, r io.Reader) *Cmd {
    result := &Cmd{}
    client := GetClient()
    req, err := http.NewRequest(http.MethodPost, url, r)
    if err != nil {
        result.err = err
        return result
    }
    req.Header.Set("Content-Type", "application/json; charset=UTF-8")

    resp, err := client.Do(req)
    if err != nil {
        result.err = err
        return result
    }
    result.resp = resp
    return result
}
func PostJsonBind(url string, r io.Reader, v interface{}) error {
    cmd := PostJson(url, r)
    if cmd.Error() != nil {
        return cmd.Error()
    }
    return cmd.BindJson(v)
}
