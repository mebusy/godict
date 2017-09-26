package sqlManager

import (
    "path"
    "log"
    "sync"
    "database/sql"
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

func GetInstance() sqlManager {
    once.Do(func() {
        instance = sqlManager{ 
            make( map[string] int ) , 
            make( []* sql.DB ,0,2 )  }
    })

    return instance
}

func ( m sqlManager ) ConnectionCount() int {
    return len( m.AllConns )
}

func checkErr(err error) {
    if err != nil {
        log.Println( err.Error() )
        // panic(err)
    }
}

//==========================================

func (m sqlManager) OpenDB( db_path string ) int  {
    _, f := path.Split( db_path )
    if v, ok := m.ConnIdx[ f ]; ok {
        return v
    }
    db , err := sql.Open("sqlite3", db_path )
    if err != nil {
        checkErr(err)
        return -1 
    }
    log.Printf( "connect to %s\n", f )
    m.ConnIdx[ f ] = len( m.AllConns  )
    m.AllConns = append( m.AllConns , db )

    return m.ConnIdx[ f ]
}

func (m sqlManager) CloseAll() {
    for _ , db := range m.AllConns  {
        db.Close()
    }
}

        // // query
        // rows, err := db.Query("SELECT * FROM userinfo")
        // checkErr(err)
        // var uid int
        // var username string
        // var department string
        // var created time.Time
        //
        // for rows.Next() {
        //     err = rows.Scan(&uid, &username, &department, &created)
        //     checkErr(err)
        //     fmt.Println(uid)
        //     fmt.Println(username)
        //     fmt.Println(department)
        //     fmt.Println(created)
        // }
        //
        // rows.Close() //good habit to close
        //



