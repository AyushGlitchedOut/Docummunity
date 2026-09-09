import {
  Avatar,
  Box,
  Button,
  CircularProgress,
  Divider,
  Typography,
} from "@mui/material";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { BACKEND_URL } from "../../consts";
import type {
  OtherUserInfo,
  RecordInfo,
  RecordMetadata,
} from "../../models/models";
import { Download, Save, Share } from "@mui/icons-material";
import { filesize } from "filesize";

export default function RecordViewerPage() {
  const { uuid } = useParams();
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [data, setData] = useState<RecordInfo>();
  const [userInfo, setUserInfo] = useState<OtherUserInfo>();
  const [metadata, setMetadata] = useState<RecordMetadata>();
  const navigator = useNavigate();

  async function downloadFile(): Promise<void> {
    const file = await fetch(
      `${BACKEND_URL}/data/FILE/${data?.FILEPATH.split("/").at(-1)}`,
      { method: "GET" },
    );
    const fileContentBlob = await file.blob();

    const fileContentURL = URL.createObjectURL(fileContentBlob);

    //Do all this so the file gets downloaded with the actual title and not the UUID name
    const a = document.createElement("a");
    a.href = fileContentURL;
    a.download = data?.NAME ? data.NAME : "Document.pdf";
    a.click();

    URL.revokeObjectURL(fileContentURL);
  }

  async function fetchRecord(): Promise<void> {
    setIsLoading(true);

    //Obtain data about the record
    const response = await fetch(`${BACKEND_URL}/data/GET/${uuid}`, {
      method: "GET",
    });

    if (!response.ok) {
      if (response.status == 404) {
        alert("Invalid Record ID");
        navigator("/home");
        return;
      }
      alert("Something went wrong!!");
      navigator("/home");
      return;
    }

    const resData: RecordInfo = (await response.json()).message;
    setData(resData);

    //obtain info about the metadata
    const metaResponse = await fetch(
      `${BACKEND_URL}/data/FILE/${resData.FILEPATH.split("/").at(-1)}/meta`,
      { method: "GET" },
    );

    if (!response.ok) {
      alert("Error fetching metadata about record");
    }

    const metaResData: RecordMetadata = (await metaResponse.json()).message;
    setMetadata(metaResData);

    //obtain info about the user
    const creatorRes = await fetch(
      `${BACKEND_URL}/user/GET/${resData.CREATOR_ID}`,
    );

    if (!response.ok) {
      alert("Something Went Wrong while fetching the user");
    }

    const creatorResData: OtherUserInfo = (await creatorRes.json()).message;
    setUserInfo(creatorResData);

    setIsLoading(false);
  }

  useEffect(() => {
    fetchRecord();
  }, []);

  return (
    <Box
      sx={{
        height: "100%",
        width: "100%",
        borderRadius: "20px",
        display: "flex",
        alignItems: "center",
        justifyContent: isLoading ? "center" : "space-between",
      }}
    >
      {isLoading ? (
        <CircularProgress enableTrackSlot aria-label="loading" size={100} />
      ) : (
        <>
          <Box
            sx={{
              height: "90%",
              width: "45%",
              margin: "1%",
            }}
          >
            <img
              src={`${BACKEND_URL}/data/PREVIEW/${data?.PREVIEW_IMG_PATH.split("/").at(-1)}`}
              style={{
                maxWidth: "100%",
                maxHeight: "50%",
              }}
            />
            <Typography variant="subtitle2">By:</Typography>
            <Box
              onClick={() => {
                alert("Displaying User");
              }}
              sx={(theme) => ({
                display: "flex",
                alignItems: "center",
                justifyContent: "left",
                flexDirection: "row",
                cursor: "pointer",
                backgroundColor: theme.palette.background.paper,
                margin: "1%",
                padding: "1%",
                transition: `${theme.transitions.create(
                  ["background-color", "transform"],
                  {
                    duration: theme.transitions.duration.standard,
                  },
                )}`,
                ":hover": (theme) => ({
                  backgroundColor: theme.palette.divider,
                }),
              })}
            >
              <Avatar
                src={`${BACKEND_URL}/user/PROFILE_PIC/${userInfo?.PROFILE_PIC.split("/").at(-1)}`}
                sx={{ marginRight: "1%" }}
              />
              <Typography>{userInfo?.DISPLAY_NAME}</Typography>
            </Box>
            <Button
              variant="contained"
              sx={{ height: "10%", width: "100%" }}
              onClick={() => {
                downloadFile();
              }}
            >
              <Download /> <Typography variant="button">DOWNLOAD</Typography>
            </Button>
            <Box
              sx={{
                margin: "1%",
                display: "flex",
                alignItems: "center",
                justifyContent: "center  ",
              }}
            >
              {/* TODO: Implement save feature */}
              <Button
                variant="contained"
                sx={{ width: "45%", marginRight: "1%" }}
                onClick={() => {
                  alert("Save feature to implement");
                }}
              >
                <Save /> <Typography variant="button">SAVE</Typography>
              </Button>
              {/* TODO: Implement sharing feature */}
              <Button
                variant="contained"
                sx={{ width: "45%" }}
                onClick={() => {
                  alert("Sharing feature to implement");
                }}
              >
                <Share /> <Typography variant="button">SHARE</Typography>
              </Button>
            </Box>
          </Box>
          <Divider
            orientation="vertical"
            sx={(theme) => ({
              backgroundColor: theme.palette.divider,
              height: "90%",
              width: "0.5%",
            })}
          />
          <Box
            sx={{
              height: "90%",
              width: "45%",
              margin: "1%",
            }}
          >
            <Typography variant="h6">
              Title:<Typography variant="h5">{data?.NAME}</Typography>
            </Typography>
            <Typography variant="h6" sx={{ marginTop: "2%" }}>
              Desc: <Typography variant="body1">{data?.DESCRIPTION}</Typography>
            </Typography>
            <Typography variant="h6" sx={{ marginTop: "2%" }}>
              Size:{" "}
              <Typography variant="body1">
                {(() => {
                  try {
                    if (!metadata?.SIZE) return "Unknown size";
                    return filesize(metadata?.SIZE);
                  } catch {
                    return "Unknown size";
                  }
                })()}
              </Typography>
            </Typography>
            <Typography variant="h6" sx={{ marginTop: "2%" }}>
              Created At:
              <Typography variant="body1">
                {(() => {
                  if (!data?.CREATION_DATE) return "Unknown Date";
                  const formattedTime = new Date(
                    Number(data?.CREATION_DATE) * 1000,
                  );
                  return formattedTime.toLocaleDateString();
                })()}
              </Typography>
            </Typography>
            <Typography variant="h6" sx={{ marginTop: "2%" }}>
              Tags: <Typography variant="body1">Coming Soon.....</Typography>
            </Typography>
          </Box>
        </>
      )}
    </Box>
  );
}
