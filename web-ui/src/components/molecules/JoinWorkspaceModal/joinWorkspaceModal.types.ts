import { IWorkspaceInfoResponse } from '../../../services/workspaces/workspacesService.types';

export interface IJoinWorkspaceModalProps {
  workspaceInfo: IWorkspaceInfoResponse;
  open: boolean;

  handleConfirm: (success: boolean) => void;
}
