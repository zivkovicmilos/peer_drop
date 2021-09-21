import { ENewWorkspaceType } from '../../context/newWorkspaceContext.types';
import { IContactResponse } from '../contacts/contactsService.types';

export interface INewWorkspaceRequest {
  workspaceName: string;
  workspaceType: ENewWorkspaceType;
  workspaceAccessControlType: 'password' | 'contacts';
  baseWorkspaceOwnerKeyID: string;

  workspaceAccessControl: {
    password?: string;
    contacts?: IContactResponse[];
  };

  workspaceAdditionalOwnerPublicKeys: string[];
}

export interface INewWorkspaceResponse {
  mnemonic: string;
}

export interface IWorkspaceInfoResponse {
  name: string;
  mnemonic: string;
  workspaceOwnerPublicKeys: string[];
  securityType: string;
  workspaceType: string;

  passwordHash?: string;
  contactsWrapper?: {
    contactPublicKeys: string[];
  };
}

export interface IJoinWorkspaceRequest {
  mnemonic: string;
  password: string;
  publicKeyID: string;
}

export interface IWorkspaceListResponse {
  workspaceWrappers: IWorkspaceWrapper[];
  count: number;
}

export interface IWorkspaceWrapper {
  workspaceName: string;
  workspaceMnemonic: string;
}

export interface IWorkspaceDetailedResponse {
  workspaceMnemonic: string;
  workspaceName: string;
  workspaceType: string;
  workspaceFiles: IWorkspaceDetailedFileResponse[];
}

export interface IWorkspaceDetailedFileResponse {
  name: string;
  extension: string;
  size: number; // in bytes
  dateModified: number; // unix time
  checksum: string;
}

export interface IWorkspaceNumPeersResponse {
  numPeers: number;
}
