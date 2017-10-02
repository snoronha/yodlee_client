package model

import (
    _ "fmt"
)

type Account struct {
    AccountId            uint
    AccountNumber        string
    AccountStatus        string
    AccountType          string
    Balance              float64
    BalanceCurrency      string
    CashValue            float64
    CashValueCurrency    string
}
