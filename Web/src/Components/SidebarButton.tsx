import { Button, Typography } from "@mui/material";

interface SidebarButtonProps {
  text: string;
  Icon: React.ReactNode;
  ActiveTab: boolean;
  callback: Function;
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
        backgroundColor: args.ActiveTab
          ? "secondary.main"
          : "background.default",
        color: args.ActiveTab ? "primary.main" : "secondary.main",
      }}
      onClick={() => args.callback()}
    >
      {args.Icon}
      <Typography sx={{ fontSize: "110%", width: "100%" }}>
        {args.text}
      </Typography>
    </Button>
  );
}

export default SidebarButton;
