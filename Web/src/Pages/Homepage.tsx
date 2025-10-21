import { SearchRounded } from "@mui/icons-material";
import {
  Box,
  Button,
  Card,
  CardActions,
  CardContent,
  CardMedia,
  Divider,
  IconButton,
  InputAdornment,
  TextField,
  Typography,
} from "@mui/material";
import PdfIcon from "../assets/Icons/pdf.png";
import TextIcon from "../assets/Icons/txt.png";
import XMLIcon from "../assets/Icons/xml.png";
import MDIcon from "../assets/Icons/md.png";
import MSwordIcon from "../assets/Icons/docx.png";
import MSexcelIcon from "../assets/Icons/xlsx.png";
import MSpowerpointIcon from "../assets/Icons/pptx.png";
import ODFtextIcon from "../assets/Icons/odt.png";
import ODFspreadsheetIcon from "../assets/Icons/ods.png";
import ODFpresentationIcon from "../assets/Icons/odp.png";
import { useEffect, useRef, useState } from "react";
import Marquee from "react-fast-marquee";
import { useFirebase } from "../services/firebase";
import { useNavigate } from "react-router-dom";

let supportedFileFormats: FileTypeIconDisplayProps[] = [
  {
    extensionInfo: ".pdf (PDF file)",
    imagePath: PdfIcon,
    linkToInfo: "https://en.wikipedia.org/wiki/PDF",
  },
  {
    extensionInfo: ".txt (Text File)",
    imagePath: TextIcon,
    linkToInfo: "https://en.wikipedia.org/wiki/Text_file",
  },
  {
    extensionInfo: ".xml (XML File)",
    imagePath: XMLIcon,
    linkToInfo: "https://en.wikipedia.org/wiki/XML",
  },
  {
    extensionInfo: ".md (Markdown File)",
    imagePath: MDIcon,
    linkToInfo: "https://en.wikipedia.org/wiki/Markdown ",
  },
  {
    extensionInfo: ".docx (MS Word)",
    imagePath: MSwordIcon,
    linkToInfo: "https://en.wikipedia.org/wiki/Office_Open_XML",
  },
  {
    extensionInfo: ".xlsx (MS Excel)",
    imagePath: MSexcelIcon,
    linkToInfo: "https://en.wikipedia.org/wiki/Office_Open_XML",
  },
  {
    extensionInfo: ".pptx (MS Powerpoint)",
    imagePath: MSpowerpointIcon,
    linkToInfo: "https://en.wikipedia.org/wiki/Office_Open_XML",
  },
  {
    extensionInfo: ".odt (Opendocument text)",
    imagePath: ODFtextIcon,
    linkToInfo: "https://en.wikipedia.org/wiki/OpenDocument",
  },
  {
    extensionInfo: ".ods (Opendocument Spreadsheet)",
    imagePath: ODFspreadsheetIcon,
    linkToInfo: "https://en.wikipedia.org/wiki/OpenDocument",
  },
  {
    extensionInfo: ".odp (Opendocument Presentation)",
    imagePath: ODFpresentationIcon,
    linkToInfo: "https://en.wikipedia.org/wiki/OpenDocument",
  },
];

