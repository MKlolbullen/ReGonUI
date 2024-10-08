package main

import (
    "bytes"
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "sync"

    "github.com/gin-gonic/gin"
)

type ToolRequest struct {
    ProjectName string `json:"projectName"`
    Tool        string `json:"tool"`
    InputType   string `json:"inputType"`
    Input       string `json:"input"`
}

var projectMutex sync.Mutex

// Structs for parsing Nmap XML output
type NmapRun struct {
    XMLName xml.Name `xml:"nmaprun"`
    Hosts   []Host   `xml:"host"`
}

type Host struct {
    Addresses []Address `xml:"address"`
    Ports     Ports     `xml:"ports"`
}

type Address struct {
    Addr string `xml:"addr,attr"`
    Type string `xml:"addrtype,attr"`
}

type Ports struct {
    Ports []Port `xml:"port"`
}

type Port struct {
    Protocol string  `xml:"protocol,attr"`
    PortID   string  `xml:"portid,attr"`
    State    State   `xml:"state"`
    Service  Service `xml:"service"`
}

type State struct {
    State string `xml:"state,attr"`
}

type Service struct {
    Name    string `xml:"name,attr"`
    Product string `xml:"product,attr"`
    Version string `xml:"version,attr"`
}

func main() {
    router := gin.Default()

    // Serve static files (CSS, JS)
    router.Static("/static", "./static")

    // Load HTML templates
    router.LoadHTMLGlob("templates/*")

    // Define routes
    router.GET("/", func(c *gin.Context) {
        projectName := c.Query("project")
        if projectName == "" {
            c.HTML(http.StatusOK, "project.html", nil)
        } else {
            projectDir := filepath.Join("projects", projectName)
            if _, err := os.Stat(projectDir); os.IsNotExist(err) {
                c.HTML(http.StatusOK, "project.html", gin.H{
                    "error": "Project does not exist. Please create a new project.",
                })
                return
            }
            c.HTML(http.StatusOK, "index.html", gin.H{
                "projectName": projectName,
            })
        }
    })

    router.POST("/start-project", startProjectHandler)
    router.POST("/run-tool", runToolHandler)
    router.GET("/output", getOutputHandler)
    router.GET("/tool-help", toolHelpHandler)

    // Network map route
    router.POST("/run-nmap", runNmapHandler)

    // Start server
    router.Run(":8080")
}

func startProjectHandler(c *gin.Context) {
    var req ToolRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    projectDir := filepath.Join("projects", req.ProjectName)
    projectMutex.Lock()
    defer projectMutex.Unlock()
    if _, err := os.Stat(projectDir); os.IsNotExist(err) {
        if err := os.MkdirAll(projectDir, os.ModePerm); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project directory"})
            return
        }
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Project already exists"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":     "Project created",
        "projectName": req.ProjectName,
    })
}

func runToolHandler(c *gin.Context) {
    var req ToolRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    output := runTool(req.Tool, req.InputType, req.Input)

    projectDir := filepath.Join("projects", req.ProjectName)
    saveOutput(projectDir, req.InputType, req.Tool, req.Input, output)

    c.JSON(http.StatusOK, gin.H{
        "output": output,
    })
}

func getOutputHandler(c *gin.Context) {
    projectName := c.Query("project")
    inputType := c.Query("inputType")
    tool := c.Query("tool")
    input := c.Query("input")

    projectDir := filepath.Join("projects", projectName)
    sanitizedInput := sanitizeFileName(input)
    filename := fmt.Sprintf("%s_%s_%s.txt", tool, inputType, sanitizedInput)
    filePath := filepath.Join(projectDir, filename)

    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read output file"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "output": string(data),
    })
}

func toolHelpHandler(c *gin.Context) {
    tool := c.Query("tool")

    helpText := getToolHelp(tool)

    c.JSON(http.StatusOK, gin.H{
        "help": helpText,
    })
}

