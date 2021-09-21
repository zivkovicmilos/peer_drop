import { IWorkspaceDetailedResponse } from '../../../services/workspaces/workspacesService.types';
import { IWorkspaceInfo } from '../../pages/ViewWorkspace/viewWorkspace.types';

export interface IViewWorkspaceFilesProps {
  workspaceInfo: IWorkspaceDetailedResponse;
}

export enum FILE_TYPE {
  FILE,
  FOLDER
}
