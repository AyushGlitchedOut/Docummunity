import {
  Avatar,
  Box,
  Divider,
  IconButton,
  Paper,
  Typography,
} from "@mui/material";
import type { OtherUserRecordInfo, RecordInfo } from "../models/models";
import { BACKEND_URL } from "../consts";
import { Edit, HideImageSharp, Visibility } from "@mui/icons-material";
import { useNavigate } from "react-router-dom";
import { useState } from "react";

interface XRecordCardProps {
  height: string;
  width: string;
  //Overloaded so the card can load both records of the current user and other users too, though it does lead to code editor showing error on lines like 120 and 129
  record: RecordInfo | OtherUserRecordInfo;
}

//Already has Margin btw
export default function XRecordCard(props: XRecordCardProps) {
  const navigator = useNavigate();

  const previewIMGFilename = props.record.PREVIEW_IMG_PATH.split("/").at(-1);
  const [imageFound, setImageFound] = useState<boolean>(
    props.record.PREVIEW_IMG_PATH != "",
  );

  var forUsersRecords: boolean = true;

  //determine if record is RecordInfo or OtherUserRecordInfo
  if ("CREATOR_NAME" in props.record && "CREATOR_PROFILE_PIC" in props.record) {
    forUsersRecords = false;
  }

  //I was using the card element, but as it turns out there was little to no use of its features and i had to write a lot of sutomization myself, so theres that
  return (
    <Paper
      sx={(theme) => ({
        height: props.height,
        width: props.width,
        //for absolute positioning of the overlay buttons
        position: "relative",
        display: "flex",
        alignItems: "center",
        justifyContent: forUsersRecords ? "start" : "space-between",
        cursor: "pointer",
        flexDirection: "column",
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
        ":hover .overlay-button": {
          opacity: 1,
          pointerEvents: "auto",
        },
      })}
      variant="outlined"
    >
      <Box
        sx={{
          width: "100%",
          height: "60%",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          flexDirection: "column",
        }}
      >
        <Box></Box>
        {imageFound ? (
          <img
            src={`${BACKEND_URL}/data/PREVIEW/${previewIMGFilename}`}
            alt="No Preview Image found"
            onError={() => {
              setImageFound(false);
            }}
            style={{ maxWidth: "80%", maxHeight: "80%", marginTop: "5%" }}
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
      <Box sx={{ width: "100%", marginLeft: "1%" }}>
        <Typography
          variant="h5"
          overflow={"hidden"}
          textOverflow={"ellipsis"}
          sx={{ maxWidth: "100%" }}
          whiteSpace={"nowrap"}
        >
          {"" + props.record.NAME}
        </Typography>
        <Typography
          variant="subtitle2"
          sx={{
            maxWidth: "100%",
            display: "-webkit-box",
            WebkitLineClamp: forUsersRecords ? 4 : 3,
            WebkitBoxOrient: "vertical",
            overflow: "hidden",
          }}
        >
          {"->" + props.record.DESCRIPTION}
        </Typography>
      </Box>
      {forUsersRecords ? (
        <></>
      ) : (
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            justifyContent: "left",
            padding: "1%",
            flexDirection: "row",
            width: "100%",
          }}
        >
          <Avatar
            sx={{ marginLeft: "1%" }}
            src={`${BACKEND_URL}/user/PROFILE_PIC/${props.record.CREATOR_PROFILE_PIC.split("/").at(-1)}`}
          ></Avatar>
          <Typography
            variant="h5"
            sx={{ marginLeft: "2%", maxWidth: "100%" }}
            overflow={"hidden"}
            textOverflow={"ellipsis"}
            whiteSpace={"nowrap"}
          >
            {props.record.CREATOR_NAME}
          </Typography>
        </Box>
      )}

      {/* Extra buttons */}

      <IconButton
        className="overlay-button"
        size="medium"
        sx={{
          position: "absolute",
          backgroundColor: "white",
          border: "2px solid black",
          left: "80%",
          margin: "2%",
          width: "15%",
          aspectRatio: 1,
          minHeight: 0,
          minWidth: 0,
          opacity: 0,
          pointerEvents: "none",
          ":hover": { backgroundColor: "grey" },
        }}
        onClick={() => {
          navigator("/home/viewRecord/" + props.record.UUID);
        }}
        aria-label="view"
      >
        <Visibility />
      </IconButton>
      {forUsersRecords ? (
        <IconButton
          className="overlay-button"
          sx={{
            position: "absolute",
            backgroundColor: "white",
            top: "15%",
            left: "80%",
            margin: "2%",
            width: "15%",
            aspectRatio: 1,
            minHeight: 0,
            minWidth: 0,
            opacity: 0,
            border: "2px solid black",
            pointerEvents: "none",
            ":hover": { backgroundColor: "grey" },
          }}
          aria-label="edit"
        >
          <Edit />
        </IconButton>
      ) : (
        <></>
      )}
    </Paper>
  );
}
