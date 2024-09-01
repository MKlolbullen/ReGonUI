# ReGonUI
A multi purpose Golang tool that is infused into a web-based GUI.

## 1. Create a directory for the tool.

```bash
mkdir recon-platform
cd recon-platform
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

```bash
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite # SQLite as an example DB
```
## 4. Ensure Required Tools Are Installed:
```bash
go install -v github.com/projectdiscovery/httpx@latest
go install -v github.com/projectdiscovery/dnsx@latest
go install -v github.com/projectdiscovery/subfinder@latest
go install -v github.com/projectdiscovery/assetfinder@latest
```

   Make sure all the reconnaissance tools (`assetfinder`, `subfinder`, `masscan`, `dnsx`, `httpx`, etc.) are installed and accessible in your system's PATH.


## 5. Run the Application:

   Navigate to your project directory and execute:

   ```bash
   go run main.go
   ```

## 6. Access the Application:

   Open your browser and navigate to [http://localhost:8080](http://localhost:8080). You should see the **"Start New Project"** page.
