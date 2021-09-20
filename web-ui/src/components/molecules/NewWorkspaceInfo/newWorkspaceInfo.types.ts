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
  accessControl:
    | INWAccessControlContacts
    | INWAccessControlPassword
    | { contacts: string[] };
  permissions?: INWPermissions;
  accessControlType: ENWAccessControl;
}
