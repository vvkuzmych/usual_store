import React, { useState, useEffect, useRef } from 'react';
import {
  Box,
  Button,
  TextField,
  Typography,
  Paper,
  IconButton,
  Badge,
  Avatar,
  Divider,
  CircularProgress,
  Chip
} from '@mui/material';
import {
  Chat as ChatIcon,
  Close as CloseIcon,
  Send as SendIcon,
  SupportAgent as SupportAgentIcon,
  Minimize as MinimizeIcon
} from '@mui/icons-material';
import axios from 'axios';

const API_URL = process.env.REACT_APP_SUPPORT_API_URL || 'http://localhost:5001';
const WS_URL = process.env.REACT_APP_SUPPORT_WS_URL || 'ws://localhost:5001';

function SupportChatWidget() {
  const [isOpen, setIsOpen] = useState(false);
  const [isMinimized, setIsMinimized] = useState(false);
  const [sessionID, setSessionID] = useState(null);
  const [userName, setUserName] = useState('');
  const [userEmail, setUserEmail] = useState('');
  const [subject, setSubject] = useState('');
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [isConnected, setIsConnected] = useState(false);
  const [isTyping, setIsTyping] = useState(false);
  const [unreadCount, setUnreadCount] = useState(0);
  const [step, setStep] = useState('form'); // 'form' | 'chat'
  const [isLoading, setIsLoading] = useState(false);

  const ws = useRef(null);
  const messagesEndRef = useRef(null);
  const typingTimeoutRef = useRef(null);

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);

  // Load session from localStorage
  useEffect(() => {
    const savedSessionID = localStorage.getItem('support_session_id');
    if (savedSessionID) {
      setSessionID(savedSessionID);
      setStep('chat');
      loadPreviousMessages(savedSessionID);
      connectWebSocket(savedSessionID);
    }
  }, []);

  // Connect to WebSocket
  const connectWebSocket = (sid) => {
    const wsUrl = `${WS_URL}/ws/support/user/${sid}?name=${encodeURIComponent(userName || 'Guest')}&email=${encodeURIComponent(userEmail || '')}`;
    
    ws.current = new WebSocket(wsUrl);

    ws.current.onopen = () => {
      console.log('WebSocket connected');
      setIsConnected(true);
    };

    ws.current.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log('Received message:', data);

      if (data.type === 'message') {
        setMessages(prev => [...prev, data.message]);
        if (!isOpen) {
          setUnreadCount(prev => prev + 1);
        }
      } else if (data.type === 'typing') {
        setIsTyping(true);
        if (typingTimeoutRef.current) {
          clearTimeout(typingTimeoutRef.current);
        }
        typingTimeoutRef.current = setTimeout(() => setIsTyping(false), 3000);
      } else if (data.type === 'supporter_joined') {
        setMessages(prev => [...prev, data.message]);
      } else if (data.type === 'system') {
        setMessages(prev => [...prev, data.message]);
      }
    };

    ws.current.onerror = (error) => {
      console.error('WebSocket error:', error);
      setIsConnected(false);
    };

    ws.current.onclose = () => {
      console.log('WebSocket disconnected');
      setIsConnected(false);
    };
  };

  // Load previous messages
  const loadPreviousMessages = async (sid) => {
    try {
      const response = await axios.get(`${API_URL}/api/support/ticket/${sid}/messages`);
      setMessages(response.data || []);
    } catch (error) {
      console.error('Error loading messages:', error);
    }
  };

  // Create a new support ticket
  const handleStartChat = async () => {
    if (!userName || !subject) {
      alert('Please enter your name and subject');
      return;
    }

    setIsLoading(true);
    try {
      const response = await axios.post(`${API_URL}/api/support/ticket`, {
        user_name: userName,
        user_email: userEmail,
        subject: subject,
        priority: 'medium'
      });

      const newSessionID = response.data.session_id;
      setSessionID(newSessionID);
      localStorage.setItem('support_session_id', newSessionID);
      setStep('chat');
      connectWebSocket(newSessionID);
    } catch (error) {
      console.error('Error creating ticket:', error);
      alert('Failed to start chat. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  // Send a message
  const handleSendMessage = () => {
    if (!newMessage.trim() || !ws.current) return;

    const messageData = {
      type: 'message',
      message: newMessage
    };

    ws.current.send(JSON.stringify(messageData));
    setNewMessage('');
  };

  // Handle typing indicator
  const handleTyping = () => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({ type: 'typing' }));
    }
  };

  // Close chat
  const handleCloseChat = () => {
    if (ws.current) {
      ws.current.close();
    }
    setIsOpen(false);
  };

  // Open chat
  const handleOpenChat = () => {
    setIsOpen(true);
    setUnreadCount(0);
    if (isMinimized) {
      setIsMinimized(false);
    }
  };

  // End chat session
  const handleEndSession = () => {
    if (ws.current) {
      ws.current.close();
    }
    localStorage.removeItem('support_session_id');
    setSessionID(null);
    setMessages([]);
    setStep('form');
    setIsConnected(false);
  };

  return (
    <>
      {/* Chat Button */}
      {!isOpen && (
        <Box
          sx={{
            position: 'fixed',
            bottom: 20,
            right: 20,
            zIndex: 1000
          }}
        >
          <Badge badgeContent={unreadCount} color="error">
            <Button
              variant="contained"
              color="primary"
              onClick={handleOpenChat}
              sx={{
                borderRadius: '50px',
                width: '60px',
                height: '60px',
                minWidth: '60px',
                boxShadow: 3
              }}
            >
              <SupportAgentIcon fontSize="large" />
            </Button>
          </Badge>
        </Box>
      )}

      {/* Chat Window */}
      {isOpen && (
        <Paper
          elevation={6}
          sx={{
            position: 'fixed',
            bottom: 20,
            right: 20,
            width: '380px',
            height: isMinimized ? 'auto' : '600px',
            display: 'flex',
            flexDirection: 'column',
            zIndex: 1000,
            borderRadius: 2
          }}
        >
          {/* Header */}
          <Box
            sx={{
              bgcolor: 'primary.main',
              color: 'white',
              p: 2,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              borderTopLeftRadius: 8,
              borderTopRightRadius: 8
            }}
          >
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <SupportAgentIcon />
              <Typography variant="h6">Live Support</Typography>
              {isConnected && (
                <Chip
                  label="Online"
                  size="small"
                  sx={{ bgcolor: 'success.light', color: 'white' }}
                />
              )}
            </Box>
            <Box>
              <IconButton size="small" onClick={() => setIsMinimized(!isMinimized)} sx={{ color: 'white' }}>
                <MinimizeIcon />
              </IconButton>
              <IconButton size="small" onClick={handleCloseChat} sx={{ color: 'white' }}>
                <CloseIcon />
              </IconButton>
            </Box>
          </Box>

          {!isMinimized && (
            <>
              {/* Content */}
              <Box sx={{ flex: 1, overflow: 'hidden', display: 'flex', flexDirection: 'column' }}>
                {step === 'form' ? (
                  // Initial Form
                  <Box sx={{ p: 3, display: 'flex', flexDirection: 'column', gap: 2 }}>
                    <Typography variant="body1" gutterBottom>
                      Hi! How can we help you today?
                    </Typography>
                    <TextField
                      fullWidth
                      label="Your Name *"
                      value={userName}
                      onChange={(e) => setUserName(e.target.value)}
                      disabled={isLoading}
                    />
                    <TextField
                      fullWidth
                      label="Email (optional)"
                      type="email"
                      value={userEmail}
                      onChange={(e) => setUserEmail(e.target.value)}
                      disabled={isLoading}
                    />
                    <TextField
                      fullWidth
                      label="Subject *"
                      value={subject}
                      onChange={(e) => setSubject(e.target.value)}
                      disabled={isLoading}
                    />
                    <Button
                      variant="contained"
                      fullWidth
                      onClick={handleStartChat}
                      disabled={isLoading || !userName || !subject}
                      startIcon={isLoading ? <CircularProgress size={20} /> : <ChatIcon />}
                    >
                      {isLoading ? 'Starting Chat...' : 'Start Chat'}
                    </Button>
                  </Box>
                ) : (
                  // Chat Interface
                  <>
                    {/* Messages */}
                    <Box
                      sx={{
                        flex: 1,
                        overflowY: 'auto',
                        p: 2,
                        bgcolor: 'grey.50'
                      }}
                    >
                      {messages.map((msg, index) => (
                        <Box
                          key={index}
                          sx={{
                            mb: 2,
                            display: 'flex',
                            flexDirection: msg.sender_type === 'user' ? 'row-reverse' : 'row',
                            alignItems: 'flex-start'
                          }}
                        >
                          {msg.sender_type !== 'system' && (
                            <Avatar
                              sx={{
                                bgcolor: msg.sender_type === 'user' ? 'primary.main' : 'secondary.main',
                                width: 32,
                                height: 32,
                                fontSize: '0.875rem',
                                mx: 1
                              }}
                            >
                              {msg.sender_name?.[0]?.toUpperCase() || '?'}
                            </Avatar>
                          )}
                          <Box
                            sx={{
                              maxWidth: '70%',
                              bgcolor: msg.sender_type === 'system' ? 'grey.200' : msg.sender_type === 'user' ? 'primary.light' : 'white',
                              color: msg.sender_type === 'user' ? 'white' : 'text.primary',
                              p: 1.5,
                              borderRadius: 2,
                              boxShadow: msg.sender_type !== 'system' ? 1 : 0,
                              textAlign: msg.sender_type === 'system' ? 'center' : 'left',
                              width: msg.sender_type === 'system' ? '100%' : 'auto'
                            }}
                          >
                            {msg.sender_type !== 'system' && (
                              <Typography variant="caption" sx={{ fontWeight: 'bold', display: 'block', mb: 0.5 }}>
                                {msg.sender_name}
                              </Typography>
                            )}
                            <Typography variant="body2">{msg.message}</Typography>
                            <Typography variant="caption" sx={{ display: 'block', mt: 0.5, opacity: 0.7 }}>
                              {new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                            </Typography>
                          </Box>
                        </Box>
                      ))}
                      {isTyping && (
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, ml: 1 }}>
                          <Avatar sx={{ bgcolor: 'secondary.main', width: 32, height: 32 }}>S</Avatar>
                          <Typography variant="caption" sx={{ fontStyle: 'italic', color: 'text.secondary' }}>
                            Support agent is typing...
                          </Typography>
                        </Box>
                      )}
                      <div ref={messagesEndRef} />
                    </Box>

                    <Divider />

                    {/* Input */}
                    <Box sx={{ p: 2, bgcolor: 'white', display: 'flex', gap: 1 }}>
                      <TextField
                        fullWidth
                        size="small"
                        placeholder="Type your message..."
                        value={newMessage}
                        onChange={(e) => setNewMessage(e.target.value)}
                        onKeyPress={(e) => {
                          if (e.key === 'Enter' && !e.shiftKey) {
                            e.preventDefault();
                            handleSendMessage();
                          }
                        }}
                        onInput={handleTyping}
                        disabled={!isConnected}
                      />
                      <IconButton
                        color="primary"
                        onClick={handleSendMessage}
                        disabled={!newMessage.trim() || !isConnected}
                      >
                        <SendIcon />
                      </IconButton>
                    </Box>

                    {/* Footer */}
                    <Box sx={{ p: 1, textAlign: 'center', bgcolor: 'grey.100' }}>
                      <Button size="small" onClick={handleEndSession} color="error">
                        End Chat
                      </Button>
                    </Box>
                  </>
                )}
              </Box>
            </>
          )}
        </Paper>
      )}
    </>
  );
}

export default SupportChatWidget;

