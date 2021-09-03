import { Box } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import AddRoundedIcon from '@material-ui/icons/AddRounded';
import { FC } from 'react';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import ContactsTable from '../../molecules/ContactsTable/ContactsTable';
import { IContactsProps } from './contacts.types';

const Contacts: FC<IContactsProps> = () => {
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
        <ActionButton text={'New contact'} startIcon={<AddRoundedIcon />} />
      </Box>

      <ContactsTable />
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {};
});

export default Contacts;
