package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    "github.com/nut-abctech/handy-scrape/libs/parser"
    "gopkg.in/xmlpath.v1"
)

type record struct {
    Name, Location, Contact, Detail string
}

var chResult chan []*record
var chDone chan bool

const MAX_PAGES int = 58

func init() {
    chDone = make(chan bool)
    chResult = make(chan []*record, 10)
}

func main() {
    var i int = 1
    err := os.Mkdir("./dataset", 0777)
    if err != nil {
        if !os.IsExist(err) {
            log.Panicln(err)
        }
    }
    if role, err := parser.Parse("./path.json"); err == nil {
        // Channels
        hitURL := role.URL + role.Get + strconv.Itoa(i)
        for ; i <= MAX_PAGES; i++ {
            go start(hitURL, role)
            hitURL = role.URL + role.Get + strconv.Itoa(i)
        }
        var c int = 0
        var aggregateData = make([]*record, 0)
        for {
            select {
            case results := <-chResult:
                aggregateData = append(aggregateData, results...)
            case <-chDone:
                c++
            }
            if c == MAX_PAGES {
                break
            }
        }

        writeFile("./dataset/"+time.Now().Format(time.Stamp)+".json", aggregateData)
        log.Printf("DONE scrape %d pages", c)
    } else {
        log.Panicln(err)
    }
}

func writeFile(filename string, content []*record) {
    bytes, err := json.MarshalIndent(content, "", "    ")
    if err != nil {
        log.Println("Unable to marshal json data ", err)
        return
    }
    if f, err := os.Create(filename); err == nil {
        defer f.Close()
        if size, err := f.Write(bytes); err == nil {
            log.Printf("Created %s size: %d bytes", filename, size)
        } else {
            log.Printf("Fail to write the result file %s Reason: %s", filename, err)
            return
        }
    } else {
        log.Println("Unable to create a file ", err)
        return
    }
}

func start(hitURL string, role *parser.Role) {
    defer func() {
        // Notify that we're done after this function
        chDone <- true
    }()
    if resp, err := http.Get(hitURL); err == nil {
        // log.Println("GET:", hitURL)
        body := resp.Body
        defer body.Close()
        entry := xmlpath.MustCompile(role.Routes.Entry)
        if root, err := xmlpath.ParseHTML(body); err == nil {
            results := crawl(entry.Iter(root), role)
            chResult <- results
        } else {
            log.Println("Unable to parse HTML :: ", err)
            chResult <- nil
        }
    } else {
        log.Println("Unable to hit request :: ", err)
        chResult <- nil
    }
}

func crawl(iterContent *xmlpath.Iter, role *parser.Role) []*record {
    info := role.Routes.Info
    // attributes
    // linkPath := xmlpath.MustCompile(info.Link)
    namePath := xmlpath.MustCompile(info.Name)
    detailPath := xmlpath.MustCompile(info.Detail)
    locationPath := xmlpath.MustCompile(info.Location)
    contactNoPath := xmlpath.MustCompile(info.ContactNo)
    var results = make([]*record, 0)

    for iterContent.Next() {
        node := iterContent.Node()
        name, _ := namePath.String(node)
        detail, _ := detailPath.String(node)
        location, _ := locationPath.String(node)
        contactNo, _ := contactNoPath.String(node)
        results = append(results, &record{name, location, contactNo, detail})
    }
    return results
}
