export interface UserInfo {
  UID: string;
  DISPLAY_NAME: string;
  BIO: string;
  PROFILE_PIC: string;
  SETTINGS: string;
  CREATION_DATE: string;
}

//CHANGE THIS IF BACKEND SCHEMA CHANGES, ALSO EVERY NAME NEEDS TO MATCH SAME AS BACKEND RESPONSES
export interface RecordInfo {
  UUID: string;
  NAME: string;
  DESCRIPTION: string;
  FILEPATH: string;
  CREATOR_ID: string;
  PREVIEW_IMG_PATH: string;
  CREATION_DATE: string;
}
export interface OtherUserRecordInfo extends RecordInfo {
  CREATOR_NAME?: string;
  CREATOR_PROFILE_PIC?: string;
}
