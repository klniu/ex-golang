package main

import (
    "log"
)

const (
    CLIENT_LONG_PASSWORD uint32 = 1 << iota
    CLIENT_FOUND_ROWS
    CLIENT_LONG_FLAG
    CLIENT_CONNECT_WITH_DB
    CLIENT_NO_SCHEMA
    CLIENT_COMPRESS
    CLIENT_ODBC
    CLIENT_LOCAL_FILES
    CLIENT_IGNORE_SPACE
    CLIENT_PROTOCOL_41
    CLIENT_INTERACTIVE
    CLIENT_SSL
    CLIENT_IGNORE_SIGPIPE
    CLIENT_TRANSACTIONS
    CLIENT_RESERVED
    CLIENT_SECURE_CONNECTION
    CLIENT_MULTI_STATEMENTS
    CLIENT_MULTI_RESULTS
    CLIENT_PS_MULTI_RESULTS
    CLIENT_PLUGIN_AUTH
    CLIENT_CONNECT_ATTRS
    CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA
)

func main() {
    log.Printf("%d\n", CLIENT_LONG_PASSWORD)
    log.Printf("%b\n", CLIENT_LONG_PASSWORD)

    log.Printf("%d\n", CLIENT_FOUND_ROWS)
    log.Printf("%b\n", CLIENT_FOUND_ROWS)

    log.Printf("%d\n", CLIENT_LONG_FLAG)
    log.Printf("%b\n", CLIENT_LONG_FLAG)

    log.Printf("%d\n", CLIENT_CONNECT_WITH_DB)
    log.Printf("%b\n", CLIENT_CONNECT_WITH_DB)

    log.Printf("%d\n", CLIENT_NO_SCHEMA)
    log.Printf("%b\n", CLIENT_NO_SCHEMA)

    capability1 := uint32(
        CLIENT_PROTOCOL_41 |
            CLIENT_SECURE_CONNECTION |
            CLIENT_LONG_PASSWORD |
            CLIENT_TRANSACTIONS |
            CLIENT_LOCAL_FILES |
            CLIENT_LONG_FLAG,
    )

    log.Printf("capability1: %b\n", capability1)

    capability2 := uint32(
        CLIENT_PROTOCOL_41 |
            CLIENT_SECURE_CONNECTION |
            CLIENT_LONG_PASSWORD |
            CLIENT_TRANSACTIONS |
            CLIENT_LOCAL_FILES,
    )

    log.Printf("capability2: %b\n", capability2)

    if (capability1 & CLIENT_LONG_FLAG) != 0 {
        log.Printf("%b\n", capability1&CLIENT_LONG_FLAG)
        capability2 |= CLIENT_LONG_FLAG
    }

    log.Printf("capability2: %b\n", capability2)
}
