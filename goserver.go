package main
import(
  "fmt"
  //"log"
  "os"
  "net/rpc/jsonrpc"
  "net/http"
  "net/rpc"
  "net"
  "io/ioutil"
  "encoding/json"
  "strings"
  "strconv"
  "math"
)

type Request struct{
  Str string
  Id int
  Budget float64
  Unvestedamt float64
}

var TempId int

type Response struct{
  Tradeid int
  Stocks[30] string
  Totalamt  float64
  Unvested float64
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

var Slice []Response

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
type St int




func (t *St) BuyingStocks(args *Request, reply *Response) error{
  
  var stock Stock
  reply.Stocks[0] = args.Str
 // reply.Id = args.Id
  reply.Totalamt = args.Budget 
  namearr := strings.Split(reply.Stocks[0],",")
  TempId++
  reply.Tradeid = TempId 

  for i:= range(namearr){
    stocksymbol := strings.Split(namearr[i],":")
    stockname := stocksymbol[0]
    reply.Stocks[i] = stockname
    percent,_ := strconv.ParseFloat(stocksymbol[1],64)
    amount := reply.Totalamt * (percent/100)
    url :=  fmt.Sprintf("http://finance.yahoo.com/webservice/v1/symbols/%s/quote?format=json", stockname)
    res, err := http.Get(url)
    if err!=nil {
       panic(err)
    }
    body, err := ioutil.ReadAll(res.Body)
    err = json.Unmarshal(body, &stock)
    price, _:= strconv.ParseFloat(stock.List.Resources[0].Resource.Fields.Price ,64)
    noofshares:= amount/price
    noofshares = math.Floor(noofshares)
    reply.Shares = noofshares
    reply.InvestedPrice = noofshares*price
    reply.Unvested += amount - (noofshares*price)
    reply.Totalamt =(reply.Totalamt-amount)+reply.Unvested
    reply.Stocks[i] = reply.Stocks[i] + ":" + strconv.FormatFloat(noofshares,'f',0,64) + ":$" + strconv.FormatFloat(price,'f',6,64) + ":" + strconv.FormatFloat(reply.Totalamt,'f',2,64) 
    reply.Count++
  //  fmt.Printf("Length: %d", len(Slice))
   //     fmt.Printf("\n***%s***", reply.Stocks[i])
  }
   Slice = append(Slice, *reply)
return nil
}








func (t *St) Portfolio(temp int, reply *PortfolioResponse) error{

  request := Slice[temp]
  tcount := request.Count
  
  for i := 0; i < tcount; i++ {
    Data := strings.Split(request.Stocks[i], ":")
    url := fmt.Sprintf("http://finance.yahoo.com/webservice/v1/symbols/%s/quote?format=json",Data[0])
        urlRes,err := http.Get(url)
        if err != nil{
            panic(err)
        }
        defer urlRes.Body.Close()

        body,err := ioutil.ReadAll(urlRes.Body)
        if err != nil{
            panic(err)
        }

       var stock Stock
        err = json.Unmarshal(body, &stock)
        if err != nil{
            panic(err)
        }
    
  

       
        tempPrice,_ := strconv.ParseFloat(stock.List.Resources[0].Resource.Fields.Price, 64) 
        reply.Stoks = reply.Stoks+Data[0]+":"+Data[1]
        prev,_ := strconv.ParseFloat(Data[2], 64)
        if tempPrice > prev{
            reply.Stoks +=":+$"+strconv.FormatFloat(tempPrice,'f', 3, 64)+"\"; "
        }else if tempPrice < prev{
            reply.Stoks +=":-$"+strconv.FormatFloat(tempPrice,'f', 3, 64)+"\"; "
        }else{
            reply.Stoks +=":==$"+strconv.FormatFloat(tempPrice,'f', 3, 64)+"\"; "
        }
        shareno,_ := strconv.ParseFloat(Data[1], 64)
        tempMarketValue :=  shareno * tempPrice
        reply.MarketValue += tempMarketValue
    }
    reply.Amount = request.Totalamt
    return nil
}













func main(){
  st := new(St)
  rpc.Register(st)
  tcpAddr, err :=net.ResolveTCPAddr("tcp","localhost:1337")
  checkError(err)
  listener, err := net.ListenTCP("tcp", tcpAddr)
  checkError(err)
 // rpc.Accept(listener)
for {
    conn, err := listener.Accept()
    if err != nil {
      continue
    }
    jsonrpc.ServeConn(conn)
  }



}
  func checkError(err error){
    if err!= nil {
       fmt.Println("Error",err.Error())
       os.Exit(1)
    }
  }

  