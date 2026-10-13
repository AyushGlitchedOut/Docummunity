import { Avatar, Box, Divider, Paper, Typography } from "@mui/material";
import type { OtherUserInfo } from "../models/models";
import TimeAgo from "../helper/TimeAgoFormatter";
import { useState } from "react";
import { HideImageSharp, Widgets } from "@mui/icons-material";
import { BACKEND_URL } from "../consts";
import { useNavigate } from "react-router-dom";

interface XUserCardProps {
  height: string;
  width: string;
  user: OtherUserInfo;
}

export default function XUserCard(props: XUserCardProps) {
  const navigator = useNavigate();

  const profilePictureFilename = props.user.PROFILE_PIC.split("/").at(-1);
  const [imageFound, setImageFound] = useState<boolean>(
    props.user.PROFILE_PIC != "",
  );
  props.user.BIO =
    "Lorem Ipsum Dolor Sit Amet blah b;ah blag b;haoi vgldx hklguybxdilughrgliudxgurd gicfthg lixtdui rdxfhbiufd hiugdifughdkfxhihbxfhjgkdbxk hdt7 vikdtrchyfiubdhd";

  return (
    <Paper
      onClick={() => {
        navigator("/home/viewUser/" + props.user.UID);
      }}
      sx={(theme) => ({
        height: props.height,
        width: props.width,
        display: "flex",
        alignItems: "center",
        justifyContent: "start",
        flexDirection: "column",
        cursor: "pointer",
        margin: "1%",
        marginLeft: "5%",
        boxShadow: "5px 5px 5px " + theme.palette.divider,
        transition: `${theme.transitions.create(
          ["background-color", "transform"],
          {
            duration: theme.transitions.duration.standard,
          },
        )}`,
        ":hover": (theme) => ({
          backgroundColor: theme.palette.primary.main,
        }),
      })}
      variant="outlined"
    >
      <Box
        sx={{
          width: "100%",
          height: "100%",
        }}
      >
        <Box
          height={"50%"}
          sx={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            flexDirection: "column",
          }}
        >
          {imageFound ? (
            <Avatar
              onError={() => setImageFound(false)}
              src={`${BACKEND_URL}/user/PROFILE_PIC/${profilePictureFilename}`}
              sx={{
                height: "80%",
                aspectRatio: 1 / 1,
                width: "auto",
                marginTop: "5%",
              }}
            />
          ) : (
            <HideImageSharp sx={{ fontSize: "500%" }} />
          )}
          <Divider
            variant="fullWidth"
            sx={{ borderWidth: "2px", width: "100%", marginTop: "5%" }}
            orientation="horizontal"
          />
        </Box>
        <Typography variant="h5">
          {" "}
          {"->"} {props.user.DISPLAY_NAME}
        </Typography>
        <Typography sx={{ marginLeft: "2%" }}>
          Joined : {TimeAgo(props.user.CREATION_DATE)}
        </Typography>
        <Typography
          sx={{
            marginLeft: "2%",
            maxWidth: "100%",
            display: "-webkit-box",
            WebkitLineClamp: 4,
            WebkitBoxOrient: "vertical",
            overflow: "hidden",
          }}
        >
          About: {props.user.BIO}
        </Typography>
      </Box>
    </Paper>
  );
}
