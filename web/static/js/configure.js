document.addEventListener('DOMContentLoaded', () => {
    const deploymentCheckboxes = document.querySelectorAll('input[type="checkbox"][name="deployments"]');
    const replicaInputs = document.querySelectorAll('input[type="number"]');
    const serviceInputs = document.querySelectorAll('.service-control input[type="checkbox"]');
    const createBtn = document.getElementById('createDeploymentBtn');


    deploymentCheckboxes.forEach((checkbox, index) => {
        checkbox.addEventListener('change', () => {
            serviceInputs[index].disabled = !checkbox.checked;
            replicaInputs[index].disabled = !checkbox.checked;
        });
    }); 


    createBtn.addEventListener('click', async () => {
        const selectedDeployments = [];
        const selectedServices = [];
        console.log('Selected Services:', selectedServices);
    
        deploymentCheckboxes.forEach((checkbox, index) => {
            if (checkbox.checked) {
                const deployment = {
                    type: checkbox.value,
                    replicas: parseInt(replicaInputs[index].value),
                    service: serviceInputs[index].checked
                };
    
                selectedDeployments.push(deployment);
    
                if (serviceInputs[index].checked) {
                    selectedServices.push({
                        type: checkbox.value,
                        replicas: parseInt(replicaInputs[index].value)
                    });
                }
            }
        });

        if (selectedDeployments.length === 0) {
            alert('Please select at least one deployment');
            return;
        }

        try {
            const urlParams = new URLSearchParams(window.location.search);
            const namespace = urlParams.get('ns');

            const response = await fetch('/api/deployments', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    namespace: namespace,
                    deployments: selectedDeployments
                })
            });

            if (!response.ok) {
                throw new Error('Failed to create deployments');
            }
            console.log('Selected Services:', selectedServices);
            if (selectedServices.length > 0) {
                const serviceResponse = await fetch('/api/services', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        namespace: namespace,
                        services: selectedServices
                    })
                });
    
                if (!serviceResponse.ok) {
                    throw new Error('Failed to create services');
                }
            }    

            window.location.href = `/dashboard?ns=${namespace}`;
            
        } catch (error) {
            console.error('Error:', error);
            alert('Failed to create deployments. Please try again.');
        }
    });
});
