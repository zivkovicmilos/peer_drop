import { RestService } from '../rest/restService';
import {
  IJoinWorkspaceRequest,
  INewWorkspaceRequest,
  INewWorkspaceResponse,
  IWorkspaceInfoResponse
} from './workspacesService.types';

class WorkspacesService {
  public static async createWorkspace(
    request: INewWorkspaceRequest
  ): Promise<INewWorkspaceResponse> {
    try {
      return await RestService.post<INewWorkspaceResponse>({
        url: `workspaces`,
        data: {
          ...request
        }
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async joinWorkspace(
    request: IJoinWorkspaceRequest
  ): Promise<string> {
    try {
      return await RestService.post<string>({
        url: `join-workspace`,
        data: {
          ...request
        }
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async getWorkspaceInfo(
    mnemonic: string
  ): Promise<IWorkspaceInfoResponse> {
    try {
      return await RestService.get<IWorkspaceInfoResponse>({
        url: `workspaces/${mnemonic.replace(/\s/g, '-')}`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }
}

export default WorkspacesService;
