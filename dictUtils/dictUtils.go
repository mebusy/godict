package dictUtils

import ( 
    "fmt"
    // "strconv"
    "github.com/mebusy/godict/sqlManager"
    "github.com/mebusy/godict/cryptor"
    . "github.com/mebusy/godict/extension"
    "log"
    "strings"
    "bytes"

)

func init() {
    log.SetPrefix( "godict " )    
}


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

const COLOR_ROOT = "#008b8b" 

func GenerateRootInterpretation(db_idx int,  root string) string {
    ch_root_mean := make(  chan string )
    ch_synroot := make( chan string )
    ch_root_example := make( chan string )

    go generateFormatedMeaning( root, ch_root_mean )
    go getSynonymsRoots( db_idx , root , ch_synroot ) 
    go generateRootWordExamples( db_idx, root, ch_root_example ) 

    rootMean := <- ch_root_mean
    synRoots := <- ch_synroot
    wordExample := <- ch_root_example 
    return fmt.Sprintf( "\n\n%s\n\n%s\n\n%s" , rootMean , synRoots , wordExample   )
}

//======= func for main dict =========================

var key_4_key = []byte("mebusy key 4 key")
var key_4_dict = []byte("mebusy 2018 dict")
var key_4_indo = "gj+uRZLqDRsVlEE5sTNEoAXMGR8ODOl2/OZyvvgkqDYpqy0uP6+Snkz65M9O11LM"

func SearchWordLike( db_idx int, _word string ) string {
    sql := sqlManager.GetInstance()
    if db_idx<0 || db_idx > sql.ConnectionCount() {
        log.Println( "SearchWordLike","wrong db index" )
        return ""
    }
    db := sql.AllConns[ db_idx ]

    MAX_CELL_NUMBER := 30
    cmdText := fmt.Sprintf("SELECT word,root1,root2 FROM dict WHERE word like \"%%%[1]s%%\"  " +
        "ORDER BY (CASE WHEN word = \"%[1]s\" COLLATE NOCASE THEN 1 WHEN word LIKE \"%[1]s%%\" THEN 2 ELSE 3 END) limit %[2]d   " , 
        _word , MAX_CELL_NUMBER )
    rows, err := db.Query( cmdText ) 
    if HasErr(err) {
        return ""
    }
    defer rows.Close()
    
    for rows.Next() {
        var word string 
        var root1 string 
        var root2 string 
        err := rows.Scan( &word , &root1, &root2 )
        if HasErr(err) {
            break
        }
        fmt.Println( word, root1, root2  )
    }
    return ""
}
var sepline = strings.Repeat( "·", 60 ) 
var sep = "\n<color=#DAEDF1>" + sepline + "</color>\n"

func GetWordInterpretation( db_idx int , word string  ) string {
    sql := sqlManager.GetInstance()
    if db_idx<0 || db_idx > sql.ConnectionCount() {
        log.Println( "GetWordInterpretation","wrong db index" )
        return ""
    }
    db := sql.AllConns[ db_idx ]

    // get main word data 
    cmdText := fmt.Sprintf( "SELECT etymology,desc_en,desc_cn,indo_roots,ex,voice,root1,root2  FROM dict WHERE word = \"%s\"  " , word  )
    rows, err := db.Query( cmdText )
    if HasErr(err) { return "" }
    defer rows.Close()

    var etymology,desc_en,desc_cn,str_indo_roots,ex,root1,root2 string 
    var voice int 
    if rows.Next() {
        err := rows.Scan( &etymology,&desc_en,&desc_cn,&str_indo_roots,&ex,&voice,&root1,&root2 )
        if HasErr(err) { return "" }
    }

    etymology = strings.TrimSpace( cryptor.Decrypt_CBC_AES( etymology, key_4_dict )   )

    dict_indo_roots := make( map[string]string ,2  )
    if str_indo_roots != "" {
        indo_roots := strings.Split(  str_indo_roots, ",")
        var s bytes.Buffer 
        s.WriteString( "SELECT word, desc  FROM indo_root WHERE " ) 
        r := strings.NewReplacer("(", "", ")", "")
        for i,v := range indo_roots {
            s.WriteString(  fmt.Sprintf( "%[1]v word=\"%[2]v\" " , If(i==0,"","or") , r.Replace( v ) ))    
        }

        cmdText := s.String()
        rows, err := db.Query( cmdText )
        if HasErr(err) { return "" }
        defer rows.Close()
        
        for  rows.Next() {
            var indo_word , indo_desc string
            err := rows.Scan( &indo_word , &indo_desc  )     
            if HasErr(err) { return "" }

            k:= key_4_key
            k2 := []byte(cryptor.Decrypt_CBC_AES( key_4_indo, k  )[:16])
            decrypt_indo_desc :=cryptor.Decrypt_CBC_AES (indo_desc, k2);

            // fmt.Println( indo_word, decrypt_indo_desc  )
            dict_indo_roots[indo_word] = decrypt_indo_desc 
        }
    }
    var entire_text string
    {
        var indo_roots  bytes.Buffer 
        for k,v := range dict_indo_roots {
            indo_roots.WriteString( fmt.Sprintf( "\n<b><color=%[3]s>%[1]s</color></b>: %[2]s", k,v, COLOR_ROOT  ) )
        }
        indo_root := indo_roots.String()

        if etymology != "" {
            etymology =  "<b><color=blue>Roots</color></b>:\n" + etymology

            // 如果有的话， 两个放一起显示，分割线现在不加
            if indo_root == "" {
                etymology += sep  
            }
        }
        if indo_root != ""  {  // 若没有不显示
            indo_root +=  sep;
        }

        ch_root1_mean := make(  chan string )
        ch_root2_mean := make(  chan string )
        root_info := "";
        if root1 != "" {
            go generateFormatedMeaning ( root1 , ch_root1_mean )
        }
        if root2 != "" {
            go generateFormatedMeaning ( root2 , ch_root2_mean )
        }
        if root1 != "" { root_info += <- ch_root1_mean }
        if root2 != "" { root_info += "\n"+ <- ch_root2_mean }

        if ex != "" {
            // 后期可能添加 单独的助记，也既 虽然没有root， 但是 有助记
            if root_info != "" { root_info += "\n" }
            root_info += fmt.Sprintf( "<color=blue>助记:  </color>%s" , ex )
        }
        
        if root_info != "" { root_info += sep }
        if desc_cn != ""  { desc_cn+= sep }

         entire_text = fmt.Sprintf( "\n%s%s%s%s%s\n\n\n\n" , etymology   , indo_root ,  root_info , desc_cn ,desc_en    )
    }

    return entire_text
}


