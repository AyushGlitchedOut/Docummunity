import { Box, CircularProgress, Typography } from "@mui/material";
import XRecordCard from "../../Components/XRecordCard";
import type { RecordInfo } from "../../models/models";
import { useContext, useEffect, useState } from "react";
import { BACKEND_URL } from "../../consts";
import { UserContext } from "../../Contexts/UserContext";
function ManageUploadsPage() {
  //TODO: will make loading screens other places too, where there needs to be loading
  const [recordList, setRecordList] = useState<RecordInfo[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const user = useContext(UserContext);

  async function fetchUploads(): Promise<void> {
    setLoading(true);
    if (!user) {
      return;
    }

    const response = await fetch(`${BACKEND_URL}/user/RECORDS/${user.UID}`, {
      method: "GET",
    });
    if (!response.ok) {
      if (response.status == 404) {
        setLoading(false);
        return;
      }

      alert("Something Went Wrong!");
      console.log(response.body);
      setLoading(false);
      return;
    }

    const recordArray: RecordInfo[] = (await response.json()).message;
    setRecordList(recordArray);

    setLoading(false);
  }

  useEffect(() => {
    if (!user) return;
    fetchUploads();
  }, [user]);

  return (
    <Box
      sx={{
        height: "100%",
        width: "100%",
        borderRadius: "20px",
        display: "flex",
        alignItems: loading ? "center" : "start",
        justifyContent: loading ? "center" : "start",
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