func getToolHelp(tool string) string {
    switch tool {
    case "assetfinder":
        return "Assetfinder helps you find domains related to a given domain. Use '--subs-only' to get subdomains only."
    case "subfinder":
        return "Subfinder is a subdomain discovery tool. Use '-d' followed by the domain name."
    case "masscan":
        return "Masscan is a fast port scanner. Use '-p1-65535' to scan all ports."
    case "dnsx":
        return "DNSx allows for DNS resolution. Use '-d' followed by the domain name."
    case "httpx":
        return "HTTPx probes for working HTTP servers. Use '-l' followed by a list of URLs."
    default:
        return "Help not available for this tool."
    }
}

func runTool(tool, inputType, input string) string {
    var cmd *exec.Cmd

    switch tool {
    case "assetfinder":
        if inputType != "domain" {
            return "Assetfinder requires a domain as input."
        }
        cmd = exec.Command("assetfinder", "--subs-only", input)
    case "subfinder":
        if inputType != "domain" {
            return "Subfinder requires a domain as input."
        }
        cmd = exec.Command("subfinder", "-d", input)
    case "masscan":
        if inputType != "ip" && inputType != "cidr" {
            return "Masscan requires an IP address or CIDR range as input."
        }
        cmd = exec.Command("masscan", "-p1-65535", input)
    case "dnsx":
        if inputType != "domain" && inputType != "ip" {
            return "DNSx requires a domain or IP address as input."
        }
        cmd = exec.Command("dnsx", "-d", input)
    case "httpx":
        if inputType != "domain" && inputType != "ip" {
            return "HTTPx requires a domain or IP address as input."
        }
        cmd = exec.Command("httpx", "-l", input)
    default:
        return "Tool not found"
    }

    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return "Error: " + err.Error()
    }

    return out.String()
}

func saveOutput(projectDir, inputType, tool, input, output string) {
    input = sanitizeFileName(input)
    filename := fmt.Sprintf("%s_%s_%s.txt", tool, inputType, input)
    filePath := filepath.Join(projectDir, filename)

    projectMutex.Lock()
    defer projectMutex.Unlock()

    f, err := os.Create(filePath)
    if err != nil {
        fmt.Println("Error creating output file:", err)
        return
    }
    defer f.Close()

    f.WriteString(output)
}

func sanitizeFileName(name string) string {
    return strings.ReplaceAll(name, "/", "_")
}

func runNmapHandler(c *gin.Context) {
    target := c.PostForm("target")
    ifaceIP := c.PostForm("ifaceIP")

    if target == "" || ifaceIP == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Target and Interface IP are required"})
        return
    }

    // Run nmap scan
    cmd := exec.Command("nmap", "-sV", "-oX", "-", target)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run nmap"})
        return
    }

    var nmapRun NmapRun
    err = xml.Unmarshal(out.Bytes(), &nmapRun)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse nmap output"})
        return
    }

    // Prepare hosts data
    hostsData := []map[string]interface{}{}
    for _, host := range nmapRun.Hosts {
        ip := ""
        for _, addr := range host.Addresses {
            if addr.Type == "ipv4" {
                ip = addr.Addr
                break
            }
        }
        if ip == "" {
            continue
        }

        ports := []map[string]string{}
        for _, port := range host.Ports.Ports {
            if port.State.State == "open" {
                ports = append(ports, map[string]string{
                    "port":     port.PortID,
                    "service":  port.Service.Name,
                    "product":  port.Service.Product,
                    "version":  port.Service.Version,
                    "protocol": port.Protocol,
                })
            }
        }

        hostsData = append(hostsData, map[string]interface{}{
            "ip":    ip,
            "ports": ports,
        })
    }

    // Return nmap results as JSON
    c.JSON(http.StatusOK, gin.H{
        "ifaceIP": ifaceIP,
        "hosts":   hostsData,
    })
}
