import { WorkspaceInfo } from '../../pages/JoinWorkspace/joinWorkspace.types';

export interface IJoinWorkspaceModalProps {
  workspaceInfo: WorkspaceInfo;
  open: boolean;

  handleConfirm: (success: boolean) => void;
}
