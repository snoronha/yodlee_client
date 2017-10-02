package model

import (
    _ "fmt"
)

type Transaction struct {
    TxnId                uint
    TxnType              string
    Amount               int64    // amount in cents
    Currency             string
    BaseType             string
    CategoryType         string
    CategoryId           uint
    Category             string
    Description          string
    DescriptionSimple    string
    AccountId            uint
}
