import { IContactResponse } from '../contacts/contactsService.types';
import { IIdentityResponse } from '../identities/identitiesService.types';
import { IWorkspaceWrapper } from '../workspaces/workspacesService.types';

export interface ISearchResults {
  workspaces: IWorkspaceWrapper[];
  contacts: IContactResponse[];
  identities: IIdentityResponse[];
}