package extension

import (
    "log"
)


func If(condition bool, trueVal, falseVal interface{}) interface{} {
    if condition {
        return trueVal
    }
    return falseVal
}


func HasErr(err error) bool {
    if err != nil {
        log.Println( "err:",err.Error() )
        // panic(err)
        return true 
    }
    return false 
}
