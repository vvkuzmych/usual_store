import React, { useState, useEffect, useRef } from 'react';
import './ChatWidget.css';

const ChatWidget = ({ apiUrl = 'http://localhost:8080' }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [sessionId] = useState(`session-${Date.now()}`);
  const [isRecording, setIsRecording] = useState(false);
  const [useVoiceResponse, setUseVoiceResponse] = useState(false);
  const messagesEndRef = useRef(null);
  const mediaRecorderRef = useRef(null);
  const audioChunksRef = useRef([]);

  useEffect(() => {
    // Welcome message
    if (messages.length === 0) {
      setMessages([{
        id: 0,
        role: 'assistant',
        content: '👋 Hi! I\'m your shopping assistant. How can I help you find the perfect product today?',
        timestamp: new Date()
      }]);
    }
  }, []);

  useEffect(() => {
    // Auto-scroll to bottom
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const sendMessage = async () => {
    const message = inputMessage.trim();
    if (!message || isLoading) return;

    // Add user message
    const userMessage = {
      id: messages.length,
      role: 'user',
      content: message,
      timestamp: new Date()
    };
    setMessages(prev => [...prev, userMessage]);
    setInputMessage('');
    setIsLoading(true);

    try {
      const response = await fetch(`${apiUrl}/api/ai/chat`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          session_id: sessionId,
          message: message,
          user_id: null // TODO: Get from auth context
        })
      });

      if (!response.ok) {
        throw new Error('Failed to get response');
      }

      const data = await response.json();

      // Add assistant message
      const assistantMessage = {
        id: messages.length + 1,
        role: 'assistant',
        content: data.message,
        timestamp: new Date(),
        products: data.products,
        suggestions: data.suggestions
      };
      setMessages(prev => [...prev, assistantMessage]);

    } catch (error) {
      console.error('Error sending message:', error);
      const errorMessage = {
        id: messages.length + 1,
        role: 'assistant',
        content: 'Sorry, I\'m having trouble connecting. Please try again in a moment.',
        timestamp: new Date()
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const mediaRecorder = new MediaRecorder(stream, {
        mimeType: 'audio/webm;codecs=opus'
      });
      
      mediaRecorderRef.current = mediaRecorder;
      audioChunksRef.current = [];

      mediaRecorder.ondataavailable = (event) => {
        if (event.data.size > 0) {
          audioChunksRef.current.push(event.data);
        }
      };

      mediaRecorder.onstop = async () => {
        const audioBlob = new Blob(audioChunksRef.current, { type: 'audio/webm' });
        await sendVoiceMessage(audioBlob);
        stream.getTracks().forEach(track => track.stop());
      };

      mediaRecorder.start();
      setIsRecording(true);
    } catch (error) {
      console.error('Error accessing microphone:', error);
      alert('Unable to access microphone. Please check your permissions.');
    }
  };

  const stopRecording = () => {
    if (mediaRecorderRef.current && isRecording) {
      mediaRecorderRef.current.stop();
      setIsRecording(false);
    }
  };

  const sendVoiceMessage = async (audioBlob) => {
    if (isLoading) return;

    setIsLoading(true);

    // Show user message indicating voice input
    const userMessage = {
      id: messages.length,
      role: 'user',
      content: '🎤 [Voice message]',
      timestamp: new Date(),
      isVoice: true
    };
    setMessages(prev => [...prev, userMessage]);

    try {
      const formData = new FormData();
      formData.append('audio', audioBlob, 'recording.webm');
      formData.append('session_id', sessionId);
      formData.append('return_audio', useVoiceResponse.toString());
      formData.append('voice', 'alloy'); // Default voice

      const response = await fetch(`${apiUrl}/api/ai/voice`, {
        method: 'POST',
        body: formData
      });

      if (!response.ok) {
        throw new Error('Failed to process voice message');
      }

      if (useVoiceResponse && response.headers.get('content-type')?.includes('audio')) {
        // Handle audio response
        const audioBlob = await response.blob();
        const audioUrl = URL.createObjectURL(audioBlob);
        const audio = new Audio(audioUrl);
        
        const assistantMessage = {
          id: messages.length + 1,
          role: 'assistant',
          content: '🔊 [Audio response]',
          timestamp: new Date(),
          audioUrl: audioUrl,
          isVoice: true
        };
        setMessages(prev => [...prev, assistantMessage]);
        
        // Auto-play audio response
        audio.play();
      } else {
        // Handle text response
        const data = await response.json();
        
        const assistantMessage = {
          id: messages.length + 1,
          role: 'assistant',
          content: data.message || data.transcription,
          timestamp: new Date(),
          products: data.products,
          suggestions: data.suggestions
        };
        setMessages(prev => [...prev, assistantMessage]);
      }

    } catch (error) {
      console.error('Error sending voice message:', error);
      const errorMessage = {
        id: messages.length + 1,
        role: 'assistant',
        content: 'Sorry, I had trouble processing your voice message. Please try again.',
        timestamp: new Date()
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  const playAudio = (audioUrl) => {
    const audio = new Audio(audioUrl);
    audio.play();
  };

  const submitFeedback = async (messageId, helpful) => {
    try {
      await fetch(`${apiUrl}/api/ai/feedback`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          message_id: messageId,
          conversation_id: 1,
          helpful: helpful,
          feedback_type: helpful ? 'helpful' : 'not_helpful'
        })
      });
    } catch (error) {
      console.error('Error submitting feedback:', error);
    }
  };

  const formatTime = (date) => {
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    return `${hours}:${minutes}`;
  };

  return (
    <>
      {/* Chat Toggle Button */}
      {!isOpen && (
        <button className="chat-toggle" onClick={() => setIsOpen(true)}>
          <span className="icon">💬</span> Need Help?
        </button>
      )}

      {/* Chat Widget */}
      {isOpen && (
        <div className="chat-widget">
          {/* Header */}
          <div className="chat-header">
            <div className="chat-header-content">
              <h3>🤖 Shopping Assistant</h3>
              <p>We're here to help!</p>
            </div>
            <button className="chat-close-btn" onClick={() => setIsOpen(false)}>
              ×
            </button>
          </div>

          {/* Messages */}
          <div className="chat-messages">
            {messages.map((msg) => (
              <div key={msg.id} className={`message ${msg.role}-message`}>
                <div className="message-content">
                  {msg.content}
                  {msg.audioUrl && (
                    <button 
                      className="play-audio-btn"
                      onClick={() => playAudio(msg.audioUrl)}
                      title="Play audio"
                    >
                      🔊 Play
                    </button>
                  )}
                </div>
                <div className="message-time">
                  {formatTime(msg.timestamp)}
                </div>

                {/* Product Cards */}
                {msg.products && msg.products.length > 0 && (
                  <div className="product-recommendations">
                    {msg.products.map((product, idx) => (
                      <div key={idx} className="product-card">
                        {product.image && <img src={product.image} alt={product.name} />}
                        <h4>{product.name}</h4>
                        <p className="product-price">${product.price}</p>
                        <p className="product-reason">{product.reason}</p>
                        <button onClick={() => window.location.href = `/product/${product.id}`}>
                          View Product
                        </button>
                      </div>
                    ))}
                  </div>
                )}

                {/* Suggestions */}
                {msg.suggestions && msg.suggestions.length > 0 && (
                  <div className="suggestions">
                    {msg.suggestions.map((suggestion, idx) => (
                      <button 
                        key={idx}
                        className="suggestion-chip"
                        onClick={() => setInputMessage(suggestion)}
                      >
                        {suggestion}
                      </button>
                    ))}
                  </div>
                )}

                {/* Feedback */}
                {msg.role === 'assistant' && msg.id > 0 && (
                  <div className="message-feedback">
                    <button onClick={() => submitFeedback(msg.id, true)} title="Helpful">
                      👍
                    </button>
                    <button onClick={() => submitFeedback(msg.id, false)} title="Not helpful">
                      👎
                    </button>
                  </div>
                )}
              </div>
            ))}

            {/* Loading indicator */}
            {isLoading && (
              <div className="message assistant-message">
                <div className="typing-indicator">
                  <span></span><span></span><span></span>
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          {/* Input */}
          <div className="chat-input">
            <input
              type="text"
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Ask me anything..."
              disabled={isLoading || isRecording}
            />
            <button
              className={`voice-btn ${isRecording ? 'recording' : ''}`}
              onClick={isRecording ? stopRecording : startRecording}
              disabled={isLoading}
              title={isRecording ? 'Stop recording' : 'Record voice message'}
            >
              {isRecording ? '⏹️' : '🎤'}
            </button>
            <button 
              className="chat-send-btn"
              onClick={sendMessage} 
              disabled={isLoading || !inputMessage.trim() || isRecording}
            >
              Send
            </button>
          </div>
          {/* Voice response toggle */}
          <div className="voice-options">
            <label>
              <input
                type="checkbox"
                checked={useVoiceResponse}
                onChange={(e) => setUseVoiceResponse(e.target.checked)}
                disabled={isLoading || isRecording}
              />
              <span>Get audio response</span>
            </label>
          </div>
        </div>
      )}
    </>
  );
};

export default ChatWidget;

