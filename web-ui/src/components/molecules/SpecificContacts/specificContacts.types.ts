import { ContactResponse } from '../../../context/newWorkspaceContext.types';

export interface ISpecificContactsProps {
  contactIDs: { contacts: ContactResponse[] };
  setContactIDs: (newContacts: { contacts: ContactResponse[] }) => void;

  errorMessage: string;

  listTitle?: string;
  disabled?: boolean;
}