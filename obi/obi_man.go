package obi

import (
    "appengine"
    "appengine/datastore"
//    "fmt"
    "net/http"
    "html/template"
)

type Excercise struct {
    Name string
    DSKeyEncoded string
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

	key := *datastore.NewKey(c, "Excercise", exc.Name, 0, nil)
	exc.DSKeyEncoded = key.Encode()
    _, err := datastore.Put(c, &key, &exc)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    c.Debugf("Saved", key, exc.DSKeyEncoded, exc.Name)

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
func DeleteExcercise(w http.ResponseWriter, r *http.Request){

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
        c.Infof("Key ",keyVal)
        c.Errorf(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/obi/list_excercises", http.StatusFound)
}

func EditExcercise(w http.ResponseWriter, r *http.Request){
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
    
    if err!=nil{
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    }
 	   
 	if r.Method == "POST" {
 		Exc.Name = r.FormValue("Name")
 	 	
 		if _, err := datastore.Put(c, keyVal, &Exc); err!=nil {
 			http.Error(w, err.Error(), http.StatusInternalServerError)
 		}
 		
 		http.Redirect(w, r, "/obi/list_excercises", http.StatusFound)
 	} else if r.Method == "GET" {
 		if err := excerciseEditTemplate.Execute(w, Exc); err != nil {
        	http.Error(w, err.Error(), http.StatusInternalServerError)
    	}
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
            <input type="text" value="{{.Name}}"/>
            <input type="submit" value="Save" style="display: inline"/>
        </form>
        
        <p/>
        <input type="button" value="Cancel"/> 
    </body>
</html>
`
