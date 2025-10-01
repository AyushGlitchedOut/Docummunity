import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";
import { ThemeProvider } from "@emotion/react";
import { AppLightTheme } from "./themes/themes.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ThemeProvider theme={AppLightTheme}>
      <App />
    </ThemeProvider>
  </StrictMode>,
);
