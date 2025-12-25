import { createTheme } from '@mui/material/styles';

// Custom theme matching your brand colors
const theme = createTheme({
  palette: {
    primary: {
      main: '#6a11cb', // Purple
      light: '#667eea',
      dark: '#4a0d8b',
      contrastText: '#fff',
    },
    secondary: {
      main: '#764ba2', // Purple gradient end
      light: '#9b6bc7',
      dark: '#5a3780',
      contrastText: '#fff',
    },
    success: {
      main: '#28a745',
      light: '#48c461',
      dark: '#1e7e34',
    },
    error: {
      main: '#dc3545',
      light: '#e4606d',
      dark: '#bd2130',
    },
    warning: {
      main: '#ffc107',
      light: '#ffd54f',
      dark: '#c79100',
    },
    background: {
      default: '#f5f7fa',
      paper: '#ffffff',
    },
    text: {
      primary: '#2c3e50',
      secondary: '#555555',
    },
  },
  typography: {
    fontFamily: [
      '-apple-system',
      'BlinkMacSystemFont',
      '"Segoe UI"',
      'Roboto',
      '"Helvetica Neue"',
      'Arial',
      'sans-serif',
    ].join(','),
    h1: {
      fontSize: '2.8rem',
      fontWeight: 700,
      lineHeight: 1.2,
    },
    h2: {
      fontSize: '2.2rem',
      fontWeight: 700,
      lineHeight: 1.3,
    },
    h3: {
      fontSize: '1.8rem',
      fontWeight: 700,
      lineHeight: 1.4,
    },
    h4: {
      fontSize: '1.5rem',
      fontWeight: 600,
      lineHeight: 1.4,
    },
    h5: {
      fontSize: '1.25rem',
      fontWeight: 600,
      lineHeight: 1.5,
    },
    h6: {
      fontSize: '1rem',
      fontWeight: 600,
      lineHeight: 1.5,
    },
    button: {
      textTransform: 'none', // Don't uppercase buttons
      fontWeight: 600,
    },
  },
  shape: {
    borderRadius: 8, // Rounded corners
  },
  shadows: [
    'none',
    '0 2px 4px rgba(0,0,0,0.05)',
    '0 4px 8px rgba(0,0,0,0.08)',
    '0 6px 12px rgba(0,0,0,0.1)',
    '0 8px 16px rgba(0,0,0,0.12)',
    '0 10px 20px rgba(0,0,0,0.15)',
    '0 12px 24px rgba(0,0,0,0.18)',
    '0 14px 28px rgba(0,0,0,0.2)',
    '0 16px 32px rgba(0,0,0,0.22)',
    '0 18px 36px rgba(0,0,0,0.25)',
    '0 20px 40px rgba(0,0,0,0.3)',
    '0 22px 44px rgba(0,0,0,0.32)',
    '0 24px 48px rgba(0,0,0,0.35)',
    '0 26px 52px rgba(0,0,0,0.37)',
    '0 28px 56px rgba(0,0,0,0.4)',
    '0 30px 60px rgba(0,0,0,0.42)',
    '0 32px 64px rgba(0,0,0,0.45)',
    '0 34px 68px rgba(0,0,0,0.47)',
    '0 36px 72px rgba(0,0,0,0.5)',
    '0 38px 76px rgba(0,0,0,0.52)',
    '0 40px 80px rgba(0,0,0,0.55)',
    '0 42px 84px rgba(0,0,0,0.57)',
    '0 44px 88px rgba(0,0,0,0.6)',
    '0 46px 92px rgba(0,0,0,0.62)',
    '0 48px 96px rgba(0,0,0,0.65)',
  ],
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          padding: '10px 24px',
          fontSize: '1rem',
        },
        contained: {
          boxShadow: '0 4px 12px rgba(106, 17, 203, 0.3)',
          '&:hover': {
            boxShadow: '0 6px 20px rgba(106, 17, 203, 0.5)',
            transform: 'translateY(-2px)',
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          boxShadow: '0 4px 15px rgba(0, 0, 0, 0.08)',
          transition: 'all 0.3s ease',
          '&:hover': {
            transform: 'translateY(-8px)',
            boxShadow: '0 8px 25px rgba(0, 0, 0, 0.15)',
          },
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            '&:hover fieldset': {
              borderColor: '#6a11cb',
            },
          },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 20,
          fontWeight: 600,
        },
      },
    },
  },
});

export default theme;

