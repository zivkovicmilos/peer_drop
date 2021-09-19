import { RestService } from '../rest/restService';
import {
  INewWorkspaceRequest,
  INewWorkspaceResponse
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
}

export default WorkspacesService;