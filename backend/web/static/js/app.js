function devGenesis() {
    return {
        installing: false,
        completed: false,
        progress: 0,
        status: 'Ready to create your development environment',
        error: null,
        ws: null,
        tools: [
            { name: 'Python 3.12', status: 'pending' },
            { name: 'Node.js LTS', status: 'pending' },
            { name: 'Git', status: 'pending' },
            { name: 'Visual Studio Code', status: 'pending' },
            { name: 'WSL2 + Ubuntu', status: 'pending' },
            { name: 'Claude Code', status: 'pending' }
        ],
        
        init() {
            // Create animated particles
            this.createParticles();
        },
        
        createParticles() {
            // Particles are created in the template
        },
        
        async startInstallation() {
            if (this.installing || this.completed) return;
            
            this.installing = true;
            this.error = null;
            this.progress = 0;
            this.status = 'Connecting to server...';
            
            try {
                await this.connectWebSocket();
            } catch (err) {
                this.error = 'Failed to connect to server: ' + err.message;
                this.installing = false;
            }
        },
        
        connectWebSocket() {
            return new Promise((resolve, reject) => {
                const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
                const wsUrl = `${protocol}//${window.location.host}/ws`;
                
                this.ws = new WebSocket(wsUrl);
                
                this.ws.onopen = () => {
                    console.log('WebSocket connected');
                    this.status = 'Connected! Starting installation...';
                    
                    // Send install command
                    setTimeout(() => {
                        this.ws.send(JSON.stringify({ type: 'install' }));
                    }, 1000);
                    
                    resolve();
                };
                
                this.ws.onmessage = (event) => {
                    const data = JSON.parse(event.data);
                    this.handleMessage(data);
                };
                
                this.ws.onerror = (err) => {
                    console.error('WebSocket error:', err);
                    reject(new Error('WebSocket connection failed'));
                };
                
                this.ws.onclose = () => {
                    console.log('WebSocket closed');
                    if (this.installing && !this.completed) {
                        this.error = 'Connection lost. Please refresh and try again.';
                        this.installing = false;
                    }
                };
            });
        },
        
        handleMessage(data) {
            switch (data.type) {
                case 'progress':
                    this.progress = data.percentage || 0;
                    this.status = data.message || 'Installing...';
                    this.updateToolStatus(data.message);
                    break;
                
                case 'status':
                    this.status = data.message;
                    break;
                
                case 'error':
                    this.error = data.error;
                    this.installing = false;
                    break;
                
                case 'complete':
                    this.completed = true;
                    this.installing = false;
                    this.progress = 100;
                    this.status = data.message || 'Installation complete!';
                    this.celebrateSuccess();
                    break;
            }
        },
        
        updateToolStatus(message) {
            // Update tool status based on message content
            this.tools.forEach(tool => {
                if (message.includes(tool.name)) {
                    if (message.includes('')) {
                        tool.status = 'completed';
                    } else if (message.includes('Installing')) {
                        tool.status = 'installing';
                    }
                }
            });
        },
        
        celebrateSuccess() {
            // Show success animation
            const successEl = document.querySelector('.success-animation');
            successEl.style.display = 'block';
            successEl.style.animation = 'celebrate 1s ease-out';
            
            // Create confetti effect
            for (let i = 0; i < 50; i++) {
                setTimeout(() => {
                    this.createConfetti();
                }, i * 30);
            }
        },
        
        createConfetti() {
            const confetti = document.createElement('div');
            confetti.style.position = 'fixed';
            confetti.style.width = '10px';
            confetti.style.height = '10px';
            confetti.style.background = `hsl(${Math.random() * 360}, 100%, 70%)`;
            confetti.style.left = Math.random() * 100 + '%';
            confetti.style.top = '-20px';
            confetti.style.transform = `rotate(${Math.random() * 360}deg)`;
            confetti.style.transition = 'all 2s ease-out';
            document.body.appendChild(confetti);
            
            setTimeout(() => {
                confetti.style.top = '100vh';
                confetti.style.transform = `rotate(${Math.random() * 720}deg) translateX(${(Math.random() - 0.5) * 200}px)`;
                confetti.style.opacity = '0';
            }, 10);
            
            setTimeout(() => {
                confetti.remove();
            }, 2000);
        }
    };
}