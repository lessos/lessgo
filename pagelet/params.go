package pagelet

import (
    "net/url"
)

type Params struct {
    url.Values // A unified view of all the individual param maps below

    // Set by the ParamsFilter
    Query url.Values // Parameters from the query string, e.g. /index?limit=10
    Form  url.Values // Parameters from the request body.
}

func ParamsFilter(c *Controller) {

    c.Params.Values = make(url.Values, 0)

    c.Params.Query = c.Request.URL.Query()
    for k, v := range c.Params.Query {
        if _, ok := c.Params.Values[k]; !ok {
            c.Params.Values[k] = v
        }
    }

    c.Params.Form = c.Request.Form
    for k, v := range c.Params.Form {
        if _, ok := c.Params.Values[k]; !ok {
            c.Params.Values[k] = v
        }
    }
}
