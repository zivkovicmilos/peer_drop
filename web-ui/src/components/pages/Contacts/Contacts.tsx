import { Box } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import AddRoundedIcon from '@material-ui/icons/AddRounded';
import { FC, useState } from 'react';
import { useHistory } from 'react-router-dom';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import ContactsConfirmModal from '../../molecules/ContactsConfirmModal/ContactsConfirmModal';
import ContactsTable from '../../molecules/ContactsTable/ContactsTable';
import { IContactConfirmInfo, IContactsProps } from './contacts.types';

const Contacts: FC<IContactsProps> = () => {
  const history = useHistory();

  const [confirmOpen, setConfirmOpen] = useState<boolean>(false);
  const [confirmInfo, setConfirmInfo] = useState<IContactConfirmInfo>({
    id: '1',
    name: 'Milos Zivkovic',
    publicKey: '12345'
  });

  const handleConfirm = (confirmed: boolean) => {
    console.log('Confirmed');

    setConfirmOpen(false);
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

      <ContactsTable handleDelete={handleDelete} />
      <ContactsConfirmModal
        contactInfo={confirmInfo}
        open={confirmOpen}
        handleConfirm={handleConfirm}
      />
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {};
});

export default Contacts;
