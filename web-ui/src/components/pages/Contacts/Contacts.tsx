import { Box } from '@material-ui/core';
import AddRoundedIcon from '@material-ui/icons/AddRounded';
import { FC, useState } from 'react';
import { useHistory } from 'react-router-dom';
import ContactsService from '../../../services/contacts/contactsService';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import ContactsConfirmModal from '../../molecules/ContactsConfirmModal/ContactsConfirmModal';
import ContactsTable from '../../molecules/ContactsTable/ContactsTable';
import useSnackbar from '../../molecules/Snackbar/useSnackbar.hook';
import { IContactConfirmInfo, IContactsProps } from './contacts.types';

const Contacts: FC<IContactsProps> = () => {
  const history = useHistory();

  const [confirmOpen, setConfirmOpen] = useState<boolean>(false);
  const [confirmInfo, setConfirmInfo] = useState<IContactConfirmInfo>({
    id: '',
    name: '',
    publicKey: ''
  });

  const { openSnackbar } = useSnackbar();

  const [fetchTrigger, setFetchTrigger] = useState<boolean>(false);

  const handleContactDelete = async (id: string) => {
    await ContactsService.deleteContact(id);
  };

  const handleConfirm = (confirmed: boolean) => {
    setConfirmOpen(false);

    if (confirmed) {
      handleContactDelete(confirmInfo.id)
        .then(() => {
          openSnackbar('Successfully deleted contact', 'success');

          setFetchTrigger(!fetchTrigger);
        })
        .catch((err) => {
          openSnackbar('Unable to delete contact', 'error');
        });
    }
  };

  const handleDelete = (contactInfo: IContactConfirmInfo) => {
    setConfirmInfo({
      id: contactInfo.id,
      name: contactInfo.name,
      publicKey: contactInfo.publicKey
    });

    setConfirmOpen(true);
  };

  const handleNewContact = () => {
    history.push('/contacts/new');
  };

  return (
    <Box
      display={'flex'}
      flexDirection={'column'}
      width={'100%'}
      height={'100%'}
    >
      <Box
        display={'flex'}
        justifyContent={'space-between'}
        width={'100%'}
        alignItems={'center'}
        mb={4}
      >
        <PageTitle title={'Contacts'} />
        <ActionButton
          text={'New contact'}
          startIcon={<AddRoundedIcon />}
          onClick={handleNewContact}
        />
      </Box>

      <ContactsTable handleDelete={handleDelete} fetchTrigger={fetchTrigger} />
      <ContactsConfirmModal
        contactInfo={confirmInfo}
        open={confirmOpen}
        handleConfirm={handleConfirm}
      />
    </Box>
  );
};

export default Contacts;
