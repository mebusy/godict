package dictUtils

import ( 
    "dict/sqlManager"
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
