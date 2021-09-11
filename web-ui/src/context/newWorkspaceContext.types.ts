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
  contactIDs: string[];
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
    contactIDs?: string[];
  };
}
