import { IContactConfirmInfo } from '../../pages/Contacts/contacts.types';

export interface IContactsTableProps {
  handleDelete: (contactInfo: IContactConfirmInfo) => void;
}
