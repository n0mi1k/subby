package main

import (
	"fmt"
	"sort"
	"reflect"
	"bufio"
	"os"
	"log"
	"strings"
	"net/http"
	"math/rand"
	"net"
	"strconv"
	"crypto/tls"
	"github.com/alexflint/go-arg"
	"time"
)

// Sample run for DNS Mode: subby -u example.com -w shubs-subdomains.txt -d 100
// Sample run for Web Mode (HTTPS): subby -u https://example.com -w shubs-subdomains.txt -d 100 -t 20 -r "200, 301, 302"
// Sample run for Web Mode (HTTP): subby -u http://example.com -w shubs-subdomains.txt -d 100 -t 20 -r "200, 301, 302"

var (
	lgreen  = "\033[92m"
	lpurple = "\033[95m"
	lred	= "\033[91m"
	yellow  = "\033[33m"
	lcyan   = "\033[96m"
	reset   = "\033[0m"
)


func readWordlist(filename string) []string {
	logger := log.New(os.Stderr, "", 0)
    file, err := os.Open(filename)
    if err != nil {
        logger.Fatal(lred + "[ERR] " + err.Error() + reset)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
	const maxCapacity int = 20000000
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	words := make([]string, 0)
    for scanner.Scan() {
		words = append(words, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        logger.Fatal("[ERR] " + err.Error())
    }
	return words
}


func writeToFile(filename, content string) error {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(content + "\n")
    if err != nil {
        return err
    }
    return nil
}


func enumHttpSubdomain(url string, subdomain string, httpsFlag bool, codes []int, delay int, timeout int, outfile string) {
	suburl := ""
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout: time.Duration(timeout) * time.Second,
	}

	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	} 
	
	if httpsFlag {
		if strings.Contains(url, "www.") {
			suburl = strings.Replace(url, "https://www.", "https://" + subdomain + ".", 1)
		} else {
			suburl = strings.Replace(url, "https://", "https://" + subdomain + ".", 1)
		}
	} else {
		if strings.Contains(url, "www.") {
			suburl = strings.Replace(url, "http://www.", "http://" + subdomain + ".", 1)
		} else {
			suburl = strings.Replace(url, "http://", "http://" + subdomain + ".", 1)
		}
	}
	
	resp, err := client.Get(suburl)
	if err == nil {
		if resp.StatusCode == 429 {
			for resp.StatusCode == 429 {
				fmt.Println(yellow + "[WRN] Rate Limited.. Reduce Threads, Increase Delay! Sleep 30s... " + reset)
				time.Sleep(30 * time.Second)
				resp, err = client.Get(suburl)
			}
		} else {
			if codes == nil {
				if resp.StatusCode >= 200 && resp.StatusCode <= 599 && resp.StatusCode != 429 {
					loggedResp := suburl + " [" + strconv.Itoa(resp.StatusCode) + "]"
					fmt.Println(loggedResp)
					writeToFile(outfile, loggedResp)
				}
			} else {
				for _, code := range codes {
					if resp.StatusCode == code {
						loggedResp := suburl + " [" + strconv.Itoa(resp.StatusCode) + "]"
						fmt.Println(loggedResp)
						writeToFile(outfile, loggedResp)
					}
				}
			}
		}
	}
}


func CheckWildcardDNS(url string) (bool, []string) {
	fmt.Println(lcyan + "[+] Checking for wildcard DNS (This may take a minute)..." + reset)

	var testCount = 5
	var returnIPCount = 0
	var wildcardIPs []string

	for i := 0; i <= testCount; i++ {
		host := strconv.Itoa(rand.Intn(9000000000) + 1000000000) + "." + url
		wildcardAddr, _:= net.LookupHost(host)
		if wildcardAddr != nil {
			wildcardIPs = wildcardAddr
			returnIPCount++
		}
	}

	if (returnIPCount == testCount) || returnIPCount >= 3 {
		return true, wildcardIPs
	}
	return false, nil
}


func enumDNSSubdomain(url string, subdomain string, delay int, wildcardDNSFlag bool, wildcardAddr []string, outfile string) {
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	} 
	host := subdomain + "." + url
	address, err:= net.LookupHost(host)

	if err == nil {
		if wildcardDNSFlag {
			if !(isEqualSlices(address, wildcardAddr)) {
				fmt.Println(host)
				writeToFile(outfile, host)
			}
		} else {
			fmt.Println(host)
			writeToFile(outfile, host)
		}
	}
}


