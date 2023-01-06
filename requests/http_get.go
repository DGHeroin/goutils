package requests

import "net/http"

func GetJson(url string) *Cmd {
    result := &Cmd{}
    client := GetClient()
    req, err := http.NewRequest(http.MethodGet, url, nil)
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
func GetJsonBind(url string, v interface{}) error {
    cmd := GetJson(url)
    if cmd.Error() != nil {
        return cmd.Error()
    }
    return cmd.BindJson(v)
}
