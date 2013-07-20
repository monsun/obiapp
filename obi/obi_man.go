package obi

import (
    "appengine"
    "appengine/datastore"
    //    "fmt"
)

type Excercise struct {
    Name         string
    DSKeyEncoded string
}

func AddExcercise(ExcName string, c appengine.Context) error {
    exc := Excercise{
        Name: ExcName,
    }

    key := *datastore.NewKey(c, "Excercise", exc.Name, 0, nil)
    exc.DSKeyEncoded = key.Encode()
    _, err := datastore.Put(c, &key, &exc)

    c.Debugf("Saved", key, exc.DSKeyEncoded, exc.Name)
    return err
}

func ListExcercises(limit int, c appengine.Context) ([]Excercise, error) {
    q := datastore.NewQuery("Excercise").Limit(limit)
    excercises := make([]Excercise, 0, limit)
    _, err := q.GetAll(c, &excercises);

    return excercises, err
}

func DeleteExcercise(KeyToDelete string, c appengine.Context) error {
    keyVal, err := datastore.DecodeKey(KeyToDelete)

    if err != nil {
        c.Errorf(err.Error())
        return err
    }

    if err := datastore.Delete(c, keyVal); err != nil {
        c.Errorf(err.Error())
        return err
    }

    return nil
}

func EditExcercise(KeyToEdit string, Exc Excercise, c appengine.Context) error {
    keyVal, err := datastore.DecodeKey(KeyToEdit)

    if err != nil {
        return err
    }

    c.Debugf("Key to load", keyVal)

    //var Exc Excercise
    //err = datastore.Get(c, keyVal, &Exc)
    //if err != nil {
    //    return err
    //}

    if _, err := datastore.Put(c, keyVal, &Exc); err != nil {
        return err
    }

    return nil
}

func LoadExcercise(KeyToLoad string, c appengine.Context, Exc *Excercise) error {

    keyVal, err := datastore.DecodeKey(KeyToLoad)

    if err != nil {
        return err
    }

    err = datastore.Get(c, keyVal, Exc)

    if err != nil {
        return err
    }

    return nil
}
