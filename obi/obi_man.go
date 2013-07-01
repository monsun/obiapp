package obi

import (
    "appengine"
    "appengine/datastore"
    "fmt"
    "net/http"
    "html/template"
)

type Excercise struct {
    Name string
}

func init(){
    http.HandleFunc("/obi/add_excercise", AddExcercise)
    http.HandleFunc("/obi/list_excercises", ListExcercises)
    http.HandleFunc("/obi/edit_excercise", EditExcercise)
    http.HandleFunc("/obi/delete_excercise", DeleteExcercise)
}

func AddExcercise(w http.ResponseWriter, r *http.Request){
    c := appengine.NewContext(r)

    exc := Excercise {
        Name: r.FormValue("ExcerciseName"),
    }

    key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Excercise", nil), &exc)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    c.Debugf("Saved", key)

    http.Redirect(w, r, "/obi/list_excercises", http.StatusFound)
}

func ListExcercises(w http.ResponseWriter, r *http.Request){
    c := appengine.NewContext(r)
    q := datastore.NewQuery("Excercise").Limit(10)
    excercises := make([]Excercise, 0, 10)
    if _, err := q.GetAll(c, &excercises); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    //XXX: debug loop
    for excercise := range excercises {
        c.Debugf("Loaded ",excercise)
    }

    if err := excerciseListTemplate.Execute(w, excercises); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var excerciseListTemplate = template.Must(template.New("ExcerciseList").Parse(excerciseListTemplateHTML))

const excerciseListTemplateHTML = `
<html>
    <body>
        <h3>View excercises</h3>
        {{range .}}
            <form method="post" action="delete_excercise">
                <input type="hidden" name="KeyToDelete" value="{{.Name}}"/>
                <pre>{{.Name}}</pre>
                <input type="submit" value="Delete!"/>
            </form>
        {{end}} 

        <h3>Or add a new one</h3>
        <form action="add_excercise" method="post">
            <div><input type="text" name="ExcerciseName"/></div>
            <div><input type="submit" value="Save excercise"/></div>
        </form> 
    </body>
</html>
`
func DeleteExcercise(w http.ResponseWriter, r *http.Request){

    c := appengine.NewContext(r)
    keyVal := r.FormValue("KeyToDelete")

    if keyVal == "" {
        http.Error(w, "Nil key to delete", http.StatusInternalServerError)
        return
    }

    key, err := datastore.DecodeKey(keyVal)

    if err != nil {
        c.Errorf(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
//        return
    }

    if err := datastore.Delete(c, key); err != nil {
        c.Infof("Key ",key)
        c.Errorf(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
//        return
    }

    http.Redirect(w, r, "/obi/list_excercises", http.StatusFound)
}

func EditExcercise(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "AddExc here!")
}