function Homepage() {
  const [query, setQuery] = useState("");
  const infoRef = useRef<HTMLDivElement | null>(null);
  const firebase = useFirebase();
  const navigator = useNavigate();

  function scrollToInfo() {
    infoRef.current?.scrollIntoView({ behavior: "smooth" });
  }
  useEffect(() => {
    const withoutHash = window.location.hash.slice(1);
    const query = withoutHash.split("?")[1];

    if (query) {
      const params = new URLSearchParams(query);
      const tab = params.get("tab");
      if (tab === "info") {
        scrollToInfo();
      }
    }
  }, []);
  function handleSubmit() {
    if (query == "") {
      return;
    }
    if (query.length > 70) {
      alert("Too Long search Query!");
      return;
    }
    alert("Searching for:" + query);
  }

  useEffect(() => {
    if (firebase.isLoggedIn) {
      navigator("/home");
    }
  }, [firebase]);

  return (
    <>
      <Box
        sx={{
          background:
            "linear-gradient(135deg,rgba(40,40,40,0.9), rgba(100,100,100,0.85), rgba(220,220,220,0.5))",
          width: "100%",
          position: "absolute",
          backgroundAttachment: "fixed",
          backgroundPosition: "center",
          height: "92vh",
        }}
      ></Box>
      <Box
        sx={{
          width: "100%",
          height: "92vh",
          overflow: "auto",
          position: "absolute",
        }}
      >
        {/* Main Content */}
        <Box
          sx={{
            width: "99.5%",
            height: "70%",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-evenly",
            flexDirection: "column",
          }}
        >
          {/* All the supported formats display */}
          <Box
            sx={{
              backgroundColor: "rgba(20,20,20,0.5)",
              width: "100%",
              height: "60%",
              display: "flex",
              flexDirection: "column",
            }}
          >
            <Box sx={{ height: "15%", width: "100%", textAlign: "center" }}>
              <Typography
                variant="h4"
                sx={{ color: "lightgray", textShadow: "2px 2px 1px grey" }}
              >
                Formats We Support!
              </Typography>
            </Box>
            <Marquee
              speed={70}
              autoFill
              style={{
                width: "100%",
                height: "100%",
              }}
            >
              {supportedFileFormats.map((args) => {
                return (
                  <FileTypeIconDisplay
                    extensionInfo={args.extensionInfo}
                    imagePath={args.imagePath}
                    linkToInfo={args.linkToInfo}
                  />
                );
              })}
            </Marquee>
          </Box>

          {/* The searchBar */}
          <Box
            component="form"
            onSubmit={() => handleSubmit()}
            sx={{ width: "60%" }}
          >
            <TextField
              variant="outlined"
              value={query}
              placeholder="Search.."
              onChange={(event) => setQuery(event.target.value)}
              sx={{ width: "100%", boxShadow: "1px 1px 5px black" }}
              slotProps={{
                input: {
                  sx: {
                    ...{
                      backgroundColor: "background.paper",
                      color: "text.secondary",
                    },
                  },
                  startAdornment: (
                    <InputAdornment position="start">
                      <IconButton onClick={() => handleSubmit()}>
                        <SearchRounded fontSize="large" />
                      </IconButton>
                    </InputAdornment>
                  ),
                },
              }}
            />
          </Box>
        </Box>
        <Divider variant="middle" />
        {/* Additional Details */}
        <Box id="info" ref={infoRef}>
          <h1>Additional Info (To be made later)</h1>
          <h1>Lorem</h1>
          <h1>Ipsum</h1>
          <h1>Dolor</h1>
          <h1>Sit</h1>
          <h1>Amet</h1>
        </Box>
      </Box>
    </>
  );
}

interface FileTypeIconDisplayProps {
  imagePath: string;
  extensionInfo: string;
  linkToInfo: string;
}

function FileTypeIconDisplay(args: FileTypeIconDisplayProps) {
  return (
    <Card
      sx={{
        width: "33%",
        aspectRatio: "3/4",
        backgroundColor: "rgba(180,180,180,0.1)",
        boxShadow: "1px 1px 1px grey",
      }}
    >
      <CardMedia
        sx={{
          margin: "5%",
          height: "50%",
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <img
          style={{ height: "100%" }}
          alt="No Icon found :("
          src={args.imagePath}
        />
      </CardMedia>
      <CardContent
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "4%",
          color: "white",
        }}
      >
        <Typography variant="body1">{args.extensionInfo}</Typography>
      </CardContent>
      <CardActions>
        <Button
          variant="text"
          sx={{ color: "text.secondary" }}
          onClick={() => {
            window.open(args.linkToInfo);
          }}
        >
          <Typography variant="button">Learn More</Typography>
        </Button>
      </CardActions>
    </Card>
  );
}

export default Homepage;
