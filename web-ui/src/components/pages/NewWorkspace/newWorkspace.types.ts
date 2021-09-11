export interface INewWorkspaceProps {}

export enum ENewWorkspaceStep {
  PARAMS = 'Parameters', // step 1
  SECURITY = 'Security', // step 2
  PERMISSIONS = 'Permissions', // step 3
  REVIEW = 'Review' // step 4
}
