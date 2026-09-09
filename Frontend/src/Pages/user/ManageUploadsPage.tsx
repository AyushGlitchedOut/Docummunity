import { Box, CircularProgress, Typography } from "@mui/material";
import XRecordCard from "../../Components/XRecordCard";
import type { RecordInfo } from "../../models/models";
import { useEffect, useState } from "react";
import { BACKEND_URL } from "../../consts";
import { useAuth } from "../../auth/fireBaseContext";
function ManageUploadsPage() {
  const auth = useAuth();

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
      if (response.status == 404) {
        setLoading(false);
        return;
      }

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
      ) : recordList.length > 0 ? (
        recordList.map((record) => {
          // Here I have put key as record.UUID because react will throw a warning in the console otherwise, and it says that using index as a key is not advised
          return (
            <XRecordCard
              height="60%"
              width="25%"
              record={record}
              key={record.UUID}
            />
          );
        })
      ) : (
        <Typography variant="h5">You dont have any uploads yet :)</Typography>
      )}
    </Box>
  );
}

export default ManageUploadsPage;
