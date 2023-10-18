# subby
An uber fast subdomain enumeration tool that automatically detects for wildcard DNS records and excludes invalid subdomains from results. Features enumeration modes with DNS or HTTP(S).

## Features
<img width="459" alt="image" src="https://github.com/n0mi1k/subby/assets/28621928/1ee5deba-85a7-4b1a-9158-53af17718b7e">

**Usage is simple, choose from the 2 enumeration modes below:**

**DNS Mode** is super fast, stealthy and utilises purely DNS requests which avoids hitting the infra and detects for wildcard DNS records.  
`./subby -u <domain> -w </path/to/wordlist>`

**Web Mode** is slower and noisier but it performs web requests and returns the corresponding status code, useful for identifying web applications.  
`./subby -u https://<domain> -w </path/to/wordlist>`

**NOTE:** *To run in web mode, include the scheme http:// or https:// in the domain*

## Installation
Subby requires Go 1.18 and above to install successfully. To install, just run the below command or download pre-compiled binary from the [release page](https://github.com/n0mi1k/subby/releases/).
```
go install github.com/n0mi1k/subby@latest
```

## Options
```console
USAGE:
  ./subby [flags]

FLAGS:
   -u, --url         Target domain to enumerate [Required]
   -w, --wordlist    Wordlist to use [Required]
   -d, --delay       Set delay in milliseconds for each request (Default 0ms)
   -r, --response    Display results with these status codes separated by commas (e.g 200,301)
   -t, --threads     Number of concurrent requester threads (Default 50)
   -s, --timeout     Maximum timeout in seconds for web requests (Default 2s)
   -o, --output      Output filename to save results
```

## Advance Usage Examples
DNS Enumeration (100 Threads, 200ms Delay, Output to results.txt):  
`./subby -u <domain> -w </path/to/wordlist> -t 100 -d 200 -o results.txt`

Web Enumeration (20 Threads, 200ms Delay, 5s Max Request Timeout, Show Codes 200 and 301, Output to results.txt):  
`./subby -u https://<domain> -w </path/to/wordlist> -t 20 -d 200 -s 5 -r "200,301" -o results.txt`

## Wildcard DNS Records
A wildcard DNS record answers DNS requests for any subdomain isn't defined. Some domains has this configured which makes subdomain enumeration tedious, as invalid subdomains still receives an answer. Subby automatically detects for wildcard DNS records and filters out false positives, accurately displaying valid and existing subdomains.

## Disclaimer
This tool is for educational and testing purposes only. Do not use it to exploit the vulnerability on any system that you do not own or have permission to test. The authors of this script are not responsible for any misuse or damage caused by its use.