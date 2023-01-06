package requests

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

var (
    ErrInvalidResponse = fmt.Errorf("invalid http response")
)

type Cmd struct {
    err  error
    resp *http.Response
}

func (c *Cmd) Error() error {
    return c.err
}

func (c *Cmd) StatusCode() int {
    if c.resp == nil {
        return 0
    }
    return c.resp.StatusCode
}
func (c *Cmd) Response() *http.Response {
    return c.resp
}
func (c *Cmd) getData() ([]byte, error) {
    if c.resp == nil {
        return nil, ErrInvalidResponse
    }
    return ioutil.ReadAll(c.resp.Body)
}
func (c *Cmd) BindJson(v interface{}) error {
    data, err := c.getData()
    if err != nil {
        return err
    }
    return json.Unmarshal(data, v)
}
func (c *Cmd) BindText(v *string) error {
    data, err := c.getData()
    if err != nil {
        return err
    }
    *v = string(data)
    return nil
}
