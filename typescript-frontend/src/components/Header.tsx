import { useState } from "react";
import { Link as RouterLink } from "react-router-dom";
import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import IconButton from "@mui/material/IconButton";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import ShoppingCartIcon from "@mui/icons-material/ShoppingCart";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";
import { useAuth } from "../contexts/AuthContext";

const Header = () => {
  const { user, isAuthenticated, logout } = useAuth();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    logout();
    handleMenuClose();
  };

  return (
    <AppBar position="sticky">
      <Toolbar>
        <Typography
          variant="h6"
          component={RouterLink}
          to="/"
          sx={{
            flexGrow: 1,
            textDecoration: "none",
            color: "inherit",
            fontWeight: 700,
            fontSize: "1.5rem",
          }}
        >
          🛍️ Usual Store <Typography component="span" sx={{ fontSize: "0.75rem", ml: 1, opacity: 0.8 }}>TypeScript</Typography>
        </Typography>

        <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
          <Button
            color="inherit"
            component={RouterLink}
            to="/"
            sx={{ fontWeight: 600 }}
          >
            Home
          </Button>
          <Button
            color="inherit"
            component={RouterLink}
            to="/products"
            sx={{ fontWeight: 600 }}
          >
            Products
          </Button>
          <IconButton color="inherit" component={RouterLink} to="/cart">
            <ShoppingCartIcon />
          </IconButton>

          {isAuthenticated && user ? (
            <>
              <IconButton onClick={handleMenuClick} sx={{ p: 0 }}>
                <Avatar
                  sx={{
                    bgcolor: "secondary.main",
                    width: 36,
                    height: 36,
                    fontSize: "1rem",
                  }}
                >
                  {user.first_name.charAt(0).toUpperCase()}
                </Avatar>
              </IconButton>
              <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleMenuClose}
                anchorOrigin={{
                  vertical: "bottom",
                  horizontal: "right",
                }}
                transformOrigin={{
                  vertical: "top",
                  horizontal: "right",
                }}
              >
                <MenuItem disabled>
                  <Typography variant="body2">
                    {user.first_name} {user.last_name}
                  </Typography>
                </MenuItem>
                <MenuItem onClick={handleLogout}>Logout</MenuItem>
              </Menu>
            </>
          ) : (
            <Button
              color="inherit"
              component={RouterLink}
              to="/login"
              startIcon={<AccountCircleIcon />}
              sx={{ fontWeight: 600 }}
            >
              Login
            </Button>
          )}
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header;

