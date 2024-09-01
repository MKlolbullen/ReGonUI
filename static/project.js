async function startProject() {
    const projectName = document.getElementById('projectName').value.trim();

    if (projectName === "") {
        alert("Project name cannot be empty.");
        return;
    }

    const response = await fetch('/start-project', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ projectName: projectName })
    });

    const data = await response.json();

    if (response.ok) {
        // Redirect to the main application page with the project name as a query parameter
        window.location.href = '/?project=' + encodeURIComponent(data.projectName);
    } else {
        // Display error message
        const errorMessage = data.error || "Failed to start project.";
        document.getElementById('error-message').textContent = errorMessage;
    }
}
