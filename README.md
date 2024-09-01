# ReGonUI
A multi purpose Golang tool that is infused into a web-based GUI with an Nmap visualizer built-in the web framework.

## 1. Create a directory for the tool.

```bash
mkdir recon-platform
cd recon-platform
# Initialize a new Go module.
go mod init recon-platform

# Get some dependencies:
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite # SQLite as an example DB
```

## 2. Make sure the project follows this structure:
   
recon-platform/
├── main.go
├── projects/             # This directory will store project data (created dynamically)
├── templates/
│   ├── index.html
│   └── project.html
└── static/
    ├── project.js
    ├── script.js
    └── style.css

## 3. install dependencies and the tools you wish to use.


## 4. Ensure Required Tools Are Installed:
If you have Go installed and configured (i.e. with $GOPATH/bin in your $PATH):
```bash
go get -u github.com/tomnomnom/assetfinder
go install -v github.com/projectdiscovery/httpx@latest
go install -v github.com/projectdiscovery/dnsx@latest
go install -v github.com/projectdiscovery/subfinder@latest
go install -v github.com/projectdiscovery/assetfinder@latest
```
Masscan comes with most pentesting distros but can also be installed like this: 
```bash
sudo apt-get --assume-yes install git make gcc
git clone https://github.com/robertdavidgraham/masscan
cd masscan
make
```
Make sure all the reconnaissance tools (`assetfinder`, `subfinder`, `masscan`, `dnsx`, `httpx`, etc.) are installed and accessible in your system's PATH.


## 5. Run the Application:

   Navigate to your project directory and execute:

   ```bash
   go run main.go
   ```

## 6. Access the Application:

   Open your browser and navigate to [http://localhost:8080](http://localhost:8080). You should see the **"Start New Project"** page.

## 7. Map out networks and visualize them
- 7.1  Run Nmap Scan:
Enter the target IP range and the interface IP, then click "Run Nmap."
The network map should render with nodes representing the discovered hosts.
- 7.2 Interact with the Network Map:
Click on a node to view the open ports and services running on that host in a pop-up window.
