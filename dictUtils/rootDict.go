package dictUtils

import ( 
    "fmt"
    "strconv"
    "github.com/mebusy/godict/sqlManager"
    // "github.com/mebusy/godict/cryptor"
    . "github.com/mebusy/godict/extension"
    "log"
    "strings"
    "bytes"

)

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
    cmdText :=  "SELECT word,en1,cn1,show_en1,en2,cn2,show_en2,en3,cn3,show_en3,en4,cn4,show_en4 from root ORDER BY word COLLATE NOCASE " 
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

    log.Println( "root dict loaded ", len( _rootDict ) , "total"  )
}

// get format root meaning , it should has en meaning, cn meaning 
func generateFormatedMeaning(root string , ch chan string) {
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
    ch <- strings.TrimSpace( sb.String() )
}


// get formated roots has the syn means with specified `root`
func getSynonymsRoots( db_idx int,  root string , ch chan string )  {
    sql := sqlManager.GetInstance()
    if db_idx<0 || db_idx > sql.ConnectionCount() {
        log.Println( "GetSynonymsRoots","wrong db index" )
        ch <- ""
    }

    synRoots:= make( []string , 0, 8 ) 
    means := _rootDict[root]
    
    cmdText := ""   
    for i:=0 ; i< len(means) ; i+=3 {
        en := means [i]
        cmdText += fmt.Sprintf(  " %[2]v en1 = \"%[1]v\" OR en2 = \"%[1]v\" OR en3 = \"%[1]v\"  OR en4 = \"%[1]v\"  " ,
            en , If(i==0,"","OR").(string) )    
    }
    cmdText = fmt.Sprintf( "SELECT word FROM root WHERE %[1]v ORDER BY word COLLATE NOCASE " , cmdText  )

    db := sql.AllConns[ db_idx ]
    rows, err := db.Query( cmdText ) 
    if HasErr(err) {
        ch <- ""
    }
    defer rows.Close()
    
    for rows.Next() {
        var word string 
        err := rows.Scan( &word )
        if HasErr(err) {
            break
        }
        if word != root {
            synRoots = append(synRoots,word)   
        }
    }

    var seeAlso string 
    for i, key := range synRoots {
        seeAlso += fmt.Sprintf( "%[3]v<b><color=%[2]v>%[1]v</color></b>", key, COLOR_ROOT , If(i==0,"",",").(string)  )   
    }
    if seeAlso != "" {
        seeAlso = "\nSee also : " + seeAlso + "\n"   
    }

    ch <- strings.TrimSpace(seeAlso) 
    
}

func generateRootWordExamples(db_idx int,  root string, ch chan string ) {
    sql := sqlManager.GetInstance()
    if db_idx<0 || db_idx > sql.ConnectionCount() {
        log.Println( "GenerateRootWordExamples","wrong db index" )
        ch <- ""
    }

    var sb bytes.Buffer 

    cmdText := fmt.Sprintf( "SELECT word , ex FROM dict WHERE  ex != \"\" AND (root1 = \"%[1]v\"  OR  root2 = \"%[1]v\" ) ORDER BY word COLLATE NOCASE " , root  )

    db := sql.AllConns[ db_idx ]
    rows, err := db.Query( cmdText ) 
    if HasErr(err) {
        ch <- ""
    }
    defer rows.Close()
    
    for rows.Next() {
        var word string 
        var ex string 
        err := rows.Scan( &word , &ex )
        if HasErr(err) {
            break
        }
        sb.WriteString( fmt.Sprintf(  "<color=blue>%v</color>:\u3000%v\n\n" , word, ex   ) )
    }

    ch <- strings.TrimSpace(sb.String() )
}



