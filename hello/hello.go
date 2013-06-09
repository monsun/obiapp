package hello

import (
    "appengine"
    "appengine/user"
    "fmt"
    "net/http"
    "html/template"
)

func init(){
    http.HandleFunc("/", root)
    http.HandleFunc("/sign", sign)
}

func root(w http.ResponseWriter, r *http.Request){
    return fmt.Frpint(w, guestbookForm)
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
    err := signTemplate.Execute(w, r.FormValue("content"))
    if err != nil{
        http.Error(w,err.Error(), http.StatusInternalServerError)
    }
}

var signTemplate = template.Must(template.New("sign").Parse(signTemplateHTML))

const signTemplateHTML = `
<html>
    <body>
        <p>You wrote: </p>
        <pre>{{.}}</pre>
    </body>
</html>
`

func handler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
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
}


