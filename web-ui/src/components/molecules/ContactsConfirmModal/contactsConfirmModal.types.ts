import { IContactConfirmInfo } from '../../pages/Contacts/contacts.types';

export interface IContactsConfirmModalProps {
  contactInfo: IContactConfirmInfo;
  open: boolean;

  handleConfirm: (success: boolean) => void;
}
