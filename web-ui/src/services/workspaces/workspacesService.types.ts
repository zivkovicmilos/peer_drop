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
