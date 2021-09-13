import { ENewWorkspaceType } from '../../../context/newWorkspaceContext.types';

export interface IViewWorkspaceProps {}

export interface IWorkspaceInfo {
  name: string;
  id: string;
  type: ENewWorkspaceType
}