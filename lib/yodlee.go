package lib

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
    "github.com/karlseguin/typed"
    "yodlee_client/model"
)

func CobrandLogin(cobrandUrl string, cobrandLogin string, cobrandPassword string) map[string]string {
    retMap     := make(map[string]string)
    paramMap   := map[string]map[string]string{
        "cobrand": {
            "cobrandLogin": cobrandLogin,
            "cobrandPassword": cobrandPassword,
            "locale": "en_US",
        },
    }
    paramStr, _ := json.Marshal(paramMap)
    req, err := http.NewRequest("POST", cobrandUrl, bytes.NewBuffer(paramStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    if nil != err {
        fmt.Println("error while reading the body", err)
        return retMap
    }

    // fmt.Printf("%s\n", string(body))
    
    typed, err := typed.JsonString(string(body))
    if err != nil {
        fmt.Println(err)
    }
    session := typed.Object("session")
    retMap["cobrandId"]  = strconv.Itoa(typed.Int("cobrandId"))
    retMap["cobSession"] = session.String("cobSession")
    return retMap
}

func UserLogin(loginUrl string, loginName string, password string, cobSession string) map[string]string {
    retMap     := make(map[string]string)
    paramMap   := map[string]map[string]string{
        "user": {
            "loginName": loginName,
            "password": password,
            "locale": "en_US",
        },
    }
    paramStr, _ := json.Marshal(paramMap)
    req, err := http.NewRequest("POST", loginUrl, bytes.NewBuffer(paramStr))
    req.Header.Set("Authorization", "cobSession=" + cobSession)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    if nil != err {
        fmt.Println("error while reading the body", err)
        return retMap
    }

    // fmt.Printf("%s\n", string(body))

    typed, err := typed.JsonString(string(body))
    if err != nil {
        fmt.Println(err)
    }
    user    := typed.Object("user")
    name    := user.Object("name")
    session := user.Object("session")
    retMap["roleType"]    = user.String("roleType")
    retMap["userId"]      = strconv.Itoa(user.Int("id"))
    retMap["loginName"]   = user.String("loginName")
    retMap["firstName"]   = name.String("first")
    retMap["lastName"]    = name.String("lastfirst")
    retMap["userSession"] = session.String("userSession")
    return retMap
}

func GetTransactions(txnUrl string, cobSession string, userSession string) []model.Transaction {
    txnArr      := []model.Transaction{}
    paramMap    := map[string]map[string]string{}
    paramStr, _ := json.Marshal(paramMap)
    req, err := http.NewRequest("GET", txnUrl, bytes.NewBuffer(paramStr))
    req.Header.Set("Authorization", "userSession=" + userSession + ", cobSession=" + cobSession)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    if nil != err {
        fmt.Println("error while reading the body", err)
        return txnArr
    }

    // fmt.Printf("%s\n", string(body))

    typed, err := typed.JsonString(string(body))
    if err != nil {
        fmt.Println(err)
    }
    transactions := typed.Objects("transaction")
    fmt.Printf("\nNUMBER OF TXNS = %d\n", len(transactions))
    for _, txn := range transactions {
        amountBlob := txn.Object("amount")
        descrBlob  := txn.Object("description")
        transaction := model.Transaction{
            TxnId:             uint(txn.Int("id")),
            Amount:            int64(amountBlob.Float("amount") * 100), // amount in cents
            Currency:          amountBlob.String("currency"),
            BaseType:          txn.String("baseType"),
            CategoryType:      txn.String("categoryType"),
            CategoryId:        uint(txn.Int("categoryId")),
            Category:          txn.String("category"),
            Description:       descrBlob.String("original"),
            DescriptionSimple: descrBlob.String("simple"),
            TxnType:           txn.String("type"),
            AccountId:         uint(txn.Int("accountId")),
        }
        txnArr = append(txnArr, transaction)
    }
    return txnArr
}

func GetAccounts(accountsUrl string, cobSession string, userSession string) []model.Account {
    accountsArr := []model.Account{}
    paramMap    := map[string]map[string]string{}
    paramStr, _ := json.Marshal(paramMap)
    req, err := http.NewRequest("GET", accountsUrl, bytes.NewBuffer(paramStr))
    req.Header.Set("Authorization", "userSession=" + userSession + ", cobSession=" + cobSession)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    if nil != err {
        fmt.Println("error while reading the body", err)
        return accountsArr
    }

    // fmt.Printf("%s\n", string(body))

    typed, err := typed.JsonString(string(body))
    if err != nil {
        fmt.Println(err)
    }
    accounts := typed.Objects("account")
    fmt.Printf("\nNUMBER OF ACCOUNTS = %d\n\n", len(accounts))
    for _, acct := range accounts {
        balance   := acct.Object("balance")
        cashValue := acct.Object("cashValue")
        account   := model.Account{
            AccountId:         uint(acct.Int("id")),
            AccountNumber:     acct.String("accountNumber"),
            AccountStatus:     acct.String("accountStatus"),
            AccountType:       acct.String("accountType"),
            Balance:           balance.Float("amount"),
            BalanceCurrency:   balance.String("currency"),
            CashValue:         cashValue.Float("amount"),
            CashValueCurrency: cashValue.String("currency"),
        }
        accountsArr = append(accountsArr, account)
    }
    return accountsArr
}

func GetFastlinkToken(fastlinkTokenUrl string,
    cobrandName string, appIds string, cobSession string, userSession string) []model.FastlinkToken {
    fastlinkTokensArr := []model.FastlinkToken{}
    paramMap    := map[string]string{}
    paramStr, _ := json.Marshal(paramMap)
    req, err := http.NewRequest("GET", fastlinkTokenUrl + "?cobrandName=" + cobrandName + "&appIds=" + appIds, bytes.NewBuffer(paramStr))
    req.Header.Set("Authorization", "userSession=" + userSession + ", cobSession=" + cobSession)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    if nil != err {
        fmt.Println("error while reading the body", err)
        return fastlinkTokensArr
    }

    // fmt.Printf("%s\n", string(body))

    typed, err := typed.JsonString(string(body))
    if err != nil {
        fmt.Println(err)
    }
    user         := typed.Object("user")
    accessTokens := user.Objects("accessTokens")
    for _, accessToken := range accessTokens {
        fastlinkToken  := model.FastlinkToken{
            AppId: accessToken.String("appId"),
            Value: accessToken.String("value"),
            Url:   accessToken.String("url"),
        }
        fastlinkTokensArr = append(fastlinkTokensArr, fastlinkToken)
    }
    return fastlinkTokensArr
}
