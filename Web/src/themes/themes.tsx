import { createTheme, responsiveFontSizes } from "@mui/material";

export const AppLightTheme = responsiveFontSizes(
  createTheme({
    palette: {
      primary: {
        main: "#dd6401", //For larger, major elements e.g Navbar
      },
      secondary: {
        main: "#244d1a", //For smaller, detailing elements e.g. Icons, Buttons, Loading screens
      },
      divider: "#242331", //outlines
      background: {
        paper: "#d8a47f", // For cards e.g. document Card
        default: "#eb723e", //Default for the entire page
      },
      warning: {
        main: "#616161", //For dialogue boxes
      },
      text: {
        primary: "#333333", //text on light backgrounds
        secondary: "#BBBBBB", //text on dark backgrounds
      },
      error: {
        main: "#c62828", // For error messages, prohibition warning dialogues, Cancel buttons, delete dialogues, etc.
      },
    },
    typography: {
      fontFamily: "Roboto",
      h3: {
        fontSize: "2.5rem",
        textShadow: "2px 2px 1px grey",
      },
      button: {
        fontWeight: 800,
      },
    },
  }),
);