func isEqualSlices(slice1, slice2 []string) bool {
    if len(slice1) != len(slice2) {
        return false
    }
	
    sort.Strings(slice1)
    sort.Strings(slice2)

    return reflect.DeepEqual(slice1, slice2)
}


type args struct {
	URL string `arg:"-u, --url, required" help:"URL to enumerate [required]"`
	Wordlist string `arg:"-w, --wordlist, required" help:"Wordlist to use [required]"`
	Delay int `arg:"-d, --delay" default:"0" help:"Delay in milliseconds for each goroutine request"`
	Response string `arg:"-r" help:"Show only these status codes separated by comma: e.g 200, 301"`
	Threads int `arg:"-t, --threads" default:"50" help:"Number of concurrent Goroutines"`
	Timeout int `arg:"-s, --timeout" default:"2" help:"Set max timeout in seconds for each request"`
	Output string `arg:"-o, --output" help:"File to output subdomains"`
}


func (args) Description() string {
	return "An uber fast next-generation subdomain enumeration toolkit"
}


func main() {
	logger := log.New(os.Stderr, "", 0)
	var httpsFlag bool = false
	var dnsFlag bool = false
	var wildcardDNSFlag bool = false
	var wildcardAddr []string
	var args args 
	arg.MustParse(&args)

	url := args.URL
	codes := args.Response
	threads := args.Threads	
	wordfile :=	args.Wordlist
	delay := args.Delay	
	timeout := args.Timeout
	output := args.Output

	wordlist := readWordlist(wordfile)

	subbyArt := lcyan + `
    ____  _     ____  ____ ___  _
   / ___\/ \ /\/  _ \/  _ \\  \//
   |    \| | ||| | //| | // \  / 
   \___ || \_/|| |_\\| |_\\ / /  
   \____/\____/\____/\____//_/ v1.0

        github.com/n0mi1k   				  
	` + reset
	fmt.Println(subbyArt)

	fmt.Println(lgreen + "[+] Target Domain: " + reset + url)
	fmt.Println(lpurple + "[+] Wordlist: " + reset + wordfile)
	fmt.Println(yellow + "[+] Max Timeout: " + reset + strconv.Itoa(timeout) + "s")
	fmt.Println(lcyan + "[+] Delay Per Req: " + reset + strconv.Itoa(delay) + "ms")
	fmt.Println(lpurple + "[+] Threads: " + reset + strconv.Itoa(threads))

	if args.Output != "" {
		fmt.Println(lcyan + "[+] Output File: " + reset + args.Output)
		file, err := os.OpenFile(args.Output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			logger.Fatal(lred + "[ERR] " + err.Error() + reset)
		}
		file.Close()
	}

	if strings.Contains(url, "https://") {
		httpsFlag = true
	} else if strings.Contains(url, "http://") {
		httpsFlag = false
	} else {
		fmt.Println(lgreen + "[NOTE] Using DNS enumeration mode" + reset)
		rand.Seed(time.Now().UnixNano())
		wildcardDNSFlag, wildcardAddr = CheckWildcardDNS(url)
		dnsFlag = true

		if wildcardDNSFlag {
			fmt.Println(yellow + "[ALERT] Wildcard DNS is enabled on domain" + reset)
		} else {
			fmt.Println(lgreen + "[NOTE] No wildcard DNS detected on domain" + reset)
		}
	}

	var acceptedCodes []int

	if len(codes) > 0 {
		parts := strings.Split(codes, ",")	
		for _, part := range parts {
			part = strings.TrimSpace(part)
			code, _ := strconv.Atoi(part)
			acceptedCodes = append(acceptedCodes, code)
		}
		
		message := lgreen  + "[+] Allowed Codes:" + reset
		for _, code := range acceptedCodes {
			message += fmt.Sprintf(" %d", code)
		}
		fmt.Println(message)
	}

	semaphore := make(chan struct{}, threads)

	for _, subdomain := range wordlist {
		semaphore <- struct{}{}

		go func(subdomain string) {
			if dnsFlag != true {
				enumHttpSubdomain(url, subdomain, httpsFlag, acceptedCodes, delay, timeout, output)
			} else {
				enumDNSSubdomain(url, subdomain, delay, wildcardDNSFlag, wildcardAddr, output)
			}	
			<-semaphore
		}(subdomain)
	}

	for i := 0; i < cap(semaphore); i++ {
		semaphore <- struct{}{}
	}
}