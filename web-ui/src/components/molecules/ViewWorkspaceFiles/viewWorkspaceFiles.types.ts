import { IWorkspaceInfo } from '../../pages/ViewWorkspace/viewWorkspace.types';

export interface IViewWorkspaceFilesProps {
  workspaceInfo: IWorkspaceInfo;
}

export enum FILE_TYPE {
  FILE,
  FOLDER
}
