import { Box, IconButton, TextField } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import ArrowBackRoundedIcon from '@material-ui/icons/ArrowBackRounded';
import { FC } from 'react';
import { useParams } from 'react-router-dom';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import Link from '../../atoms/Link/Link';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import KeyManager from '../../organisms/KeyManager/KeyManager';
import { EKeyInputType } from '../../organisms/KeyManager/keyManager.types';
import { IContactEditParams, IContactEditProps } from './contactEdit.types';

const ContactEdit: FC<IContactEditProps> = (props) => {
  const { type } = props;

  const { contactId } = useParams() as IContactEditParams;

  return (
    <Box display={'flex'} flexDirection={'column'}>
      <Box display={'flex'} alignItems={'center'}>
        {
          // TODO add unsaved changes modal
        }
        <Link to={'/contacts'}>
          <IconButton>
            <ArrowBackRoundedIcon
              style={{
                fill: 'black'
              }}
            />
          </IconButton>
        </Link>
        <Box>
          <PageTitle title={'New Contact'} />
        </Box>
      </Box>
      <Box width={'100%'}>
        <form autoComplete={'off'}>
          <Box
            display={'flex'}
            flexDirection={'column'}
            maxWidth={'30%'}
            mt={2}
          >
            <Box mb={2}>
              <FormTitle title={'Basic info'} />
            </Box>
            <TextField id={'name'} label={'Name'} variant={'outlined'} />
          </Box>
          <Box
            display={'flex'}
            flexDirection={'column'}
            mt={2}
          >
            <Box mb={2}>
              <FormTitle title={'Public keys'} />
            </Box>
            <KeyManager
              visibleTypes={[EKeyInputType.IMPORT, EKeyInputType.ENTER]}
            />
          </Box>
        </form>
      </Box>
    </Box>
  );
};

export default ContactEdit;
