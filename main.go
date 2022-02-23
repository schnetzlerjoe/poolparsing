package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
  "strings"
  "github.com/schnetzlerjoe/poolparsing/cosmos"
)

var client = &http.Client{Timeout: 10 * time.Second}

type OsmosisAsset struct {
  Base string `json:"base"`
  Name string `json:"name"`
  Display string `json:"display"`
  Symbol string `json:"symbol"`
}

type OsmosisRequest struct {
  Assets []OsmosisAsset `json:"assets"`
}

type CosmosAsset struct {
}

type Pool struct {
  Id string `json:"id"`
  CoinDenoms []string `json:"reserve_coin_denoms"`
}

type Pools struct {
  Pools []Pool `json:"pools"`
}

type DenomTrace struct {
  Path string `json:"path"`
  BaseDenom string `json:"base_denom"`
}

type DenomDetail struct {
  DenomTrace DenomTrace `json:"denom_trace"`
}

func denomDetails(denom string) (ret DenomDetail) {
  hash := strings.Split(denom, "/")[1]
  response, err := client.Get("https://api.cosmos.network/ibc/apps/transfer/v1/denom_traces/" + hash)
  if err != nil {
    log.Fatal(err)
  }
  responseData, err := ioutil.ReadAll(response.Body)
  if err != nil {
    log.Fatal(err)
  }

  responseString := string(responseData)

  data := DenomDetail{}

  json.Unmarshal([]byte(responseString), &data)

  return data
}

func getOsmosisAssets() (ret []OsmosisAsset) {
  url := "https://raw.githubusercontent.com/osmosis-labs/assetlists/main/osmosis-1/osmosis-1.assetlist.json"

  response, err := client.Get(url)
  if err != nil {
      log.Fatal(err)
  }
  defer response.Body.Close()

  responseData, err := ioutil.ReadAll(response.Body)
  if err != nil {
      log.Fatal(err)
  }

  responseString := string(responseData)

  data := OsmosisRequest{}
  json.Unmarshal([]byte(responseString), &data)

  return data.Assets
}

func getCosmosPoolsRaw() (ret []Pool) {
  response, err := client.Get("https://api.cosmos.network/cosmos/liquidity/v1beta1/pools")
  
  if err != nil {
    fmt.Println(err)
  }

  defer response.Body.Close()

  responseData, err := ioutil.ReadAll(response.Body)

  responseString:= string(responseData)

  data := Pools{}
  json.Unmarshal([]byte(responseString), &data)

  return data.Pools
}

func getCosmosPools() (ret []CosmosAsset) {
  pools := getCosmosPoolsRaw()
  for i := range pools {
    denoms := pools[i].CoinDenoms
    for i := range denoms  {
      fmt.Println(denoms[i])
    }
  }
  return ret
}

func getCosmosPoolsState() (state Pools) {
  return state
}

func main() {
  cosmos.cosmos()
}
