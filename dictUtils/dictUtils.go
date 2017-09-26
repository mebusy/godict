package dictUtils

import ( 
    "github.com/mebusy/godict/sqlManager"
)


func DBConnect( db_path string ) int {
    sql := sqlManager.GetInstance()
    idx := sql.OpenDB( db_path )
    return idx 
}

func DBCloseAll() {
    sql := sqlManager.GetInstance()
    sql.CloseAll()    
}
