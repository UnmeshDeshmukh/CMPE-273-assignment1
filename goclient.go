package main

import (
	
	"fmt"
	"log"
	"net/rpc/jsonrpc"
//	"encoding/json"
//	"net/http"
 //	"io/ioutil"
 	"os"
 	"strconv"
)


type Stock struct {
  List struct {
    Resources []struct {
      Resource struct {
        Fields struct {
          Name    string `json:"name"`
          Price   string `json:"price"`
          Symbol  string `json:"symbol"`
          Ts      string `json:"ts"`
          Type    string `json:"type"`
          UTCTime string `json:"utctime"`
          Volume  string `json:"volume"`
        } `json:"fields"`
      } `json:"resource"`
    } `json:"resources"`
  } `json:"list"`
}
type Request struct{
	Str string
	Id int
	Budget float64
}

type Response struct{
  Tradeid int
  Stocks[30] string
  Unvested float64
  Totalamt float64
  Shares float64
  Count int
  InvestedPrice float64
}
type PortfolioRequest struct{
  Unvestedamt float64
  ID int
}

type PortfolioResponse struct{
    Stoks string
    MarketValue float64
    Amount float64
}
type St int

func BuyStocks(){
	var args Request
	//var stock Stock
	client, err := jsonrpc.Dial("tcp", "localhost:1337")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args.Str = os.Args[1]
	temp := os.Args[2]

	var floatvalue float64
	floatvalue, _ = strconv.ParseFloat(temp, 64)
	//fmt.Println(floatvalue)
	args.Budget = floatvalue
	var reply Response

	err = client.Call("St.BuyingStocks",args, &reply)
	if err != nil {
		log.Fatal("st error:", err)
	}
	//fmt.Println(reply)
	//fmt.Println("ID-Symbol-ShareNo-Price-InvestedAmt")
  for loop:= 0;loop< reply.Count; loop++{
     fmt.Printf("%s\n",reply.Stocks[loop])
  }

}

func ViewPortfolio(){
    var portfolioRequest PortfolioRequest
    
    integer,_ := strconv.ParseInt(os.Args[1], 10, 32)
    portfolioRequest.ID = int(integer)
    client, err := jsonrpc.Dial("tcp", "localhost:1337")
    if err != nil {
        log.Fatal("dialing:", err)
    }
    
    var Res PortfolioResponse

    
    
        err = client.Call("St.Portfolio", portfolioRequest.ID - 1, &Res)
        if err != nil {
            log.Fatal("St error:", err)
        }

        fmt.Printf("Your Portfolio is---- \n")
        fmt.Printf("\nStocks Bought : ")
        fmt.Printf("\"%s", Res.Stoks)        
        fmt.Printf("\nMarket Value : %f", Res.MarketValue)
        fmt.Printf("\nUnvested Amount : %f\n", Res.Amount)

}
func main() {
	
	if len(os.Args) == 3 {
        BuyStocks()
    }else if len(os.Args) == 2{
        ViewPortfolio()
    }else{
        fmt.Println("localhost:1337")
        log.Fatal(1)
    }
}
