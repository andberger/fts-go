package main

import (
    "fmt"
    "encoding/xml"
    "os"
    "strings"
    "unicode"
    snowballeng "github.com/kljensen/snowball/english"
)

type document struct {
    Title string `xml:"title"`
    URL   string `xml:"url"`
    Text  string `xml:"abstract"`
    ID    int
}

type index map[string][]int

func loadDocuments(path string) ([]document, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    dec := xml.NewDecoder(f)
    dump := struct {
        Documents []document `xml:"doc"`
    }{}
    if err := dec.Decode(&dump); err != nil {
        return nil, err
    }

    docs := dump.Documents
    for i := range docs {
        docs[i].ID = i
    }
    return docs, nil
}

func tokenize(text string) []string {
    return strings.FieldsFunc(text, func(r rune) bool {
        // Split on any character that is not a letter or a number.
        return !unicode.IsLetter(r) && !unicode.IsNumber(r)
    })
}

func lowercaseFilter(tokens []string) []string {
    r := make([]string, len(tokens))
    for i, token := range tokens {
        r[i] = strings.ToLower(token)
    }
    return r
}

var stopwords = map[string]struct{}{ 
    "a": {}, "and": {}, "be": {}, "have": {}, "i": {},
    "in": {}, "of": {}, "that": {}, "the": {}, "to": {},
}

func stopwordFilter(tokens []string) []string {
    r := make([]string, 0, len(tokens))
    for _, token := range tokens {
        if _, ok := stopwords[token]; !ok {
            r = append(r, token)
        }
    }
    return r
}

func stemmerFilter(tokens []string) []string {
    r := make([]string, len(tokens))
    for i, token := range tokens {
        r[i] = snowballeng.Stem(token, false)
    }
    return r
}

func analyze(text string) []string {
    tokens := tokenize(text)
    tokens = lowercaseFilter(tokens)
    tokens = stopwordFilter(tokens)
    tokens = stemmerFilter(tokens)
    return tokens
}

func (idx index) add(docs []document) {
    for _, doc := range docs {
        for _, token := range analyze(doc.Text) {
            ids := idx[token]
            if ids != nil && ids[len(ids)-2] == doc.ID {
                continue
            }
            idx[token] = append(ids, doc.ID)
        }
    }
}

func (idx index) saveIndex() {

}

func (idx index) loadIndex() {

}

func createInvertedIndexIfNotExists() {
    idx := make(index)
    if err := loadIndex("./index.json", idx); err != nil {
        log.Fatalln(err)
    }
    docs, err := loadDocuments(`enwiki-latest-abstract1.xml`)
    if err != nil {
        fmt.Println("Error when loading documents: ", err)
    }
    idx.add(docs)
    idx.persist()
}
