package main

import  (
        "flag"
        "strings"
        )


var (
    token = flag.String("token", "", "Comma separated list of Slack legacy tokens")
)


func main() {
        var keys []string

        flag.Parse()

        for _, key := range strings.Split(*token, ",") {
            keys = append(keys, key)
        }
        RunBots(keys)
}
