async function runNmap() {
    const target = document.getElementById('target').value.trim();
    const ifaceIP = document.getElementById('ifaceIP').value.trim();

    if (target === "" || ifaceIP === "") {
        alert("Both target and interface IP are required.");
        return;
    }

    const formData = new FormData();
    formData.append("target", target);
    formData.append("ifaceIP", ifaceIP);

    const response = await fetch('/run-nmap', {
        method: 'POST',
        body: formData
    });

    const data = await response.json();

    if (response.ok) {
        renderNetworkMap(data.ifaceIP, data.hosts);
    } else {
        alert("Failed to run nmap: " + (data.error || "Unknown error"));
    }
}

function renderNetworkMap(ifaceIP, hosts) {
    const elements = [];

    // Add the root node (interface IP)
    elements.push({
        data: { id: ifaceIP, label: ifaceIP },
        position: { x: 300, y: 300 },
        style: { 'background-color': 'red', 'shape': 'square' }
    });

    // Add nodes and edges
    hosts.forEach((host, index) => {
        const ip = host.Addresses.find(addr => addr.Type === 'ipv4').Addr;

        // Add host node
        elements.push({
            data: { id: ip, label: ip }
        });

        // Add edge from interface IP to host
        elements.push({
            data: { id: 'edge' + index, source: ifaceIP, target: ip }
        });
    });

    const cy = cytoscape({
        container: document.getElementById('cy'),
        elements: elements,
        style: [
            {
                selector: 'node',
                style: {
                    'label': 'data(label)',
                    'text-valign': 'center',
                    'text-halign': 'center',
                    'background-color': '#3498db',
                    'color': '#fff',
                    'width': '40px',
                    'height': '40px',
                    'font-size': '12px'
                }
            },
            {
                selector: 'edge',
                style: {
                    'width': 2,
                    'line-color': '#000',
                    'curve-style': 'bezier'
                }
            }
        ],
        layout: {
            name: 'circle'
        }
    });

    // Add click event for showing host details
    cy.nodes().on('click', function (event) {
        const nodeId = event.target.id();
        const host = hosts.find(h => h.Addresses.some(addr => addr.Addr === nodeId));

        if (host) {
            showHostDetails(host);
        }
    });
}

function showHostDetails(host) {
    const ip = host.Addresses.find(addr => addr.Type === 'ipv4').Addr;
    let details = `Host: ${ip}\n\nOpen Ports:\n`;

    host.Ports.Ports.forEach(port => {
        details += `Port: ${port.PortID}, Service: ${port.Service.Name}, Version: ${port.Service.Product} ${port.Service.Version}\n`;
    });

    alert(details);
}
