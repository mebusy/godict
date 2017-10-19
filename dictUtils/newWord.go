package dictUtils

import ( 
    "fmt"
    // "strconv"
    "github.com/mebusy/godict/sqlManager"
    // "github.com/mebusy/godict/cryptor"
    . "github.com/mebusy/godict/extension"
    "log"
    // "strings"
    // "bytes"
    // "encoding/json"
    "database/sql"
)

var db_conn_newword * sql.DB = nil
var table_name_newword = "newword"
func InitNewWordDB( db_path string  ) {
    log.Println( "initialing newword DB... " )
    db_idx := DBConnect( db_path )
    db := sqlManager.GetInstance().GetDBConnectionByIndex( db_idx )
    if db == nil {
        return
    }

    stmt, err := db.Prepare( "create table if not exists " + table_name_newword  + " ( word text primary key , desc text NOT NULL )" )

    if HasErr(err) { return }

    _ , err = stmt.Exec( )
    if HasErr(err) { return }

    db_conn_newword = db
    log.Println( "newword DB initialiezd" )
}

func AddNewWord( word string, desc string  ) {
    if db_conn_newword == nil { 
        log.Println( "can not find connection to newword DB" )
        return 
    }

    db := db_conn_newword 

    stmt, err := db.Prepare( "INSERT or replace INTO " + table_name_newword + " values (?,?)"  )
    if HasErr(err) { return }

    _ , err = stmt.Exec( word, desc  )
    if HasErr(err) { return }

    log.Println( fmt.Sprintf( "new word %s added." , word  ) )
    
}

func RemoveNewWord( word string ) {
    if db_conn_newword == nil { 
        log.Println( "can not find connection to newword DB" )
        return 
    }

    db := db_conn_newword 
    stmt, err := db.Prepare( "DELETE FROM " + table_name_newword + " where word=?") 
    if HasErr(err) { return }          
    _ , err = stmt.Exec( word  )
    if HasErr(err) { return }

    log.Println( fmt.Sprintf( "new word %s deleted." , word  ) )  

}

