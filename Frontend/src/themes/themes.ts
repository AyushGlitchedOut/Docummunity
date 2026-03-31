import { createTheme } from "@mui/material";

export const defaultTheme = createTheme({
  palette: {
    text: {
      primary: "#1a1523",
      secondary: "#464141",
    },
    primary: {
      main: "#d36f3d",
    },
    secondary: {
      main: "#e0c67c",
    },
    background: {
      default: "#f2f6D5",
      paper: "#E4BE9E",
    },
    divider: "#777777",
  },
  typography: {
    fontFamily: "Inter, sans-serif",
  },
});
