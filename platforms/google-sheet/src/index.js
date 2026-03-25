document.getElementById('requestForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const httpMethod = document.getElementById('httpMethod').value;
            const url = document.getElementById('urlInput').value;
            const statusDiv = document.getElementById('statusMessage');
            const runButton = document.getElementById('runButton');

            // Validation
            if (!url.trim()) {
                statusDiv.textContent = 'Please enter a valid URL';
                statusDiv.className = 'status error';
                return;
            }

            // Show loading state
            statusDiv.textContent = 'Sending request...';
            statusDiv.className = 'status loading';
            runButton.disabled = true;

            try {
            // Call the Apps Script backend instead
            const response = await new Promise((resolve, reject) => {
                google.script.run
                    .withSuccessHandler(resolve)
                    .withFailureHandler(reject)
                    .makeHttpRequest(httpMethod, url);
            });

            if (response.success) {
                const data = response.data;
                statusDiv.textContent = `✓ Request successful (${response.status})`;
                statusDiv.className = 'status success';
                
                // Write to sheet
                google.script.run.writeResponseToSheet(data);
                
                console.log('Response:', data);
            } else {
                statusDiv.textContent = `✗ Request failed: ${response.error || response.status}`;
                statusDiv.className = 'status error';
            }
            } catch (error) {
                statusDiv.textContent = `✗ Error: ${error.message}`;
                statusDiv.className = 'status error';
            } finally {
                runButton.disabled = false;
            }
        });