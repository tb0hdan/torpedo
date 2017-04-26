package main

import (
        "io/ioutil"
        "net/http"
       )

func GetURLBytes(url string) (result []byte, err error) {
    response, err := http.DefaultClient.Get(url)
    if err != nil {
        return
    }
    defer response.Body.Close()
    result, err = ioutil.ReadAll(response.Body)
    return
}
