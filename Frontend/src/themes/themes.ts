import { createTheme } from "@mui/material";

export const defaultTheme = createTheme({
  palette: {
    text: {
      primary: "#3C3744",
      secondary: "#777777",
    },
    primary: {
      main: "#c67750",
    },
    secondary: {
      main: "#F2F6D0",
    },
    background: {
      default: "#f2f6D5",
      paper: "#E4BE9E",
    },
    divider: "#000000",
  },
  typography: {
    fontFamily: "Inter, sans-serif",
  },
});
