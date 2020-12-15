package main

import (
	"net"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bufio"
	"os"
	"sync"
	"strings"
)

var gopath = os.Getenv("GOPATH")

type Cloudflareapi struct {
	Result struct {
		Ipv4Cidrs []string `json:"ipv4_cidrs"`
		Ipv6Cidrs []string `json:"ipv6_cidrs"`
		Etag      string   `json:"etag"`
	} `json:"result"`
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
}

type Cloudfrontapi struct {
	CLOUDFRONTGLOBALIPLIST       []string `json:"CLOUDFRONT_GLOBAL_IP_LIST"`
	CLOUDFRONTREGIONALEDGEIPLIST []string `json:"CLOUDFRONT_REGIONAL_EDGE_IP_LIST"`
}

type Fastyapi struct {
	Addresses     []string `json:"addresses"`
	Ipv6Addresses []string `json:"ipv6_addresses"`
}

func Cloudflarev4() []string {
	bytes, _ := ioutil.ReadFile(gopath + "/src/github.com/kawakatz/kawacdncheck/json/cloudflare.json")
	cloudflareapi := new(Cloudflareapi)
	json.Unmarshal(bytes, &cloudflareapi)
	
	return cloudflareapi.Result.Ipv4Cidrs
}

func Cloudfrontv4() []string {
	bytes, _ := ioutil.ReadFile(gopath + "/src/github.com/kawakatz/kawacdncheck/json/cloudfront.json")
	cloudfrontapi := new(Cloudfrontapi)
	json.Unmarshal(bytes, &cloudfrontapi)
	
	return append(cloudfrontapi.CLOUDFRONTGLOBALIPLIST, cloudfrontapi.CLOUDFRONTREGIONALEDGEIPLIST...)
}

func Fastyv4() []string {
	bytes, _ := ioutil.ReadFile(gopath + "/src/github.com/kawakatz/kawacdncheck/json/fasty.json")
	fastyapi := new(Fastyapi)
	json.Unmarshal(bytes, &fastyapi)
	
	return fastyapi.Addresses
}

func check(ip string) string {
	cloudflarev4 := Cloudflarev4()
	cloudfrontv4 := Cloudfrontv4()
	fastyv4 := Fastyv4()
	targetip := net.ParseIP(ip)

	for _, iprange := range cloudflarev4 {
		_, subnet, _ := net.ParseCIDR(iprange)
		if subnet.Contains(targetip){
			return "CloudFlare"
		}
	}

	for _, iprange := range cloudfrontv4 {
		_, subnet, _ := net.ParseCIDR(iprange)
		if subnet.Contains(targetip){
			return "CloudFront"
		}
	}

	for _, iprange := range fastyv4 {
		_, subnet, _ := net.ParseCIDR(iprange)
		if subnet.Contains(targetip){
			return "Fasty"
		}
	}

	revdns, err := net.LookupAddr(ip)
	if err == nil {
		revdomain := revdns[0]
		if strings.Contains(revdomain, "deploy.static.akamaitechnologies.com") {
			return "Akamai"
		}
	}

	return "None"
}

func main(){
	domains := make(chan string)

	var wg sync.WaitGroup
	naked := []string{}

	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			for domain := range domains {
				ipv4, err := net.ResolveIPAddr("ip4", domain)
				if err != nil {
					continue
				}
				result := check(ipv4.String())
				if result=="None"{
					naked = append(naked, domain)
				}else{
					//fmt.Println(domain + ": " + result)
				}
			}

			wg.Done()
		}()
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		domains <- sc.Text()
	}
	close(domains)

	wg.Wait()
	for _, domain := range naked {
		fmt.Println(domain)
	}

}
