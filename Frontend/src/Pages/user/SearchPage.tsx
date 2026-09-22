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
} from "@mui/material";
import { useState } from "react";

function SearchPage() {
  const [typeofResults, setTypeofResults] = useState<string>("records");
  const [searchTags, setsearchTags] = useState<string[]>(["Coming Soon....."]);
  const [useDateTimeFilter, setUseDateTimeFilter] = useState<boolean>(false);
  const [dateFilterStart, setDateFilterStart] = useState<number>();
  const [dateFilterEnd, setDateFilterEnd] = useState<number>();
  const theme = useTheme();

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
      <Box
        sx={{
          width: "95%",
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
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <Search />
                </InputAdornment>
              ),
            },
          }}
          sx={{ margin: "1%" }}
        ></TextField>
        <Box
          sx={{
            margin: "1%",
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
          <Select
            value={searchTags}
            multiple
            disabled
            sx={{ width: "20%", marginLeft: "1%" }}
          >
            <MenuItem value={"Coming Soon....."}> Coming Soon.....</MenuItem>
          </Select>

          <Box
            sx={{
              marginLeft: "1%",
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
      <Box
        sx={{
          width: "95%",
          height: "80%",
          backgroundColor: "red",
        }}
      ></Box>
    </Box>
  );
}

export default SearchPage;
