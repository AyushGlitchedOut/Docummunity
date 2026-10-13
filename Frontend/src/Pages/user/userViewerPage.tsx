import { Box } from "@mui/material";
import { useParams } from "react-router-dom";

export default function UserViewerPage() {
  const { uid } = useParams();

  return <Box>{uid}</Box>;
}
