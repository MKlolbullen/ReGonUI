document.addEventListener("DOMContentLoaded", function () {
    const projectNameElement = document.getElementById('project-name');
    const projectName = projectNameElement ? projectNameElement.textContent.trim() : '';

    if (!projectName) {
        console.error("Project name not found.");
        return;
    }

    // Initialize inputs list from existing files (optional enhancement)
    // You can implement fetching existing inputs from the server if needed
});

async function runTool() {
    const inputType = document.getElementById('inputType').value;
    const tool = document.getElementById('tool').value;
    const input = document.getElementById('input').value.trim();
    const projectName = document.getElementById('project-name').textContent.trim();

    if (input === "") {
        alert("Input cannot be empty.");
        return;
    }

    const response = await fetch('/run-tool', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            projectName: projectName,
            tool: tool,
            inputType: inputType,
            input: input
        })
    });

    const data = await response.json();

    if (response.ok) {
        document.getElementById('output').textContent = data.output;

        const inputsList = document.getElementById('inputs-list');
        const listItem = document.createElement('li');
        listItem.textContent = `${inputType.toUpperCase()}: ${input}`;
        listItem.setAttribute('data-tool', tool);
        listItem.setAttribute('data-input-type', inputType);
        listItem.setAttribute('data-input', input);
        listItem.onclick = viewPreviousOutput;
        inputsList.appendChild(listItem);
    } else {
        const errorOutput = data.error || "Failed to run tool.";
        document.getElementById('output').textContent = errorOutput;
    }
}

async function viewPreviousOutput(event) {
    const projectName = document.getElementById('project-name').textContent.trim();
    const tool = event.currentTarget.getAttribute('data-tool');
    const inputType = event.currentTarget.getAttribute('data-input-type');
    const input = event.currentTarget.getAttribute('data-input');

    const response = await fetch(`/output?project=${encodeURIComponent(projectName)}&tool=${encodeURIComponent(tool)}&inputType=${encodeURIComponent(inputType)}&input=${encodeURIComponent(input)}`);

    const data = await response.json();

    if (response.ok) {
        document.getElementById('output').textContent = data.output;
    } else {
        const errorOutput = data.error || "Failed to load previous output.";
        document.getElementById('output').textContent = errorOutput;
    }
}

async function showHelp() {
    const tool = document.getElementById('tool').value;

    const response = await fetch(`/tool-help?tool=${encodeURIComponent(tool)}`);

    const data = await response.json();

    if (response.ok) {
        alert(data.help);
    } else {
        alert("Failed to retrieve help information.");
    }
}
