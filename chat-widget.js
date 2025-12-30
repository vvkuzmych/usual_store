/**
 * AI Chat Widget
 * Standalone JavaScript chat widget for e-commerce sites
 */

class ChatWidget {
    constructor(config = {}) {
        this.config = {
            apiUrl: config.apiUrl || 'http://localhost:8080',
            position: config.position || 'bottom-right',
            theme: config.theme || 'blue',
            welcomeMessage: config.welcomeMessage || "👋 Hi! I'm your shopping assistant. How can I help you find the perfect product today?"
        };

        this.sessionId = this.generateSessionId();
        this.messages = [];
        this.isOpen = false;
        this.isLoading = false;
        this.isRecording = false;
        this.useVoiceResponse = false;
        this.mediaRecorder = null;
        this.audioChunks = [];

        this.init();
    }

    generateSessionId() {
        return 'session-' + Date.now() + '-' + Math.random().toString(36).substr(2, 9);
    }

    init() {
        this.createElements();
        this.attachEventListeners();
        this.addWelcomeMessage();
    }

    createElements() {
        const container = document.getElementById('chat-widget-container') || document.body;

        // Create toggle button
        this.toggleBtn = document.createElement('button');
        this.toggleBtn.className = 'chat-toggle';
        this.toggleBtn.innerHTML = '<span class="icon">💬</span> Need Help?';
        container.appendChild(this.toggleBtn);

        // Create chat widget
        this.widget = document.createElement('div');
        this.widget.className = 'chat-widget hidden';
        this.widget.innerHTML = `
            <div class="chat-header">
                <div class="chat-header-content">
                    <h3>🤖 Shopping Assistant</h3>
                    <p>We're here to help!</p>
                </div>
                <button class="chat-close-btn">×</button>
            </div>
            <div class="chat-messages" id="chat-messages"></div>
            <div class="chat-input">
                <input type="text" id="chat-input-field" placeholder="Ask me anything..." />
                <button class="voice-btn" id="voice-btn" title="Record voice message">🎤</button>
                <button class="chat-send-btn" id="chat-send-btn">Send</button>
            </div>
            <div class="voice-options">
                <label>
                    <input type="checkbox" id="voice-response-toggle" />
                    <span>Get audio response</span>
                </label>
            </div>
        `;
        container.appendChild(this.widget);

        // Get references
        this.messagesContainer = this.widget.querySelector('#chat-messages');
        this.inputField = this.widget.querySelector('#chat-input-field');
        this.sendBtn = this.widget.querySelector('#chat-send-btn');
        this.voiceBtn = this.widget.querySelector('#voice-btn');
        this.voiceToggle = this.widget.querySelector('#voice-response-toggle');
        this.closeBtn = this.widget.querySelector('.chat-close-btn');
    }

