import { Box, Button, Divider, Typography } from "@mui/material";
import XTextField from "../../Components/XTextField";
import { useMemo, useState } from "react";
import { useDropzone } from "react-dropzone";
import { useAuth } from "../../auth/fireBaseContext";
import { BACKEND_URL } from "../../consts";
import { useNavigate } from "react-router-dom";
import { UploadFile } from "@mui/icons-material";

function CreatePage() {
  const [title, setTitle] = useState("");
  const [Description, setDescription] = useState("");
  const [DocumentFile, setDocumentFile] = useState<File | null>(null);
  const [PreviewImageFile, setPreviewImageFile] = useState<File | null>(null);
  const navigator = useNavigate();
  const auth = useAuth();

  //for imageURL to change whenever the imageFile changes
  const imageURL = useMemo(() => {
    if (!PreviewImageFile) return null;
    return URL.createObjectURL(PreviewImageFile);
  }, [PreviewImageFile]);

  const DocumentPickerDropzone = useDropzone({
    onDrop(files) {
      setDocumentFile(files[0]);
    },
    multiple: false,
    maxFiles: 1,
    maxSize: 40 << 20,
    accept: {
      "application/pdf": [],
    },
  });
  const PreviewIMGPickerDropzone = useDropzone({
    onDrop(files) {
      setPreviewImageFile(files[0]);
    },
    multiple: false,
    maxFiles: 1,
    maxSize: 2 << 20,
    accept: {
      "image/jpeg": [],
      "image/png": [],
      "image/webp": [],
    },
  });

  async function handleSubmit(): Promise<void> {
    console.log("Title:", title);
    console.log("Description:", Description);
    console.log("DocumentFile Name:", DocumentFile?.name);
    console.log("Preview Image Name:", PreviewImageFile?.name);

    if (!title) {
      alert("Please Enter a Title");
      return;
    }
    if (!DocumentFile) {
      alert("Please select a Document File");
      return;
    }
    if (!auth || !auth.currentUser) {
      alert("Authentication Error");
      return;
    }
    const token = await auth!.currentUser!.getIdToken();
    if (!token) {
      alert("Authentication Error");
      return;
    }

    const form = new FormData();
    console.log(DocumentFile.type);
    console.log(PreviewImageFile?.type);
    form.append("FILE", DocumentFile);
    form.append("NAME", title);
    form.append("DESCRIPTION", Description);
    if (PreviewImageFile) {
      form.append("PREVIEW", PreviewImageFile);
    }

    const response = await fetch(BACKEND_URL + "/data/CREATE", {
      method: "POST",
      body: form,
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (!response.ok) {
      alert("Something Went Wrong!");
      if (response.status == 415) {
        alert("Unsupported File type. Please try valid filetypes");
      }
      return;
    }

    alert("Created Record Successfully");

    navigator("/home/manage_uploads");
  }

  //DOnt use Form input here since Dropzone doesnt support forms
  return (
    <Box
      sx={{
        height: "100%",
        width: "100%",
        borderRadius: "20px",
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        flexDirection: "row",
      }}
    >
      <Box
        sx={{
          height: "90%",
          width: "45%",
          margin: "1%",
        }}
      >
        <Box
          {...DocumentPickerDropzone.getRootProps()}
          sx={(theme) => ({
            backgroundColor: theme.palette.background.paper,
            borderStyle: "dotted",
            borderWidth: "5px",
            borderColor: theme.palette.divider,
            height: "50%",
            width: "100%",
            margin: "5%",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            cursor: "pointer",
          })}
        >
          {/* React Dropzone package's input  */}
          <input {...DocumentPickerDropzone.getInputProps()} />
          <UploadFile />
          <Typography variant="subtitle1">
            {DocumentFile ? DocumentFile.name : "Please enter a valid pdf File"}
          </Typography>
        </Box>
        <Box
          {...PreviewIMGPickerDropzone.getRootProps()}
          sx={(theme) => ({
            backgroundColor: "lightgray",
            borderStyle: "dotted",
            borderWidth: "5px",
            borderColor: theme.palette.divider,
            height: "40%",
            width: "100%",
            margin: "5%",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            cursor: "pointer",
          })}
        >
          {/* React Dropzone package's input  */}
          <input {...PreviewIMGPickerDropzone.getInputProps()} />
          {PreviewImageFile ? (
            <img
              src={imageURL!}
              style={{ maxWidth: "80%", maxHeight: "100%" }}
            ></img>
          ) : (
            <Typography>
              "Please enter a valid JPEG, PNG OR WEBP File (optional)"
            </Typography>
          )}
        </Box>
      </Box>
      <Divider
        sx={(theme) => ({
          backgroundColor: theme.palette.divider,
          width: "0.5%",
          height: "90%",
        })}
      />
      <Box
        sx={{
          height: "60%",
          width: "45%",
          margin: "1%",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-evenly",
          flexDirection: "column",
        }}
      >
        <XTextField
          labelName="Title:"
          id="record-title"
          label="Upload's Title"
          value={title}
          onChange={(event) => {
            setTitle(event.target.value);
          }}
        />
        <XTextField
          labelName="Description:"
          id="record-desc"
          label="Describe Your upload"
          value={Description}
          multiline
          minRows={2}
          maxRows={4}
          onChange={(event) => {
            setDescription(event.target.value);
          }}
        />
        <Button
          sx={{ width: "100%", height: "15%" }}
          variant="contained"
          onClick={() => {
            handleSubmit();
          }}
        >
          CREATE
        </Button>
      </Box>
    </Box>
  );
}

export default CreatePage;
