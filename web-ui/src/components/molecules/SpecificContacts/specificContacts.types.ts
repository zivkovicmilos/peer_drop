import { IContactResponse } from '../../../services/contacts/contactsService.types';

export interface ISpecificContactsProps {
  contactIDs: { contacts: IContactResponse[] };
  setContactIDs: (newContacts: { contacts: IContactResponse[] }) => void;

  errorMessage: string;

  listTitle?: string;
  disabled?: boolean;
}
