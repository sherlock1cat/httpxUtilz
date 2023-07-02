# httpxUtilz

httpxUtilz is a basic tool for target information gathering.

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

## Examples

### Process a Single URL

```
./httpxUtilz -url https://www.hackerone.com
```

### Process URL List File

```
./httpxUtilz -urls urls.txt
```

### More Options

```
./httpxUtilz -url https://www.hackerone.com -proxy http://127.0.0.1:1080 -usehttps true -followredirects true -maxredirects 5 -method POST -randomuseragent true -headers "Authorization: Bearer token" -followsamehost true -processes 2 -rateLimit 100 -res true -resultFile result.json
```

## Notes

- Make sure you have Go programming language environment installed.
- Ensure sufficient permissions to execute httpxUtilz.
- Customize the command-line arguments as per your requirements.

## Contributing

Feel free to raise issues, report bugs, or contribute code. Create an issue or submit a pull request on GitHub.

## License

This project is distributed under the MIT License. See the LICENSE file for more information.