    attachEventListeners() {
        this.toggleBtn.addEventListener('click', () => this.open());
        this.closeBtn.addEventListener('click', () => this.close());
        this.sendBtn.addEventListener('click', () => this.sendMessage());
        this.voiceBtn.addEventListener('click', () => {
            if (this.isRecording) {
                this.stopRecording();
            } else {
                this.startRecording();
            }
        });
        this.voiceToggle.addEventListener('change', (e) => {
            this.useVoiceResponse = e.target.checked;
        });
        this.inputField.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage();
            }
        });
    }

    open() {
        this.isOpen = true;
        this.widget.classList.remove('hidden');
        this.toggleBtn.classList.add('hidden');
        this.inputField.focus();
    }

    close() {
        this.isOpen = false;
        this.widget.classList.add('hidden');
        this.toggleBtn.classList.remove('hidden');
    }

    addWelcomeMessage() {
        this.addMessage({
            role: 'assistant',
            content: this.config.welcomeMessage,
            timestamp: new Date()
        });
    }

    addMessage(message) {
        this.messages.push(message);
        this.renderMessage(message);
        this.scrollToBottom();
    }

    renderMessage(message) {
        const messageDiv = document.createElement('div');
        messageDiv.className = `message ${message.role}-message`;

        const time = this.formatTime(message.timestamp);

        let html = `
            <div class="message-content">
                ${this.escapeHtml(message.content)}
                ${message.audioUrl ? `<button class="play-audio-btn" data-audio-url="${message.audioUrl}" title="Play audio">🔊 Play</button>` : ''}
            </div>
            <div class="message-time">${time}</div>
        `;

        // Add product recommendations if available
        if (message.products && message.products.length > 0) {
            html += '<div class="product-recommendations">';
            message.products.forEach(product => {
                html += `
                    <div class="product-card">
                        ${product.image ? `<img src="${product.image}" alt="${product.name}">` : ''}
                        <h4>${this.escapeHtml(product.name)}</h4>
                        <div class="product-price">$${product.price.toFixed(2)}</div>
                        ${product.reason ? `<p class="product-reason">${this.escapeHtml(product.reason)}</p>` : ''}
                        <button onclick="window.location.href='/product/${product.id}'">View Product</button>
                    </div>
                `;
            });
            html += '</div>';
        }

        // Add suggestions if available
        if (message.suggestions && message.suggestions.length > 0) {
            html += '<div class="suggestions">';
            message.suggestions.forEach((suggestion, idx) => {
                html += `<button class="suggestion-chip" data-suggestion="${this.escapeHtml(suggestion)}">${this.escapeHtml(suggestion)}</button>`;
            });
            html += '</div>';
        }

        // Add feedback buttons for assistant messages
        if (message.role === 'assistant' && message.id) {
            html += `
                <div class="message-feedback">
                    <button data-message-id="${message.id}" data-helpful="true" title="Helpful">👍</button>
                    <button data-message-id="${message.id}" data-helpful="false" title="Not helpful">👎</button>
                </div>
            `;
        }

        messageDiv.innerHTML = html;

        // Attach suggestion click handlers
        messageDiv.querySelectorAll('.suggestion-chip').forEach(btn => {
            btn.addEventListener('click', (e) => {
                this.inputField.value = e.target.dataset.suggestion;
                this.inputField.focus();
            });
        });

        // Attach feedback handlers
        messageDiv.querySelectorAll('.message-feedback button').forEach(btn => {
            btn.addEventListener('click', (e) => {
                this.submitFeedback(
                    parseInt(e.target.dataset.messageId),
                    e.target.dataset.helpful === 'true'
                );
                e.target.classList.add('active');
                // Disable both buttons after feedback
                e.target.parentElement.querySelectorAll('button').forEach(b => b.disabled = true);
            });
        });

        // Attach play audio handlers
        messageDiv.querySelectorAll('.play-audio-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const audioUrl = e.target.dataset.audioUrl;
                if (audioUrl) {
                    this.playAudio(audioUrl);
                }
            });
        });

        this.messagesContainer.appendChild(messageDiv);
    }

    async sendMessage() {
        const message = this.inputField.value.trim();
        if (!message || this.isLoading) return;

        // Add user message
        this.addMessage({
            role: 'user',
            content: message,
            timestamp: new Date()
        });

        // Clear input
        this.inputField.value = '';
        this.setLoading(true);

        try {
            const response = await fetch(`${this.config.apiUrl}/api/ai/chat`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    session_id: this.sessionId,
                    message: message,
                    user_id: this.getUserId()
                })
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();

            // Add assistant response
            this.addMessage({
                id: Date.now(), // Temporary ID for feedback
                role: 'assistant',
                content: data.message,
                timestamp: new Date(),
                products: data.products || [],
                suggestions: data.suggestions || []
            });

        } catch (error) {
            console.error('Error sending message:', error);
            this.addMessage({
                role: 'assistant',
                content: '❌ Sorry, I\'m having trouble connecting. Please try again in a moment.',
                timestamp: new Date()
            });
        } finally {
            this.setLoading(false);
        }
    }

    async startRecording() {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
            this.mediaRecorder = new MediaRecorder(stream, {
                mimeType: 'audio/webm;codecs=opus'
            });
            
            this.audioChunks = [];

            this.mediaRecorder.ondataavailable = (event) => {
                if (event.data.size > 0) {
                    this.audioChunks.push(event.data);
                }
            };

            this.mediaRecorder.onstop = async () => {
                const audioBlob = new Blob(this.audioChunks, { type: 'audio/webm' });
                await this.sendVoiceMessage(audioBlob);
                stream.getTracks().forEach(track => track.stop());
            };

            this.mediaRecorder.start();
            this.isRecording = true;
            this.voiceBtn.textContent = '⏹️';
            this.voiceBtn.classList.add('recording');
            this.voiceBtn.title = 'Stop recording';
            this.inputField.disabled = true;
        } catch (error) {
            console.error('Error accessing microphone:', error);
            alert('Unable to access microphone. Please check your permissions.');
        }
    }

    stopRecording() {
        if (this.mediaRecorder && this.isRecording) {
            this.mediaRecorder.stop();
            this.isRecording = false;
            this.voiceBtn.textContent = '🎤';
            this.voiceBtn.classList.remove('recording');
            this.voiceBtn.title = 'Record voice message';
            this.inputField.disabled = false;
        }
    }

    async sendVoiceMessage(audioBlob) {
        if (this.isLoading) return;

        this.setLoading(true);

        // Show user message indicating voice input
        this.addMessage({
            role: 'user',
            content: '🎤 [Voice message]',
            timestamp: new Date(),
            isVoice: true
        });

        try {
            const formData = new FormData();
            formData.append('audio', audioBlob, 'recording.webm');
            formData.append('session_id', this.sessionId);
            formData.append('return_audio', this.useVoiceResponse.toString());
            formData.append('voice', 'alloy');

            const response = await fetch(`${this.config.apiUrl}/api/ai/voice`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                throw new Error('Failed to process voice message');
            }

            if (this.useVoiceResponse && response.headers.get('content-type')?.includes('audio')) {
                // Handle audio response
                const audioBlob = await response.blob();
                const audioUrl = URL.createObjectURL(audioBlob);
                const audio = new Audio(audioUrl);
                
                this.addMessage({
                    role: 'assistant',
                    content: '🔊 [Audio response]',
                    timestamp: new Date(),
                    audioUrl: audioUrl,
                    isVoice: true
                });
                
                // Auto-play audio response
                audio.play();
            } else {
                // Handle text response
                const data = await response.json();
                
                this.addMessage({
                    id: Date.now(),
                    role: 'assistant',
                    content: data.message || data.transcription,
                    timestamp: new Date(),
                    products: data.products || [],
                    suggestions: data.suggestions || []
                });
            }

        } catch (error) {
            console.error('Error sending voice message:', error);
            this.addMessage({
                role: 'assistant',
                content: '❌ Sorry, I had trouble processing your voice message. Please try again.',
                timestamp: new Date()
            });
        } finally {
            this.setLoading(false);
        }
    }

    playAudio(audioUrl) {
        const audio = new Audio(audioUrl);
        audio.play();
    }

    setLoading(loading) {
        this.isLoading = loading;
        this.sendBtn.disabled = loading;
        this.inputField.disabled = loading || this.isRecording;
        this.voiceBtn.disabled = loading;

        if (loading) {
            // Add typing indicator
            const typingDiv = document.createElement('div');
            typingDiv.className = 'message assistant-message';
            typingDiv.id = 'typing-indicator';
            typingDiv.innerHTML = `
                <div class="typing-indicator">
                    <span></span>
                    <span></span>
                    <span></span>
                </div>
            `;
            this.messagesContainer.appendChild(typingDiv);
            this.scrollToBottom();
        } else {
            // Remove typing indicator
            const typingIndicator = document.getElementById('typing-indicator');
            if (typingIndicator) {
                typingIndicator.remove();
            }
        }
    }

    async submitFeedback(messageId, helpful) {
        try {
            await fetch(`${this.config.apiUrl}/api/ai/feedback`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    message_id: messageId,
                    conversation_id: 1, // TODO: Track conversation ID
                    helpful: helpful,
                    rating: helpful ? 5 : 1,
                    feedback_type: helpful ? 'helpful' : 'not_helpful'
                })
            });
        } catch (error) {
            console.error('Error submitting feedback:', error);
        }
    }

    getUserId() {
        // TODO: Get from authentication system
        return null;
    }

    scrollToBottom() {
        setTimeout(() => {
            this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
        }, 100);
    }

    formatTime(date) {
        const hours = date.getHours().toString().padStart(2, '0');
        const minutes = date.getMinutes().toString().padStart(2, '0');
        return `${hours}:${minutes}`;
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// Auto-initialize if config exists
if (typeof window !== 'undefined') {
    window.ChatWidget = ChatWidget;
}

