import {
  CalendarMonth,
  CalendarMonthOutlined,
  FindInPage,
  PersonSearch,
  Search,
} from "@mui/icons-material";
import {
  Box,
  InputAdornment,
  MenuItem,
  Select,
  TextField,
  Checkbox,
  useTheme,
  FormControlLabel,
  IconButton,
  CircularProgress,
  Typography,
} from "@mui/material";
import { useState } from "react";
import type { OtherUserInfo, OtherUserRecordInfo } from "../../models/models";
import { BACKEND_URL } from "../../consts";
import XRecordCard from "../../Components/XRecordCard";
import XUserCard from "../../Components/XUserCard";

function SearchPage() {
  const [typeofResults, setTypeofResults] = useState<string>("records");
  const [useDescription, setUseDescription] = useState<boolean>(false);
  const [searchTags, setsearchTags] = useState<string[]>(["Coming Soon....."]);

  const [useDateTimeFilter, setUseDateTimeFilter] = useState<boolean>(false);
  const [dateFilterStart, setDateFilterStart] = useState<number>();
  const [dateFilterEnd, setDateFilterEnd] = useState<number>();
  const [loading, setLoading] = useState<boolean>(false);

  const [query, setQuery] = useState<string>("");

  const [recordSearchResults, setRecordSearchResults] = useState<
    OtherUserRecordInfo[]
  >([]);
  const [userSearchResults, setUserSearchResults] = useState<OtherUserInfo[]>(
    [],
  );

  const theme = useTheme();

  async function handleSearch(): Promise<void> {
    try {
      if (query == "") {
        return;
      }

      setLoading(true);
      var recordsSearched: boolean = false;

      //Clean Up other arrays just for good practice
      if (typeofResults == "records") {
        setUserSearchResults([]);
        recordsSearched = true;
      } else {
        setRecordSearchResults([]);
      }

      const params = new URLSearchParams();
      params.set("useDescription", useDescription ? "true" : "false");

      //fetch Results
      const searchRes = await fetch(
        `${BACKEND_URL}/${recordsSearched ? "data" : "user"}/SEARCH/${encodeURIComponent(query)}?${recordsSearched ? params.toString() : ""}`,
      );

      if (!searchRes.ok) {
        if (searchRes.status != 404) {
          alert("Something went Wrong");
          setLoading(false);
          return;
        }
        setLoading(false);
        return;
      }

      if (recordsSearched) {
        const records: OtherUserRecordInfo[] = (await searchRes.json()).message;
        setRecordSearchResults(records);
        console.log(records);
      } else {
        const records: OtherUserInfo[] = (await searchRes.json()).message;
        setUserSearchResults(records);
        console.log(records);
      }

      setLoading(false);
    } catch (error) {
      console.error("Lol");
    }
  }

  return (
    <Box
      sx={{
        height: "100%",
        width: "100%",
        borderRadius: "20px",
        display: "flex",
        alignItems: "center",
        justifyContent: "start",
        flexDirection: "column",
      }}
    >
      <form
        onSubmit={(event) => {
          event.preventDefault();
          handleSearch();
        }}
        style={{
          width: "95%",
        }}
      >
        <Box
          sx={{
            width: "100%",
            display: "flex",
            alignItems: "center",
            flexDirection: "column",
            justifyContent: "space-between",
          }}
        >
          <TextField
            fullWidth
            variant="outlined"
            type="search"
            label="Enter your Query"
            value={query}
            onChange={(event) => {
              setQuery(event.target.value);
            }}
            slotProps={{
              input: {
                startAdornment: (
                  <InputAdornment position="start">
                    {/* NOTE TO SELF: Remember that for a form to even have a submit action, it must have a submit button. If a submit button doesnt exist, even clicking enter will not trigger onSubmit */}
                    <IconButton type="submit">
                      <Search />
                    </IconButton>
                  </InputAdornment>
                ),
              },
            }}
            sx={{ margin: "0.5%" }}
          ></TextField>
          <Box
            sx={{
              margin: "0.5%",
              width: "100%",
              display: "flex",
              alignItems: "center",
              flexDirection: "row",
              justifyContent: "left",
            }}
          >
            <Select
              value={typeofResults}
              onChange={(event) => {
                setTypeofResults(event.target.value);
              }}
              sx={{
                width: "15%",
              }}
            >
              <MenuItem value={"records"}>
                <FindInPage /> Records
              </MenuItem>
              <MenuItem value={"users"}>
                <PersonSearch /> Users
              </MenuItem>
            </Select>
            <FormControlLabel
              control={
                <Checkbox
                  disabled={typeofResults != "records"}
                  checked={useDescription}
                  onChange={(event) => {
                    setUseDescription(event.target.checked);
                  }}
                />
              }
              label={"Description"}
              sx={{ marginLeft: "0.5%" }}
            />
            <Select
              value={searchTags}
              multiple
              disabled
              sx={{ width: "20%", marginLeft: "0.5%" }}
            >
              <MenuItem value={"Coming Soon....."}> Coming Soon.....</MenuItem>
            </Select>

            <Box
              sx={{
                marginLeft: "0.5%",
                //A webkit property, suppoted everywhere except firefox but its only cosmetic
                "& .calendar-picker::-webkit-calendar-picker-indicator": {
                  cursor: "pointer",
                },
              }}
            >
              <Checkbox
                checked={useDateTimeFilter}
                onChange={(event) => {
                  setUseDateTimeFilter(event.target.checked);
                }}
                icon={<CalendarMonthOutlined />}
                checkedIcon={<CalendarMonth />}
              />
              <input
                type="datetime-local"
                min={"1970-01-01T00:00"}
                disabled={!useDateTimeFilter}
                className="calendar-picker"
                style={{
                  backgroundColor: theme.palette.background.paper,
                }}
                onChange={(event) => {
                  const epochTime = Math.floor(
                    new Date(event.target.value).getTime() / 1000,
                  );
                  setDateFilterStart(epochTime);
                }}
              />{" "}
              to {/* Empty space here for padding*/}
              <input
                type="datetime-local"
                min={"1970-01-01T00:00"}
                disabled={!useDateTimeFilter}
                className="calendar-picker"
                style={{
                  backgroundColor: theme.palette.background.paper,
                }}
                onChange={(event) => {
                  const epochTime = Math.floor(
                    new Date(event.target.value).getTime() / 1000,
                  );
                  setDateFilterEnd(epochTime);
                }}
              />
            </Box>
          </Box>
        </Box>
      </form>
      <Box
        sx={(theme) => ({
          width: "100%",
          flex: 1,
          borderTop: "5px solid" + theme.palette.divider,
          display: "flex",
          alignItems: loading ? "center" : "start",
          justifyContent: loading ? "center" : "start",
          flexDirection: "row",
          flexWrap: "wrap",
          overflow: "auto",
        })}
      >
        {loading ? (
          <CircularProgress enableTrackSlot aria-label="loading" size={100} />
        ) : typeofResults == "records" ? (
          !recordSearchResults ? (
            <Typography variant="h6">No Records Found!</Typography>
          ) : (
            recordSearchResults.map((recordInfo) => {
              return (
                <XRecordCard
                  height="70%"
                  width="25%"
                  record={recordInfo}
                  key={recordInfo.UUID}
                />
              );
            })
          )
        ) : !userSearchResults ? (
          <Typography variant="h6">No Users Found</Typography>
        ) : (
          userSearchResults.map((userInfo) => {
            return (
              <XUserCard
                height="70%"
                width="25%"
                user={userInfo}
                key={userInfo.UID}
              />
            );
          })
        )}
      </Box>
    </Box>
  );
}

export default SearchPage;
