import {
  ENewWorkspaceType,
  ENWAccessControl
} from '../../context/newWorkspaceContext.types';
import { IContactResponse } from '../contacts/contactsService.types';

export interface INewWorkspaceRequest {
  workspaceName: string;
  workspaceType: ENewWorkspaceType;
  workspaceAccessControlType: ENWAccessControl;
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