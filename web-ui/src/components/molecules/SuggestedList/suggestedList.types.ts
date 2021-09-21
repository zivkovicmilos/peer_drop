import { IWorkspaceDetailedFileResponse } from '../../../services/workspaces/workspacesService.types';

export interface ISuggestedListProps {
  files: IWorkspaceDetailedFileResponse[];
}

export interface IFileInfo {
  id: string;
  name: string;
  extension: string;
}
