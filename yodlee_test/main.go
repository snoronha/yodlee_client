package main

import (
	"fmt"
    "github.com/spf13/viper"
    "yodlee_client/lib"
)

func main() {
    viper.SetConfigName("app")     // no need for file extension
    viper.AddConfigPath("../config")  // set path to config file

    err := viper.ReadInConfig()
    if err != nil {
        fmt.Println("Config file not found ...")
        return
    }
    cobrandUrl       := viper.GetString("development.cobrandUrl")
    cobrandLogin     := viper.GetString("development.cobrandLogin")
    cobrandPassword  := viper.GetString("development.cobrandPassword")
    loginUrl         := viper.GetString("development.loginUrl")
    fastlinkTokenUrl := viper.GetString("development.fastlinkTokenUrl")
    loginName        := "sbMemsnoronha2"
    password         := "sbMemsnoronha2#123"
    accountUrl       := viper.GetString("development.accountUrl")

    //------ Cobrand Login ------//
    cobMap          := lib.CobrandLogin(cobrandUrl, cobrandLogin, cobrandPassword)
    cobSession      := cobMap["cobSession"]
    fmt.Printf("cobSession: %s\n", cobSession)

    if len(cobSession) <= 0 {
        fmt.Printf("No valid Cobrand Session: %s\n", cobSession)
        return
    }
    
    //------ User Login ------//
    loginMap        := lib.UserLogin(loginUrl, loginName, password, cobSession)
    userSession     := loginMap["userSession"]
    fmt.Printf("userSession: %s\n", userSession)

    if len(userSession) <= 0 {
        fmt.Printf("No valid User Session: %s\n", userSession)
        return
    }

    //------ Get FastlinkToken ------//
    cobrandName      := "yodlee"
    appIds           := "10003600"
    fastlinkTokenArr := lib.GetFastlinkToken(fastlinkTokenUrl, cobrandName, appIds, cobSession, userSession)
    fmt.Printf("--------- FASTLINK TOKEN ---------\n%v\n\n", fastlinkTokenArr)
    
    //------ Get Accounts ------//
    // _ = accountUrl
    accountArr := lib.GetAccounts(accountUrl, cobSession, userSession)
    fmt.Printf("--------- ACCT ARRAY ---------\n%v\n\n", accountArr)
    
    //------ Get Transactions ------//
    txnUrl  := "https://developer.api.yodlee.com/ysl/restserver/v1/transactions"
    txnArr  := lib.GetTransactions(txnUrl, cobSession, userSession)
    fmt.Printf("-------- TXN ARRAY --------\n%v\n\n", txnArr)
}
