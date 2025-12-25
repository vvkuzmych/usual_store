import { createTheme, ThemeOptions } from "@mui/material/styles";

const themeOptions: ThemeOptions = {
  palette: {
    mode: "light",
    primary: {
      main: "#2196f3", // Blue theme for TypeScript frontend
      light: "#64b5f6",
      dark: "#1976d2",
      contrastText: "#fff",
    },
    secondary: {
      main: "#00bcd4", // Cyan
      light: "#4dd0e1",
      dark: "#0097a7",
      contrastText: "#fff",
    },
    success: {
      main: "#4caf50",
    },
    error: {
      main: "#f44336",
    },
    warning: {
      main: "#ff9800",
    },
    background: {
      default: "#f5f5f5",
      paper: "#ffffff",
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontSize: "2.5rem",
      fontWeight: 700,
      lineHeight: 1.2,
    },
    h2: {
      fontSize: "2rem",
      fontWeight: 600,
      lineHeight: 1.3,
    },
    h3: {
      fontSize: "1.75rem",
      fontWeight: 600,
    },
    h4: {
      fontSize: "1.5rem",
      fontWeight: 600,
    },
    h5: {
      fontSize: "1.25rem",
      fontWeight: 500,
    },
    h6: {
      fontSize: "1rem",
      fontWeight: 500,
    },
    body1: {
      fontSize: "1rem",
      lineHeight: 1.5,
    },
    button: {
      textTransform: "none",
      fontWeight: 600,
    },
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          padding: "8px 24px",
          boxShadow: "none",
          "&:hover": {
            boxShadow: "0 4px 8px rgba(0,0,0,0.1)",
          },
        },
        containedPrimary: {
          background: "linear-gradient(45deg, #2196f3 30%, #21cbf3 90%)",
          "&:hover": {
            background: "linear-gradient(45deg, #1976d2 30%, #00bcd4 90%)",
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          boxShadow: "0 2px 12px rgba(0,0,0,0.08)",
          transition: "all 0.3s ease-in-out",
          "&:hover": {
            transform: "translateY(-4px)",
            boxShadow: "0 6px 20px rgba(0,0,0,0.12)",
          },
        },
      },
    },
    MuiTextField: {
      defaultProps: {
        variant: "outlined",
        fullWidth: true,
      },
      styleOverrides: {
        root: {
          marginBottom: 16,
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          background: "linear-gradient(45deg, #2196f3 30%, #21cbf3 90%)",
          boxShadow: "0 2px 4px rgba(0,0,0,0.1)",
        },
      },
    },
  },
};

const theme = createTheme(themeOptions);

export default theme;

