import { Box, CircularProgress } from "@mui/material";
import XRecordCard from "../../Components/XRecordCard";
import type { RecordInfo } from "../../models/models";
import { useEffect, useState } from "react";
import { BACKEND_URL } from "../../consts";
import { getAuth } from "firebase/auth";
import { useAuth } from "../../auth/fireBaseContext";
import { useNavigate } from "react-router-dom";

//DEBUG: test record
const testRecord: RecordInfo = {
  UUID: "01a05b89-35f0-7d91-9bf0-0d30f330a119",
  CREATION_DATE: "1788242114",
  CREATOR_ID: "XElQrGhyUQfXpK97U2Vrel8mc7k1",
  DESCRIPTION: "Lorem ipsum ",
  FILEPATH: "uploads/FILES/01a05b89-35f0-7d91-9bf0-0d30f330a119.pdf",
  NAME: "Parallelogram lorem ipsum dolor sit amet",
  PREVIEW_IMG_PATH: "uploads/PREVIEW/01a05b89-35f0-7d91-9bf0-0d30f330a119.jpg",
};

function ManageUploadsPage() {
  const auth = useAuth();
  const navigator = useNavigate();

  //TODO: will make loading screens other places too, where there needs to be loading
  const [recordList, setRecordList] = useState<RecordInfo[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  async function fetchUploads(): Promise<void> {
    setLoading(true);
    if (!auth || !auth.currentUser) {
      return;
    }

    const response = await fetch(
      `${BACKEND_URL}/user/RECORDS/${auth.currentUser.uid}`,
      {
        method: "GET",
      },
    );
    if (!response.ok) {
      alert("Something Went Wrong!");
      console.log(response.body);
      return;
    }

    const recordArray: RecordInfo[] = (await response.json()).message;
    setRecordList(recordArray);

    setLoading(false);
  }

  useEffect(() => {
    if (!auth) return;
    fetchUploads();
  }, [auth]);

  return (
    <Box
      sx={{
        height: "100%",
        width: "100%",
        borderRadius: "20px",
        display: "flex",
        alignItems: "start",
        justifyContent: "start",
        flexDirection: "row",
        flexWrap: "wrap",
        overflow: "auto",
      }}
    >
      {loading ? (
        <CircularProgress enableTrackSlot aria-label="loading" size={100} />
      ) : (
        recordList.map((record) => {
          return <XRecordCard height="60%" width="25%" record={record} />;
        })
      )}
    </Box>
  );
}

export default ManageUploadsPage;
