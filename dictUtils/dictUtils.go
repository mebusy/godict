package dictUtils

import ( 
    "fmt"
    "strconv"
    "github.com/mebusy/godict/sqlManager"
    . "github.com/mebusy/godict/extension"
    "log"
    "strings"
    "bytes"

)


func DBConnect( db_path string ) int {
    sql := sqlManager.GetInstance()
    idx := sql.OpenDB( db_path )
    // fmt.Printf( "opened, total conns:%d,%p\n", sql.ConnectionCount() , sql )
    return idx 
}

func DBCloseAll() {
    sql := sqlManager.GetInstance()
    sql.CloseAll()    
}

var _rootDict  map[string] []string
var __allkeys = []string{}


func LoadRootDict( db_idx int ) {
    if _rootDict != nil {
        return
    }
    sql := sqlManager.GetInstance()
    if db_idx<0 || db_idx > sql.ConnectionCount() {
        log.Println( "LoadRootDict","wrong db index" )
        return 
    }
    // load dict
    cmdText :=  "SELECT * from root ORDER BY word COLLATE NOCASE " 
    db := sql.AllConns[ db_idx ]
    rows, err := db.Query( cmdText ) 
    if HasErr(err) {
        return    
    }
    defer rows.Close()

    _rootDict = make( map[string] []string , 1200 ) 
    for rows.Next() {
        items := make( []string, 0, 6 )
        var word string 
        
        var en [4]string  
        var cn [4]string  
        var show [4]int
        
        err := rows.Scan( &word,
            &en[0],&cn[0],&show[0],
            &en[1],&cn[1],&show[1],
            &en[2],&cn[2],&show[2],
            &en[3],&cn[3],&show[3]) 
        if HasErr(err) {
            break
        }

        for i:=0; i<4 ; i++ {
            if en[i] == "" {
                break
            }
            items = append(items , en[i])
            items = append(items , cn[i] )
            items = append(items ,strconv.Itoa(show[i]) )
        }
        _rootDict [word] = items
        __allkeys = append( __allkeys, word )
    }

    // fmt.Printf( "%v \n", _rootDict  ) 
    
    log.Println( "root dict loaded ", len( _rootDict ) , "total"  )
}

const COLOR_ROOT = "#008b8b" 

func GetSynonymsRoots(root string) string {
    var sb bytes.Buffer
    sb.WriteString( fmt.Sprintf( "<color=blue>词根</color>:  <b><color=%v>%v</color></b>\n" , COLOR_ROOT, root   ) )

    means := _rootDict[root]
    for i:= 0 ; i< len(means) ; i+=3 {
        en := fmt.Sprintf ( "=<color=%v>%v</color>,", COLOR_ROOT , means [i+0]  )    
        sepCnMeans := strings.Split( means [i+1],"--" ) 
        mean := sepCnMeans [0]
        extendedMeaning := If(strings.Contains( means [i+1],"--" ) , ( ","+sepCnMeans [ len(sepCnMeans)-1 ] ) , "").(string)
        bShowEn := means [i+2] == "1"
        cn := fmt.Sprintf ( "表示\"%v\"%v" , mean , extendedMeaning   )

        meanIdx  := (i + 3) / 3
        sb.WriteString( fmt.Sprintf( "\u3000%v%v%v\n" ,
            If( len(means) > 3 , fmt.Sprintf("%v%v",meanIdx , ". ") ,"").(string),
            If(bShowEn,en,"").(string),cn ))

    }
    return sb.String()
}







