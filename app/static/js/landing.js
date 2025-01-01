async function deleteResource(type, name) {
    if (!confirm(`Are you sure you want to delete sandbox ${name}?`)) {
        return;
    }

    try {
        const response = await fetch(`/api/resources/${type}/${name}?ns={{ namespace }}`, {
            method: 'DELETE'
        });
        
        if (!response.ok) {
            throw new Error('Failed to delete resource');
        }
        
        setTimeout(function(){ window. location. reload(); }, 1000);
        
    } catch (error) {
        console.error('Error deleting resource:', error);
        alert('Failed to delete resource. Please try again.');
    }
}

document.addEventListener('DOMContentLoaded', async () => {
    const namespacesList = document.getElementById('namespacesList');

    async function fetchNamespaces() {
        try {
            const response = await fetch('/api/namespaces');
            if (!response.ok) throw new Error('Failed to fetch namespaces');
            
            const data = await response.json();
            
            if (!data.namespaces || data.namespaces.length === 0) {
                namespacesList.innerHTML = '<p class="no-sandboxes">No sandboxes found</p>';
                return;
            }
            
            namespacesList.innerHTML = data.namespaces
                .map(ns => `
                    <div class="namespace-item">
                        <div onclick="handleNamespaceClick('${ns.name}')">
                        <span class="namespace-name">${ns.name}</span>
                        <!--  <span class="namespace-created">${ns.created}</span>    -->
                        </div>       

                            <button class="delete-btn" onclick="deleteResource('namespace', '${ ns.name }')">
                                <span>🗑️</span>
                            </button>
                    </div>       
                `).join('');
        } catch (error) {
            console.error('Error:', error);
            namespacesList.innerHTML = '<p class="error-message">Error loading sandboxes</p>';
        }
    }

    if (namespacesList) {
        await fetchNamespaces();
    } else {
        console.error('namespacesList element not found!');
    }
});

function handleNamespaceClick(namespace) {
    NamespaceManager.setActiveNamespace(namespace);
    window.location.href = `/dashboard?ns=${namespace}`;
}