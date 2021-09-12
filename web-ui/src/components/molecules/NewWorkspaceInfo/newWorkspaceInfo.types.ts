import {
  ENewWorkspaceType,
  ENWAccessControl,
  INWAccessControlContacts,
  INWAccessControlPassword,
  INWPermissions
} from '../../../context/newWorkspaceContext.types';

export interface INewWorkspaceInfoProps {
  workspaceName: string;
  workspaceType: ENewWorkspaceType;
  accessControl: INWAccessControlContacts | INWAccessControlPassword;
  permissions: INWPermissions;
  accessControlType: ENWAccessControl;
}
