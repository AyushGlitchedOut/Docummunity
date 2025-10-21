import { Button, Typography } from "@mui/material";

interface SidebarButtonProps {
  text: string;
  Icon?: React.ReactNode;
}

function SidebarButton(args: SidebarButtonProps) {
  return (
    <Button
      sx={{
        margin: "5%",
        width: "80%",
        height: "8%",
        display: "flex",
        alignItems: "center",
        justifyContent: "space-evenly",
        backgroundColor: "background.default",
        color: "secondary.main",
      }}
    >
      {args.Icon}
      <Typography sx={{ fontSize: "110%", width: "100%" }}>
        {args.text}
      </Typography>
    </Button>
  );
}

export default SidebarButton;
