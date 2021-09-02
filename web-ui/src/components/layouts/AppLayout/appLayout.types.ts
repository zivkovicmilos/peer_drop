import {ReactNode} from "react";

export interface IAppLayoutProps {
  children?: ReactNode;
}

export enum EActiveAppTab {
  WORKSPACES = "Workspaces",
  CONTACTS = "Contacts",
  IDENTITIES = "Identities",
  SETTINGS = "Settings",
  OTHER = "Other"
}