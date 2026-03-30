import { Box, Button } from "@mui/material";

function SideBar() {
  return (
    <Box
      sx={(theme) => ({
        border: `1px solid ${theme.palette.divider}`,
        backgroundColor: theme.palette.secondary.main,
        borderRadius: "20px",
        margin: "1%",
        width: "25%",
        height: "80%",
      })}
    >
      SIDEBAR
    </Box>
  );
}
export default SideBar;
