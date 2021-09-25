import { RestService } from '../rest/restService';
import { IPagination } from '../rest/restService.types';
import {
  IDownloadFileRequest,
  IJoinWorkspaceRequest,
  INewWorkspaceRequest,
  INewWorkspaceResponse,
  IWorkspaceDetailedResponse,
  IWorkspaceInfoResponse,
  IWorkspaceListResponse,
  IWorkspaceNumPeersResponse
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

  public static async leaveWorkspace(mnemonic: string): Promise<string> {
    try {
      return await RestService.delete<string>({
        url: `workspaces/${mnemonic}`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async downloadFile(
    request: IDownloadFileRequest
  ): Promise<any> {
    try {
      return await RestService.post<any>({
        url: `workspaces/download`,
        data: {
          ...request
        },
        config: {
          responseType: 'blob'
        }
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async uploadWorkspaceFile(formData: FormData): Promise<string> {
    try {
      return await RestService.post<string>({
        url: `workspaces/upload`,
        data: formData
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
        url: `workspaces/${mnemonic}`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async getWorkspaces(
    pagination: IPagination
  ): Promise<IWorkspaceListResponse> {
    try {
      return await RestService.get<IWorkspaceListResponse>({
        url: `workspaces?page=${pagination.page}&limit=${pagination.limit}`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async getWorkspaceFiles(
    mnemonic: string
  ): Promise<IWorkspaceDetailedResponse> {
    try {
      return await RestService.get<IWorkspaceDetailedResponse>({
        url: `workspaces/${mnemonic}/files`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async getWorkspacePeers(
    mnemonic: string
  ): Promise<IWorkspaceNumPeersResponse> {
    try {
      return await RestService.get<IWorkspaceNumPeersResponse>({
        url: `workspaces/${mnemonic}/peers`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }
}

export default WorkspacesService;
