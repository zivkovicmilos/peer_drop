import { createContext } from 'react';
import {
  ENewWorkspaceType,
  ENWAccessControl,
  INWAccessControlContacts,
  INWAccessControlPassword,
  INWPermissions
} from './newWorkspaceContext.types';

export interface INewWorkspaceContext {
  step: number;
  setStep: (newStep: number) => void;

  workspaceName: string;
  workspaceType: ENewWorkspaceType;
  accessControl: INWAccessControlContacts | INWAccessControlPassword;
  accessControlType: ENWAccessControl;
  permissions: INWPermissions;

  setWorkspaceName: (name: string) => void;
  setWorkspaceType: (type: ENewWorkspaceType) => void;
  setAccessControl: (
    control: INWAccessControlContacts | INWAccessControlPassword
  ) => void;
  setPermissions: (permissions: INWPermissions) => void;
  setAccessControlType: (type: ENWAccessControl) => void;

  handleBack: () => void;
  handleNext: () => void;
}

const NewWorkspaceContext = createContext<INewWorkspaceContext>({
  step: 0,
  setStep: (newStep: number) => {},

  workspaceName: null,
  workspaceType: null,
  accessControl: null,
  accessControlType: null,
  permissions: null,

  setWorkspaceName: (name: string) => {},
  setWorkspaceType: (type: ENewWorkspaceType) => {},
  setAccessControl: (
    control: INWAccessControlContacts | INWAccessControlPassword
  ) => {},
  setPermissions: (permissions: INWPermissions) => {},
  setAccessControlType: (type: ENWAccessControl) => {},

  handleBack: () => {},
  handleNext: () => {}
});

export default NewWorkspaceContext;

export const NewWorkspaceContextConsumer = NewWorkspaceContext.Consumer;
