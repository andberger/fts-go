package main

import (
    "os"
    "strings"
)

func search(docs []document, term string) []document {
    var r []document
    for _, doc := range docs {
        if strings.Contains(doc.Text, term) {
            r = append(r, doc)
        }
    }
    return r
}

func main() {
    args := os.Args[1:]
    searchTerm := args[0]

    createInvertedIndexIfNotExists()
}
