package hello

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
    "html/template"
    "fmt"
    "net/http"
    "time"
)

type Greeting struct {
    Author string
    Content string
    Date time.Time
}

func init(){
    http.HandleFunc("/hello", handle)
    http.HandleFunc("/sign", sign)
}

func root(w http.ResponseWriter, r *http.Request){
    fmt.Fprint(w, guestbookForm)
}

const guestbookForm = `
<html>
    <body>
        <form action="/sign" method="post">
            <div><textarea name="content" rows="3" cols="60"></textarea></div>
            <div><input type="submit" value="Sign guestbook"/></div>
        </form>
    </body>
</html>
`

func sign(w http.ResponseWriter, r *http.Request){
    c := appengine.NewContext(r)
    g := Greeting {
        Content: r.FormValue("content"),
        Date: time.Now(),
    }
    if u:= user.Current(c); u != nil {
        g.Author = u.String()
    }

    _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Greeting", nil),  &g)
    if err != nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/hello", http.StatusFound)
}

func handle(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    q := datastore.NewQuery("Greeting").Order("-Date").Limit(10)
    greetings := make([]Greeting, 0, 10)
    if _, err := q.GetAll(c, &greetings); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err := guestbookTemplate.Execute(w, greetings); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
/*
    u := user.Current(c)
    if u == nil {
        url, err := user.LoginURL(c, r.URL.String())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Location", url)
        w.WriteHeader(http.StatusFound)
        return
    }
    fmt.Fprint(w, "Hello, world")
*/
}

var guestbookTemplate = template.Must(template.New("book").Parse(guestbookTemplateHTML))

const guestbookTemplateHTML = `
<html>
    <body>
        {{range .}}
            {{with .Author}}
                <p><b>{{.}}</b> wrote:</p>
            {{else}}
                <p>An anonymous person wrote:</p>
            {{end}}
            <pre>{{.Content}}</pre>
        {{end}}
        <form action="/sign" method="post">
            <div><textarea name="content" rows="3" cols="60"></textarea></div>
            <div><input type="submit" value="Sign Guestbook"/></div>
        </form>        
    </body>
</html>
`
