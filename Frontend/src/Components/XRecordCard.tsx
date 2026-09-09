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
import { Edit, Visibility } from "@mui/icons-material";
import { useNavigate } from "react-router-dom";

interface XRecordCardProps {
  height: string;
  width: string;
  //Overloaded so the card can load both records of the current user and other users too, though it does lead to code editor showing error on lines like 120 and 129
  record: RecordInfo | OtherUserRecordInfo;
}

//Already has Margin btw
export default function XRecordCard(props: XRecordCardProps) {
  const navigator = useNavigate();
  const previewURLparts = props.record.PREVIEW_IMG_PATH.split("/");
  const previewIMGFilename = previewURLparts[previewURLparts.length - 1];

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
          maxHeight: "60%",
          minHeight: "40%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "column",
        }}
      >
        <img
          src={`${BACKEND_URL}/data/PREVIEW/${previewIMGFilename}`}
          style={{ minWidth: "40", maxWidth: "80%", marginTop: "5%" }}
        />
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
            src={props.record.CREATOR_PROFILE_PIC}
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
