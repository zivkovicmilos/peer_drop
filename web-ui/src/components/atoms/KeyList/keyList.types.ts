import { IKeyPair } from '../../pages/ContactEdit/contactEdit.types';

export interface IKeyListProps {
  addedKey: IKeyPair;
  handleKeyRemove: (key: IKeyPair) => void;
}
