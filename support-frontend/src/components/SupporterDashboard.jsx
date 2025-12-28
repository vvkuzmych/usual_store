import React, { useState, useEffect, useRef } from 'react';
import {
  Box,
  Container,
  Typography,
  Paper,
  List,
  ListItem,
  ListItemText,
  TextField,
  Button,
  IconButton,
  Avatar,
  Chip,
  Divider,
  Badge,
  Grid,
  Card,
  CardContent,
  AppBar,
  Toolbar,
  Menu,
  MenuItem
} from '@mui/material';
import {
  Send as SendIcon,
  Refresh as RefreshIcon,
  SupportAgent as SupportAgentIcon,
  Person as PersonIcon,
  Logout as LogoutIcon,
  PersonAdd as PersonAddIcon,
  AccountCircle as AccountCircleIcon,
  People as PeopleIcon
} from '@mui/icons-material';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import CreateSupporterAccount from './CreateSupporterAccount';

const API_URL = process.env.REACT_APP_SUPPORT_API_URL || 'http://localhost:5001';
const WS_URL = process.env.REACT_APP_SUPPORT_WS_URL || 'ws://localhost:5001';

function SupporterDashboard() {
  const { user, logout, isAdmin, isSuperAdmin } = useAuth();
  const navigate = useNavigate();
  
  const [tickets, setTickets] = useState([]);
  const [selectedTicket, setSelectedTicket] = useState(null);
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [supporterName, setSupporterName] = useState(user ? `${user.firstName} ${user.lastName}` : 'Support Agent');
  const [supporterID] = useState(user?.id || 1);
  const [isConnected, setIsConnected] = useState(false);
  const [isTyping, setIsTyping] = useState(false);
  const [showCreateAccount, setShowCreateAccount] = useState(false);
  const [anchorEl, setAnchorEl] = useState(null);

  const ws = useRef(null);
  const messagesEndRef = useRef(null);
  const typingTimeoutRef = useRef(null);

  const handleLogout = () => {
    if (ws.current) {
      ws.current.close();
    }
    logout();
    window.location.href = 'http://localhost:8000';
  };

  const handleMenuClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);

  // Load tickets on mount
  useEffect(() => {
    loadTickets();
    const interval = setInterval(loadTickets, 10000); // Refresh every 10 seconds
    return () => clearInterval(interval);
  }, []);

  // Load tickets
  const loadTickets = async () => {
    try {
      const response = await axios.get(`${API_URL}/api/support/tickets`);
      setTickets(response.data || []);
    } catch (error) {
      console.error('Error loading tickets:', error);
    }
  };

  // Select a ticket and connect to WebSocket
  const handleSelectTicket = async (ticket) => {
    setSelectedTicket(ticket);
    loadMessages(ticket.session_id);

    // Disconnect previous WebSocket if any
    if (ws.current) {
      ws.current.close();
    }

    // Assign ticket to supporter
    try {
      await axios.post(`${API_URL}/api/support/ticket/${ticket.id}/assign`, {
        supporter_id: supporterID
      });
    } catch (error) {
      console.error('Error assigning ticket:', error);
    }

    // Connect to WebSocket
    connectWebSocket(ticket.session_id);
  };

  // Connect to WebSocket
  const connectWebSocket = (sessionID) => {
    const wsUrl = `${WS_URL}/ws/support/supporter/${sessionID}?supporter_id=${supporterID}&name=${encodeURIComponent(supporterName)}`;
    
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
      } else if (data.type === 'typing') {
        setIsTyping(true);
        if (typingTimeoutRef.current) {
          clearTimeout(typingTimeoutRef.current);
        }
        typingTimeoutRef.current = setTimeout(() => setIsTyping(false), 3000);
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

  // Load messages for a ticket
  const loadMessages = async (sessionID) => {
    try {
      const response = await axios.get(`${API_URL}/api/support/ticket/${sessionID}/messages`);
      setMessages(response.data || []);
    } catch (error) {
      console.error('Error loading messages:', error);
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

  // Get priority color
  const getPriorityColor = (priority) => {
    switch (priority) {
      case 'urgent': return 'error';
      case 'high': return 'warning';
      case 'medium': return 'info';
      case 'low': return 'success';
      default: return 'default';
    }
  };

  // Get status color
  const getStatusColor = (status) => {
    switch (status) {
      case 'open': return 'error';
      case 'assigned': return 'warning';
      case 'in_progress': return 'info';
      case 'resolved': return 'success';
      case 'closed': return 'default';
      default: return 'default';
    }
  };

  return (
    <Box>
      {/* AppBar with User Info and Actions */}
      <AppBar position="static">
        <Toolbar>
          <SupportAgentIcon sx={{ mr: 2 }} />
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Support Dashboard
          </Typography>

          <Button
            color="inherit"
            startIcon={<RefreshIcon />}
            onClick={loadTickets}
            sx={{ mr: 2 }}
          >
            Refresh
          </Button>

          {isSuperAdmin() && (
            <Button
              color="inherit"
              startIcon={<PeopleIcon />}
              onClick={() => navigate('/support/users')}
              sx={{ mr: 2 }}
            >
              Manage Users
            </Button>
          )}

          {isAdmin() && (
            <Button
              color="inherit"
              startIcon={<PersonAddIcon />}
              onClick={() => setShowCreateAccount(true)}
              sx={{ mr: 2 }}
            >
              Create Supporter
            </Button>
          )}

          <IconButton
            color="inherit"
            onClick={handleMenuClick}
          >
            <AccountCircleIcon />
          </IconButton>
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
          >
            <MenuItem disabled>
              <Box>
                <Typography variant="body2">{supporterName}</Typography>
                <Typography variant="caption" color="text.secondary">
                  {user?.role === 'super_admin' 
                    ? 'Super Administrator' 
                    : user?.role === 'admin' 
                    ? 'Administrator' 
                    : 'Support Agent'}
                </Typography>
              </Box>
            </MenuItem>
            <Divider />
            <MenuItem onClick={() => { handleMenuClose(); handleLogout(); }}>
              <LogoutIcon fontSize="small" sx={{ mr: 1 }} />
              Logout
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>

      <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>

      <Grid container spacing={3}>
        {/* Ticket List */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ height: '80vh', display: 'flex', flexDirection: 'column' }}>
            <Box sx={{ p: 2, bgcolor: 'primary.main', color: 'white' }}>
              <Typography variant="h6">
                Open Tickets ({tickets.length})
              </Typography>
            </Box>
            <List sx={{ flex: 1, overflowY: 'auto' }}>
              {tickets.length === 0 ? (
                <Box sx={{ p: 3, textAlign: 'center' }}>
                  <Typography variant="body2" color="text.secondary">
                    No open tickets
                  </Typography>
                </Box>
              ) : (
                tickets.map((ticket) => (
                  <ListItem
                    key={ticket.id}
                    button
                    selected={selectedTicket?.id === ticket.id}
                    onClick={() => handleSelectTicket(ticket)}
                    sx={{
                      borderBottom: 1,
                      borderColor: 'divider',
                      '&.Mui-selected': {
                        bgcolor: 'primary.light',
                        color: 'white'
                      }
                    }}
                  >
                    <ListItemText
                      primary={
                        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
                          <Typography variant="subtitle1" sx={{ fontWeight: 'bold' }}>
                            {ticket.user_name}
                          </Typography>
                          <Chip
                            label={ticket.priority}
                            size="small"
                            color={getPriorityColor(ticket.priority)}
                          />
                        </Box>
                      }
                      secondary={
                        <Box>
                          <Typography variant="body2" sx={{ mb: 0.5 }}>
                            {ticket.subject}
                          </Typography>
                          <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
                            <Chip
                              label={ticket.status}
                              size="small"
                              color={getStatusColor(ticket.status)}
                              variant="outlined"
                            />
                            <Typography variant="caption" color="text.secondary">
                              {new Date(ticket.created_at).toLocaleString()}
                            </Typography>
                          </Box>
                        </Box>
                      }
                    />
                  </ListItem>
                ))
              )}
            </List>
          </Paper>
        </Grid>

        {/* Chat Area */}
        <Grid item xs={12} md={8}>
          {selectedTicket ? (
            <Paper sx={{ height: '80vh', display: 'flex', flexDirection: 'column' }}>
              {/* Header */}
              <Box sx={{ p: 2, bgcolor: 'primary.main', color: 'white' }}>
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                  <Box>
                    <Typography variant="h6">{selectedTicket.user_name}</Typography>
                    <Typography variant="body2">{selectedTicket.subject}</Typography>
                  </Box>
                  {isConnected && (
                    <Chip
                      label="Connected"
                      size="small"
                      sx={{ bgcolor: 'success.light', color: 'white' }}
                    />
                  )}
                </Box>
              </Box>

              {/* Messages */}
              <Box sx={{ flex: 1, overflowY: 'auto', p: 2, bgcolor: 'grey.50' }}>
                {messages.map((msg, index) => (
                  <Box
                    key={index}
                    sx={{
                      mb: 2,
                      display: 'flex',
                      flexDirection: msg.sender_type === 'supporter' ? 'row-reverse' : 'row',
                      alignItems: 'flex-start'
                    }}
                  >
                    {msg.sender_type !== 'system' && (
                      <Avatar
                        sx={{
                          bgcolor: msg.sender_type === 'supporter' ? 'secondary.main' : 'primary.main',
                          width: 40,
                          height: 40,
                          mx: 1
                        }}
                      >
                        {msg.sender_type === 'supporter' ? <SupportAgentIcon /> : <PersonIcon />}
                      </Avatar>
                    )}
                    <Box
                      sx={{
                        maxWidth: '70%',
                        bgcolor: msg.sender_type === 'system' ? 'grey.200' : msg.sender_type === 'supporter' ? 'secondary.light' : 'white',
                        color: msg.sender_type === 'supporter' ? 'white' : 'text.primary',
                        p: 2,
                        borderRadius: 2,
                        boxShadow: msg.sender_type !== 'system' ? 2 : 0,
                        textAlign: msg.sender_type === 'system' ? 'center' : 'left',
                        width: msg.sender_type === 'system' ? '100%' : 'auto'
                      }}
                    >
                      {msg.sender_type !== 'system' && (
                        <Typography variant="caption" sx={{ fontWeight: 'bold', display: 'block', mb: 0.5 }}>
                          {msg.sender_name}
                        </Typography>
                      )}
                      <Typography variant="body1">{msg.message}</Typography>
                      <Typography variant="caption" sx={{ display: 'block', mt: 0.5, opacity: 0.7 }}>
                        {new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                      </Typography>
                    </Box>
                  </Box>
                ))}
                {isTyping && (
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, ml: 1 }}>
                    <Avatar sx={{ bgcolor: 'primary.main', width: 40, height: 40 }}>
                      <PersonIcon />
                    </Avatar>
                    <Typography variant="caption" sx={{ fontStyle: 'italic', color: 'text.secondary' }}>
                      User is typing...
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
                  multiline
                  maxRows={4}
                  placeholder="Type your response..."
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
                  color="secondary"
                  onClick={handleSendMessage}
                  disabled={!newMessage.trim() || !isConnected}
                  sx={{ alignSelf: 'flex-end' }}
                >
                  <SendIcon />
                </IconButton>
              </Box>
            </Paper>
          ) : (
            <Card sx={{ height: '80vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              <CardContent sx={{ textAlign: 'center' }}>
                <SupportAgentIcon sx={{ fontSize: 80, color: 'text.secondary', mb: 2 }} />
                <Typography variant="h5" gutterBottom>
                  Select a ticket to start chatting
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Choose a ticket from the list on the left
                </Typography>
              </CardContent>
            </Card>
          )}
        </Grid>
      </Grid>
    </Container>

    {/* Create Supporter Account Dialog */}
    <CreateSupporterAccount
      open={showCreateAccount}
      onClose={() => setShowCreateAccount(false)}
      onSuccess={() => {
        setShowCreateAccount(false);
        loadTickets(); // Reload in case new supporter needs to be shown
      }}
    />
  </Box>
  );
}

export default SupporterDashboard;

