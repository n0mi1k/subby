# subby
An uber fast subdomain enumeration tool that automatically detects for wildcard DNS records and excludes invalid subdomains from results. Features enumeration modes with DNS or HTTP(S).

**Featured on BHEH:** https://www.blackhatethicalhacking.com/tools/subby/

## Features
<img width="451" alt="image" src="https://github.com/n0mi1k/subby/assets/28621928/ad2834a4-bd98-4851-9df9-400ac26bcc80">

**Usage is simple, choose from the 2 enumeration modes below:**

**DNS Mode** is fast, stealthy and utilises purely DNS requests and auto detects for wildcard DNS records.  
`subby -u <domain> -w </path/to/wordlist>`

**Web Mode** performs concurrent web requests and returns the corresponding status code, useful for identifying web applications.  
`subby -u https://<domain> -w </path/to/wordlist>` (For HTTPS)  
`subby -u http://<domain> -w </path/to/wordlist>`  (For HTTP)  

**NOTE:** *Setting a delay is highly reccomended to prevent rate limiting or DoSing the DNS server*

## Installation
Subby requires Go 1.18 and above to install successfully. To install, just run the below command or download pre-compiled binary from the [release page](https://github.com/n0mi1k/subby/releases/).
```
go install github.com/n0mi1k/subby@latest
```
If your Go binaries are not added to PATH on Kali, do this:
```bash
nano ~/.zshrc
export PATH="$PATH:/home/kali/go/bin"
source ~/.zshrc
```

## Options
```console
USAGE:
  subby [flags]

FLAGS:
   -u, --url         Target domain to enumerate [Required]
   -w, --wordlist    Wordlist to use [Required]
   -d, --delay       Set delay in milliseconds for each request (Default 0ms)
   -r, --response    Only display results with these status codes separated by commas (e.g 200,301)
   -t, --threads     Number of concurrent requester threads (Default 50)
   -s, --timeout     Maximum timeout in seconds for web requests (Default 2s)
   -o, --output      Output filename to save results
```
Using `-d` to set a delay is highly recommended to avoid getting blocked or affecting your DNS queries. 

## Advance Usage Examples
DNS Enumeration (100 Threads, 200ms Delay, Output to results.txt):  
`subby -u <domain> -w </path/to/wordlist> -t 100 -d 200 -o results.txt`

Web Enumeration (20 Threads, 200ms Delay, 5s Max Request Timeout, Show Codes 200 and 301, Output to results.txt):  
`subby -u https://<domain> -w </path/to/wordlist> -t 20 -d 200 -s 5 -r "200,301" -o results.txt`

## Wildcard DNS Records
A wildcard DNS record answers DNS requests for any subdomain isn't defined. Some domains has this configured which makes subdomain enumeration tedious, as invalid subdomains still receives an answer. Subby automatically detects for wildcard DNS records and filters out false positives, accurately displaying valid and existing subdomains.

## Disclaimer
This tool is for educational and testing purposes only. Do not use it to exploit the vulnerability on any system that you do not own or have permission to test. The authors of this script are not responsible for any misuse or damage caused by its use.
