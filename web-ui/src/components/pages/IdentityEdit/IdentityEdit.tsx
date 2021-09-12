import {
  Avatar,
  Box,
  IconButton,
  TextField,
  Typography
} from '@material-ui/core';
import AddRoundedIcon from '@material-ui/icons/AddRounded';
import ArrowBackRoundedIcon from '@material-ui/icons/ArrowBackRounded';
import CameraAltRoundedIcon from '@material-ui/icons/CameraAltRounded';
import SaveRoundedIcon from '@material-ui/icons/SaveRounded';
import { Formik } from 'formik';
import { FC, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import newContactValidationSchema from '../../../shared/schemas/contactSchemas';
import theme from '../../../theme/theme';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import Link from '../../atoms/Link/Link';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import useSnackbar from '../../molecules/Snackbar/useSnackbar.hook';
import KeyManager from '../../organisms/KeyManager/KeyManager';
import { EKeyInputType } from '../../organisms/KeyManager/keyManager.types';
import { IKeyPair } from '../ContactEdit/contactEdit.types';

import {
  EIdentityEditType,
  IIdentityEditParams,
  IIdentityEditProps
} from './identityEdit.types';

const IdentityEdit: FC<IIdentityEditProps> = (props) => {
  const { type } = props;

  const pageTitle =
    type === EIdentityEditType.NEW ? 'New Identity' : 'Edit Identity';

  const [errorMessage, setErrorMessage] = useState<string>();
  const [addedKey, setAddedKey] = useState<IKeyPair>(null);
  const [initialValues, setInitialValues] = useState<{
    name: string;
    keyPair: IKeyPair;
  }>({
    name: '',
    keyPair: null
  });

  const { identityId } = useParams() as IIdentityEditParams;
  const { openSnackbar } = useSnackbar();

  useEffect(() => {
    if (type === EIdentityEditType.EDIT) {
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
        <Link to={'/identities'}>
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
        <Formik
          initialValues={initialValues}
          enableReinitialize={true}
          validationSchema={newContactValidationSchema}
          onSubmit={(values: any, { resetForm }: any) => {
            setErrorMessage('');

            if (values.keyPair == null) {
              setErrorMessage('A key pair is required');
              return;
            }

            console.log('Identity form submitted');
          }}
        >
          {(formik) => (
            <form autoComplete={'off'} onSubmit={formik.handleSubmit}>
              <Box display={'flex'} maxWidth={'100%'} mt={2}>
                <Box display={'flex'} flexDirection={'column'} maxWidth={'30%'}>
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
                <Box
                  display={'flex'}
                  flexDirection={'column'}
                  width={'auto'}
                  alignItems={'center'}
                  justifyContent={'center'}
                  ml={6}
                >
                  <Avatar
                    style={{
                      width: '70px',
                      height: '70px',
                      marginBottom: '16px'
                    }}
                  />
                  <ActionButton
                    text={'Select image'}
                    startIcon={<CameraAltRoundedIcon />}
                    shouldSubmit={false}
                  />
                </Box>
              </Box>
              <Box display={'flex'} flexDirection={'column'} mt={2}>
                <Box mb={2}>
                  <FormTitle title={'Key pair'} />
                </Box>
                <KeyManager
                  addedKey={addedKey}
                  keyListTitle={
                    type == EIdentityEditType.EDIT
                      ? 'Existing pair'
                      : 'Created pair'
                  }
                  setAddedKey={setAddedKey}
                  visibleTypes={[
                    EKeyInputType.IMPORT,
                    EKeyInputType.ENTER,
                    EKeyInputType.GENERATE
                  ]}
                  formik={formik}
                />
              </Box>
              <Box display={errorMessage ? 'flex' : 'none'} mt={2}>
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
                  text={type == EIdentityEditType.EDIT ? 'Save' : 'Create'}
                  square
                  startIcon={
                    type == EIdentityEditType.EDIT ? (
                      <SaveRoundedIcon />
                    ) : (
                      <AddRoundedIcon />
                    )
                  }
                />
              </Box>
            </form>
          )}
        </Formik>
      </Box>
    </Box>
  );
};

export default IdentityEdit;
