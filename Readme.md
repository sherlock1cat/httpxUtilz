# httpxUtilz

httpxUtilz is a basic tool for target information gathering and attack surface.

## Dependencies

httpxUtilz relies on the following third-party packages:

- github.com/projectdiscovery/utils
- github.com/projectdiscovery/dnsx
- github.com/projectdiscovery/cdncheck
- github.com/projectdiscovery/asnmap

Make sure to install these dependencies before building and running httpxUtilz.

You can install the dependencies using the following command:

   ```
   go mod tidy
   ```

## Installation

1. Clone or download this repository.

   ```
   git clone https://github.com/sherlock1cat/httpxUtilz.git
   ```

2. Navigate to the project directory.

   ```
   cd httpxUtilz
   ```

3. Build the project using the following command:

   ```
   make
   ```

4. Run the generated executable file.

   ```
   # Determine the ./data directory in the current directory.
   ./cmd/httpxUtilz -h
   ```

## Usage

httpxUtilz supports the following command-line arguments:

- `-url`: Single URL to process.
- `-urls`: File containing a list of URLs to process.
- `-proxy`: Proxy URL.
- `-usehttps`: Initiate an HTTPS request (default: true).
- `-followredirects`: Perform URL request redirection (default: true).
- `-maxredirects`: Maximum number of redirections (default: 10).
- `-method`: Default request method (default: GET).
- `-randomuseragent`: Use a random User-Agent header (default: true).
- `-headers`: Customize the request headers.
- `-followsamehost`: Follow Same Host (default: true).
- `-processes`: Number of processes (default: 1).
- `-rateLimit`: Rate limit (default: 100).
- `-res`: Save the result (default: false).
- `-resultFile`: File to save the result (default: ./result.json).
- `-passive`: Default not get passive info data.
- `-mayvul`: Default not get may vul info data.

## Examples

### Process a Single URL
"Determine the ./data directory in the current directory"
```
./httpxUtilz -url=https://www.hackerone.com
```

### Process URL List File

```
./httpxUtilz -urls=urls.txt
```

### More Options

```
cat url.txt | httpx -slient | ./httpxUtilz -proxy=http://127.0.0.1:1080 -maxredirects=5 -method=POST -randomuseragent=true -processes=50 -rateLimit=100 -res=true -resultFile=./result.json
```

- search vul information by waybackurl

```
echo "hackerone.com" | waybackurls -no-subs | ./httpxUtilz -randomuseragent=true -processes=50 -rateLimit=100 -base=false -mayvul=true -res=true -resultFile=./mayvul_result.json
```

## As Library

```
go get -u https://github.com/sherlock1cat/httpxUtilz@latest
```

## Notes

- Make sure you have Go programming language environment installed.
- Ensure sufficient permissions to execute httpxUtilz.
- Customize the command-line arguments as per your requirements.

## Contributing

Feel free to raise issues, report bugs, or contribute code. Create an issue or submit a pull request on GitHub.

## License

This project is distributed under the MIT License. See the LICENSE file for more information.