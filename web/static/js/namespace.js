const NamespaceManager = {
    getActiveNamespace: function() {
        return localStorage.getItem('activeNamespace');
    },
    
    setActiveNamespace: function(namespace) {
        localStorage.setItem('activeNamespace', namespace);
        this.updateUI();
    },

    endSession: function() {
        localStorage.removeItem('activeNamespace');
        window.location.href = '/';
    },
    
    updateUI: function() {
        const sessionContainer = document.querySelector('.session-container');
        const activeWorkspace = document.querySelector('.active-workspace');
        const currentPath = window.location.pathname;
        const activeNamespace = this.getActiveNamespace();
        
        if (sessionContainer && activeWorkspace) {
            if (currentPath === '/') {
                sessionContainer.style.display = 'none';
            } else if (activeNamespace) {
                sessionContainer.style.display = 'flex';
                activeWorkspace.textContent = `Active Sandbox: ${activeNamespace}`;
            } else {
                sessionContainer.style.display = 'none';
            }
        }
    },
    
    redirectToCorrectUrl: function() {
        const currentPath = window.location.pathname;
        const namespace = this.getActiveNamespace();
        
        if (namespace) {
            if (currentPath === '/') {
                window.location.href = `/dashboard?ns=${namespace}`;
                return;
            }
            
            const urlParams = new URLSearchParams(window.location.search);
            const urlNamespace = urlParams.get('ns');
            
            if (urlNamespace !== namespace) {
                window.location.href = `${currentPath}?ns=${namespace}`;
                return;
            }
        } else if (currentPath !== '/') {
            window.location.href = '/';
            return;
        }
        
        this.updateUI();
    }
};

document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('namespaceForm');
    if (form) {
        form.addEventListener('submit', function(e) {
            e.preventDefault();
            
            const input = document.getElementById('namespaceInput');
            const namespace = input.value.trim();
            
            if (namespace) {
                fetch('/api/namespace', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ namespace: namespace })
                })
                .then(response => response.json())
                .then(data => {
                    if (data.error) {
                        throw new Error(data.error);
                    }
                    NamespaceManager.setActiveNamespace(namespace);
                    window.location.href = `/configure?ns=${namespace}`;
                })
                .catch(error => {
                    alert('Error: ' + error.message);
                });
            }
        });
    }
    
    NamespaceManager.redirectToCorrectUrl();
});