package main

import ( 
    "fmt"
    "github.com/mebusy/godict/dictUtils"
)

func main() {
    path := "../../../mebusy_git_dict/dict_client/Assets/StreamingAssets/dict.db"

    idx1 := dictUtils.DBConnect(path)
    idx2 := dictUtils.DBConnect(path)
    idx3 := dictUtils.DBConnect(path)
    fmt.Println( idx1,idx2,idx3 )



    dictUtils.DBCloseAll()
    fmt.Println( "done" )
}
