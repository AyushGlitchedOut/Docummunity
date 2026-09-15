import type { UserInfo } from "../models/models";
import { createContext } from "react";

export const UserContext = createContext<UserInfo | null>(null);
