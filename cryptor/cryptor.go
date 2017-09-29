package cryptor

import (
    "encoding/base64"
    "crypto/aes"
    "crypto/cipher"
    "log"
)

func Decrypt_CBC_AES( base64str string , key[]  byte ) string {
    // _entireText = System.Convert.FromBase64String (base64str)    
    ciphertext,err := base64.StdEncoding.DecodeString( base64str)
    if err != nil {
        log.Println( "base64 decode error" )       
        return ""
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        log.Println( "aes.NewCipher error" )       
        return ""
    }

    // The IV needs to be unique, but not secure. Therefore it's common to
    // include it at the beginning of the ciphertext.
    if len(ciphertext) < aes.BlockSize {
        log.Println("ciphertext too short")
        return ""
    }
    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    // CBC mode always works in whole blocks.
    if len(ciphertext)%aes.BlockSize != 0 {
        log.Println("ciphertext is not a multiple of the block size")
        return ""
    }

    mode := cipher.NewCBCDecrypter(block, iv)

    // CryptBlocks can work in-place if the two arguments are the same.
    mode.CryptBlocks(ciphertext, ciphertext)

    // If the original plaintext lengths are not a multiple of the block
    // size, padding would have to be added when encrypting, which would be
    // removed at this point. For an example, see
    // https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
    // critical to note that ciphertexts must be authenticated (i.e. by
    // using crypto/hmac) before being decrypted in order to avoid creating
    // a padding oracle.    


    return string(ciphertext)
}
