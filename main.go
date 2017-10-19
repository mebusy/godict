package main

import ( 
    "fmt"
    "github.com/mebusy/godict/dictUtils"
    "os"
    "path/filepath"
    "runtime"
    "strings"
)

func main() {
    ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)
    fmt.Println("pwd", exPath)


    path := "/Volumes/WORK/WORK/mebusy_git_dict/dict_client/Assets/StreamingAssets/dict.db"
    path2 := "/Volumes/WORK/WORK/mebusy_git_dict/dict_client/Assets/StreamingAssets/pron.db"
    path3 := "/Volumes/WORK/WORK/mebusy_git_dict/dict_client/newword.db"

    if runtime.GOOS == "linux" {
        fmt.Println( runtime.GOOS )
        path = "/home/qibinyi/WORK/dict_dbs/dict.db"
        path2 = "/home/qibinyi/WORK/dict_dbs/pron.db"
        path3 = "/home/qibinyi/WORK/dict_dbs/newword.db"
    }

    //*/
    db_idx_dict := dictUtils.DBConnect(path)
    idx2 := dictUtils.DBConnect(path2)
    dictUtils.InitNewWordDB( path3 )

    fmt.Println( db_idx_dict ,idx2)
    
    dictUtils.LoadRootDict( db_idx_dict )


    // fmt.Println(    dictUtils.GenerateRootInterpretation  ( db_idx_dict , "arch" ) )

    dictUtils.SearchWordLike( db_idx_dict, "percei" )

    // fmt.Println( dictUtils.GetWordInterpretation(  db_idx_dict , "economy" ) )

    dictUtils.AddNewWord( "testword" , "hahaha" )
    dictUtils.AddNewWord( "testword2" , "hahaha" )
    dictUtils.RemoveNewWord( "testword" )

    dictUtils.DBCloseAll()

    fmt.Println( "test..." )
    r := strings.NewReplacer("(", "", ")", "")
    fmt.Println( r.Replace( "a(b)c)d" )  )

    fmt.Println( "done" )
}
