package sqlManager

import (
    // "fmt"
    "path"
    "log"
    "sync"
    "database/sql"
    . "github.com/mebusy/godict/extension"
    _ "github.com/mattn/go-sqlite3"
)

type sqlManager struct {
    ConnIdx  map[string] int 
    AllConns []* sql.DB 
}
var (
    once sync.Once

    instance sqlManager
)

func GetInstance() *sqlManager {
    once.Do(func() {
        log.Println( "instantialize sqlManager" )
        instance = sqlManager{ 
            make( map[string] int ) , 
            make( []* sql.DB ,0,2 )  }
    })

    return &instance
}



//==========================================

func ( m sqlManager ) ConnectionCount() int {
    return len( m.AllConns )
}

func (m *sqlManager) OpenDB( db_path string ) int  {
    _, f := path.Split( db_path )
    if v, ok := m.ConnIdx[ f ]; ok {
        return v
    }
    db , err := sql.Open("sqlite3", db_path )
    if HasErr( err ) {
        return -1 
    }
    m.ConnIdx[f] = len( m.AllConns  )
    m.AllConns = append( m.AllConns , db )

    log.Printf( "connect to %s, conn:%d\n", f ,m.ConnIdx[f] )

    return m.ConnIdx[f]
}

func (m sqlManager) CloseAll() {
    for i , db := range m.AllConns  {
        db.Close()
        log.Printf( "conn %d closed\n" ,i )
    }
}



