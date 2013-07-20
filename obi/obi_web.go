package obi

import (
    "appengine"
    "appengine/datastore"
    //    "fmt"
    "html/template"
    "net/http"
)

const ExcercisesLimit = 10

func init() {
    http.HandleFunc("/obi/add_excercise", AddExcerciseAction)
    http.HandleFunc("/obi/list_excercises", ListExcercisesAction)
    http.HandleFunc("/obi/edit_excercise", EditExcerciseAction)
    http.HandleFunc("/obi/delete_excercise", DeleteExcerciseAction)
}

func AddExcerciseAction(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    err := AddExcercise(r.FormValue("ExcerciseName"), c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/obi/list_excercises", http.StatusFound)
}

func ListExcercisesAction(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    excercises, err := ListExcercises(ExcercisesLimit, c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    if err := excerciseListTemplate.Execute(w, excercises); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var excerciseListTemplate = template.Must(template.New("ExcerciseList").Parse(excerciseListTemplateHTML))

const excerciseListTemplateHTML = `
<html>
    <body>
		<script src="//ajax.googleapis.com/ajax/libs/jquery/2.0.3/jquery.min.js"></script>
        <h3>View excercises</h3>
        {{range .}}
            <form method="post" action="delete_excercise">
                <input type="hidden" name="KeyToDelete" value="{{.DSKeyEncoded}}"/>
                <pre>{{.Name}}</pre>
                <input type="submit" value="Delete!" style="display: inline"/>
		<input type="button" value="Edit" style="display: inline" onclick="document.location.href='edit_excercise?KeyToEdit={{.DSKeyEncoded}}'"/>
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

func DeleteExcerciseAction(w http.ResponseWriter, r *http.Request) {

    c := appengine.NewContext(r)
    keyVal, err := datastore.DecodeKey(r.FormValue("KeyToDelete"))

    if keyVal == nil {
        http.Error(w, "Nil key to delete", http.StatusInternalServerError)
        return
    }

    if err != nil {
        c.Errorf(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := datastore.Delete(c, keyVal); err != nil {
        c.Infof("Key ", keyVal)
        c.Errorf(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/obi/list_excercises", http.StatusFound)
}

func EditExcerciseAction(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    keyVal, err := datastore.DecodeKey(r.FormValue("KeyToEdit"))

    if err != nil {
        c.Errorf(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if keyVal == nil {
        http.Error(w, "Nil key to edit", http.StatusInternalServerError)
        return
    }

    c.Debugf("Key to load", keyVal)

    var Exc Excercise
    err = datastore.Get(c, keyVal, &Exc)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    if r.Method == "POST" {
        Exc.Name = r.FormValue("Name")

        if _, err := datastore.Put(c, keyVal, &Exc); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }

        http.Redirect(w, r, "/obi/list_excercises", http.StatusFound)
    } else if r.Method == "GET" {
        if err := excerciseEditTemplate.Execute(w, Exc); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    } else {
        http.Error(w, "Unsupported method: "+r.Method, http.StatusInternalServerError)
    }

}

var excerciseEditTemplate = template.Must(template.New("ExcerciseEdit").Parse(excerciseEditTemplateHTML))

const excerciseEditTemplateHTML = `
<html>
    <body>
		<script src="//ajax.googleapis.com/ajax/libs/jquery/2.0.3/jquery.min.js"></script>
        <h3>Edit excercise</h3>
        <form method="post" action="edit_excercise">
		<input type="hidden" name="KeyToEdit" value="{{.DSKeyEncoded}}"/>
            <input type="text" value="{{.Name}}" name="Name"/>
            <input type="submit" value="Save" style="display: inline"/>
        </form>
        
        <p/>
        <input type="button" value="Cancel" onclick="document.location.href='list_excercises'"/> 
    </body>
</html>
`
