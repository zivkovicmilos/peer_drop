import { IContactResponse } from '../services/contacts/contactsService.types';

export enum ENewWorkspaceType {
  SEND_ONLY = 'Send only',
  RECEIVE_ONLY = 'Receive only',
  SEND_RECEIVE = 'Send & Receive'
}

export interface INewWorkspaceTypeWrapper {
  type: ENewWorkspaceType;
  description: string;
}

export enum ENWAccessControl {
  SPECIFIC_CONTACTS = 'Specific contacts',
  PASSWORD = 'Password'
}

export interface INWAccessControlContacts {
  contacts: IContactResponse[];
}

export interface INWAccessControlPassword {
  password: string;
}

export interface INWPermissions {
  autocloseWorkspace: {
    active: boolean;
    date?: Date;
  };

  enforcePeerLimit: {
    active: boolean;
    limit?: number;
  };

  additionalOwners: {
    active: boolean;
    contactIDs?: IContactResponse[];
  };
}
