import { Box, IconButton, TextField, Typography } from '@material-ui/core';
import ArrowBackRoundedIcon from '@material-ui/icons/ArrowBackRounded';
import SaveRoundedIcon from '@material-ui/icons/SaveRounded';
import { useFormik } from 'formik';
import { FC, useEffect, useState } from 'react';
import { useHistory, useParams } from 'react-router-dom';
import ContactsService from '../../../services/contacts/contactsService';
import { INewContactRequest } from '../../../services/contacts/contactsService.types';
import newContactValidationSchema from '../../../shared/schemas/contactSchemas';
import theme from '../../../theme/theme';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import Link from '../../atoms/Link/Link';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import useSnackbar from '../../molecules/Snackbar/useSnackbar.hook';
import KeyManager from '../../organisms/KeyManager/KeyManager';
import {
  EKeyInputType,
  EKeyType
} from '../../organisms/KeyManager/keyManager.types';
import {
  EContactEditType,
  IContactEditParams,
  IContactEditProps,
  IKeyPair
} from './contactEdit.types';

const ContactEdit: FC<IContactEditProps> = (props) => {
  const { type } = props;

  const pageTitle =
    type === EContactEditType.NEW ? 'New Contact' : 'Edit Contact';

  const [errorMessage, setErrorMessage] = useState<string>();
  const [addedKey, setAddedKey] = useState<IKeyPair>(null);
  const [initialValues, setInitialValues] = useState<{
    name: string;
    keyPair: IKeyPair;
  }>({
    name: '',
    keyPair: null
  });

  const { contactId } = useParams() as IContactEditParams;
  const { openSnackbar } = useSnackbar();

  const formik = useFormik({
    initialValues,
    enableReinitialize: true,
    validationSchema: newContactValidationSchema,
    onSubmit: (values, { resetForm }) => {
      setErrorMessage('');
      console.log(values);
      if (values.keyPair == null) {
        setErrorMessage('A key pair is required');
        return;
      }

      // Create the contact
      if (type == EContactEditType.NEW) {
        handleNewContact({
          name: values.name,
          publicKey: values.keyPair.publicKey
        }).catch((err) => {
          resetForm();

          openSnackbar('Unable to create new contact', 'error');
        });
      } else {
        // TODO
      }
    }
  });

  const history = useHistory();

  const handleNewContact = async (newContactRequest: INewContactRequest) => {
    await ContactsService.createContact(newContactRequest);

    history.push('/contacts');
    openSnackbar('Contact successfully created!', 'success');
  };

  useEffect(() => {
    if (type === EContactEditType.EDIT) {
      // Fetch contact information based on the passed in ID
      // TODO change mock

      setInitialValues({
        name: 'Milos Zivkovic',
        keyPair: { keyID: '123', publicKey: '123', privateKey: '123' }
      });

      setAddedKey({ keyID: '123', publicKey: '123', privateKey: '123' });
    }
  }, []);

  return (
    <Box display={'flex'} flexDirection={'column'}>
      <Box display={'flex'} alignItems={'center'}>
        {
          // TODO add unsaved changes modal
        }
        <Link to={'/contacts'}>
          <IconButton
            classes={{
              root: 'iconButtonRoot'
            }}
          >
            <ArrowBackRoundedIcon
              style={{
                fill: 'black'
              }}
            />
          </IconButton>
        </Link>
        <Box ml={2}>
          <PageTitle title={pageTitle} />
        </Box>
      </Box>
      <Box width={'100%'} display={'flex'} flexDirection={'column'}>
        <form autoComplete={'off'} onSubmit={formik.handleSubmit}>
          <Box
            display={'flex'}
            flexDirection={'column'}
            maxWidth={'30%'}
            mt={2}
          >
            <Box mb={2}>
              <FormTitle title={'Basic info'} />
            </Box>
            <TextField
              id={'name'}
              label={'Name'}
              variant={'outlined'}
              value={formik.values.name}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.name && Boolean(formik.errors.name)}
              helperText={formik.touched.name && formik.errors.name}
            />
          </Box>
          <Box display={'flex'} flexDirection={'column'} mt={2}>
            <Box mb={2}>
              <FormTitle title={'Public keys'} />
            </Box>
            <KeyManager
              addedKey={addedKey}
              setAddedKey={setAddedKey}
              formik={formik}
              expectedType={EKeyType.PUBLIC}
              visibleTypes={[EKeyInputType.IMPORT, EKeyInputType.ENTER]}
            />
          </Box>
          <Box display={errorMessage ? 'flex' : 'none'} mt={4}>
            <Typography
              variant={'body1'}
              style={{
                marginTop: '1rem',
                fontFamily: 'Montserrat',
                color: theme.palette.error.main
              }}
            >
              {errorMessage}
            </Typography>
          </Box>
          <Box mt={5} display={'flex'}>
            <ActionButton
              text={'Save'}
              square
              startIcon={<SaveRoundedIcon />}
            />
          </Box>
        </form>
      </Box>
    </Box>
  );
};

export default ContactEdit;
