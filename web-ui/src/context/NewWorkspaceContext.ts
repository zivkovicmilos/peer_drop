import { createContext } from 'react';
import {
  ENewWorkspaceType,
  INWAccessControlContacts,
  INWAccessControlPassword,
  INWPermissions
} from './newWorkspaceContext.types';

export interface INewWorkspaceContext {
  workspaceName: string;
  workspaceType: ENewWorkspaceType;
  accessControl: INWAccessControlContacts | INWAccessControlPassword;
  permissions: INWPermissions;

  setWorkspaceName: (name: string) => void;
  setWorkspaceType: (type: ENewWorkspaceType) => void;
  setAccessControl: (
    control: INWAccessControlContacts | INWAccessControlPassword
  ) => void;
  setPermissions: (permissions: INWPermissions) => void;
}

const NewWorkspaceContext = createContext<INewWorkspaceContext>({
  workspaceName: null,
  workspaceType: null,
  accessControl: null,
  permissions: null,

  setWorkspaceName: (name: string) => {},
  setWorkspaceType: (type: ENewWorkspaceType) => {},
  setAccessControl: (
    control: INWAccessControlContacts | INWAccessControlPassword
  ) => {},
  setPermissions: (permissions: INWPermissions) => {}
});

export default NewWorkspaceContext;

export const NewWorkspaceContextConsumer = NewWorkspaceContext.Consumer;
